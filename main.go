package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"time"
	"fmt"
	"net"
	"net/textproto"
	"regexp"
	"strings"
)

const ESTFormat = "Jan 8 10:28:00 EST"
var msgRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG) #\w+(?: :(.*))?$`)
var cmdRegex *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)



// TODO:
// 1. Connect to a Twitch.tv Chat channel.
//  a. Pass along necessary information for the connection.
//   i.   The IRC (chat) server.
//   ii.  The port on the server.
//   iii. The channel we want the bot to join.
//   iv.  The bot's name.
//   v.   A secure key to allow the bot to connect indirectly (not through the website).
//   vi.  A maximum speed at which the bot can respond.
// 2. Listen for messages in the chat.
// 3. Do things based on what is happening in the chat.

func TimeStamp(format string) string {
	return time.Now().Format(format)
}

func timeStamp() string {
	return TimeStamp(ESTFormat)
}

type OAuthCred struct {
	Password string `json:"password,omitempty"`
}

type TwitchBot interface {
	Connect()
	Disconnect()
	HandleChat()
	JoinChannel()
	ReadCredentials() error
	Say(msg string) error
	Start()
}

type BasicBot struct {
	Channel string
	conn net.Conn
	Credentials *OAuthCred
	MsgRate time.Duration
	Name string
	Port string
	PrivatePath string
	Server string
	startTime time.Time
}

func (bb *BasicBot) Connect() {
	var err error
	fmt.Printf("[%s] Connecting to %s...\n", timeStamp(), bb.Server)

	bb.conn, err = net.Dial("tcp", bb.Server+":"+bb.Port)
	if err != nil {
		fmt.Printf("[%s] Cannont cont to %s, retrying.\n", timeStamp(), bb.Server)
		bb.Connect()
	}

	fmt.Printf("[%s] Connected to %s", timeStamp(), bb.Server)
	bb.startTime = time.Now()
}

func (bb *BasicBot) Disconnect() {
	bb.conn.Close()
	upTime := time.Now().Sub(bb.startTime).Seconds()
	fmt.Printf("[%s] Closed connected from %s | Live for: %fs\n", timeStamp(), bb.Server, upTime)
}

func (bb *BasicBot) ReadCredentials() error {
	credFile, err := ioutil.ReadFile(bb.PrivatePath)
	if err != nil {
		return err
	}

	bb.Credentials = &OAuthCred{}

	dec := json.NewDecoder(strings.NewReader(string(credFile)))
	
	if err = dec.Decode(bb.Credentials); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (bb *BasicBot) JoinChannel() {
	fmt.Printf("[%s] Joining #%s...\n", timeStamp(), bb.Channel)
	bb.conn.Write([]byte("PASS " + bb.Credentials.Password + "\r\n"))
    bb.conn.Write([]byte("NICK " + bb.Name + "\r\n"))
	bb.conn.Write([]byte("JOIN #" + bb.Channel + "\r\n"))
	
	fmt.Printf("[%s] Joined #%s as @%s!\n", timeStamp(), bb.Channel, bb.Name)
}

func (bb *BasicBot) Say(msg string) error {
	if msg == "" {
		return errors.New("BasicBot.Say: msg was empty.")
	}
	_, err := bb.conn.Write([]byte(fmt.Sprintf("PRIVMSG #%s %s\r\n", bb.Channel, msg)))
	if err != nil {
		return err
	}
	return nil
}

func (bb *BasicBot) HandleChat() error {
	fmt.Printf("[%s] Watching #%s...\n", timeStamp(), bb.Channel)

	tp := textproto.NewReader(bufio.NewReader(bb.conn))

	for {
		line, err := tp.ReadLine()

		if nil != err {
            bb.Disconnect()
            return errors.New("bb.Bot.HandleChat: Failed to read line from channel. Disconnected.")
        }

		fmt.Printf("[%s] %s\n", timeStamp(), line)

		if "PING :tmi.twitch.tv" == line {
            bb.conn.Write([]byte("PONG :tmi.twitch.tv\r\n"))
            continue
		} else {
			matches := msgRegex.FindStringSubmatch(line)

			if matches != nil {
				userName := matches[1]
				msgType := matches[2]

				switch msgType {
				case "PRIVMSG":
					msg := matches[3]
					fmt.Printf("[%s] %s: %s\n", timeStamp(), userName, msg)

					cmdMatches := cmdRegex.FindStringSubmatch(msg)

					if cmdMatches != nil {
						cmd := cmdMatches[1]

						switch cmd {
						case "hello":
							fmt.Printf("[%s] %s said hello!", timeStamp(), userName)
							bb.Say("Hello!")
						}
					}
				}
			}
		}
		time.Sleep(bb.MsgRate)
	}
}

func (bb *BasicBot) Start() {
	err := bb.ReadCredentials()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Aborting...")
		return
	}

	for {
		bb.Connect()
		bb.JoinChannel()
		err = bb.HandleChat()
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
            fmt.Println(err)
            fmt.Println("Starting bot again...")
		} else {
			return
		}
	}
}

func main() {

	myBot := BasicBot{
		Channel:     "nirespire",
		MsgRate:     time.Duration(20/30) * time.Millisecond,
		Name:        "NirespireBot",
		Port:        "6667",
		PrivatePath: "oauth.json",
		Server:      "irc.chat.twitch.tv",
	}
	myBot.Start()
}
