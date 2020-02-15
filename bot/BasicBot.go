package bot

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/textproto"
	"regexp"
	"strings"
	"time"

	"github.com/Nirespire/twitchbot/util"
)

var msgRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG) #\w+(?: :(.*))?$`)
var cmdRegex *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)

type twitchBot interface {
	connect()
	disconnect()
	handleChat()
	joinChannel()
	readCredentials() error
	say(msg string) error
	Start()
}

type oAuthCred struct {
	Password string `json:"password,omitempty"`
}

// BasicBot defines the basic twitchbot details to connect and join a channel
// as well as behavior configuration
type BasicBot struct {
	Channel     string
	conn        net.Conn
	Credentials *oAuthCred
	MsgRate     time.Duration
	Name        string
	Port        string
	PrivatePath string
	Server      string
	startTime   time.Time
}

func (bb *BasicBot) connect() {
	var err error
	fmt.Printf("[%s] connecting to %s...\n", util.TimeStamp(), bb.Server)

	bb.conn, err = net.Dial("tcp", bb.Server+":"+bb.Port)
	if err != nil {
		fmt.Printf("[%s] Cannont cont to %s, retrying.\n", util.TimeStamp(), bb.Server)
		bb.connect()
	}

	fmt.Printf("[%s] connected to %s", util.TimeStamp(), bb.Server)
	bb.startTime = time.Now()
}

func (bb *BasicBot) disconnect() {
	bb.conn.Close()
	upTime := time.Now().Sub(bb.startTime).Seconds()
	fmt.Printf("[%s] Closed connected from %s | Live for: %fs\n", util.TimeStamp(), bb.Server, upTime)
}

func (bb *BasicBot) readCredentials() error {
	credFile, err := ioutil.ReadFile(bb.PrivatePath)
	if err != nil {
		return err
	}

	bb.Credentials = &oAuthCred{}

	dec := json.NewDecoder(strings.NewReader(string(credFile)))

	if err = dec.Decode(bb.Credentials); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (bb *BasicBot) joinChannel() {
	fmt.Printf("[%s] Joining #%s...\n", util.TimeStamp(), bb.Channel)
	bb.conn.Write([]byte("PASS " + bb.Credentials.Password + "\r\n"))
	bb.conn.Write([]byte("NICK " + bb.Name + "\r\n"))
	bb.conn.Write([]byte("JOIN #" + bb.Channel + "\r\n"))

	fmt.Printf("[%s] Joined #%s as @%s!\n", util.TimeStamp(), bb.Channel, bb.Name)
}

func (bb *BasicBot) say(msg string) error {
	if msg == "" {
		return errors.New("BasicBot.say: msg was empty")
	}
	_, err := bb.conn.Write([]byte(fmt.Sprintf("PRIVMSG #%s %s\r\n", bb.Channel, msg)))
	if err != nil {
		return err
	}
	return nil
}

func (bb *BasicBot) handleChat() error {
	fmt.Printf("[%s] Watching #%s...\n", util.TimeStamp(), bb.Channel)

	tp := textproto.NewReader(bufio.NewReader(bb.conn))

	for {
		line, err := tp.ReadLine()

		if nil != err {
			bb.disconnect()
			return errors.New("bb.Bot.handleChat: Failed to read line from channel. disconnected")
		}

		fmt.Printf("[%s] %s\n", util.TimeStamp(), line)

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
					fmt.Printf("[%s] %s: %s\n", util.TimeStamp(), userName, msg)

					cmdMatches := cmdRegex.FindStringSubmatch(msg)

					if cmdMatches != nil {
						cmd := cmdMatches[1]

						switch cmd {
						case "hello":
							fmt.Printf("[%s] %s said hello!", util.TimeStamp(), userName)
							bb.say("Hello!")
						}
					}
				}
			}
		}
		time.Sleep(bb.MsgRate)
	}
}

// Start initializes and runs the twitchbot with the provided configuration
func (bb *BasicBot) Start() {
	err := bb.readCredentials()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Aborting...")
		return
	}

	for {
		bb.connect()
		bb.joinChannel()
		err = bb.handleChat()
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			fmt.Println(err)
			fmt.Println("Starting bot again...")
		} else {
			return
		}
	}
}
