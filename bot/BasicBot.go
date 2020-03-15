package bot

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/textproto"
	"regexp"
	"strings"
	"time"

	"github.com/Nirespire/twitchbot/types"
	"github.com/Nirespire/twitchbot/web"
)

var msgRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG) #\w+(?: :(.*))?$`)
var cmdRegex *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)

type ChatBot interface {
	connect()
	disconnect()
	handleChat()
	joinChannel()
	readCredentials() error
	say(string) error
	getChatConfig() types.ChatConfig
	Start()
	startWebServer() error
}

// TwitchBot defines the basic twitchbot details to connect and join a channel
type TwitchBot struct {
	Channel     string
	conn        net.Conn
	Credentials *types.OAuthCred
	MsgRate     time.Duration
	Name        string
	Port        string
	PrivatePath string
	Server      string
	startTime   time.Time
	ServerPort  string
	ChatConfig  types.ChatConfig
}

func (bb *TwitchBot) connect() {
	var err error
	log.Printf("connecting to %s...\n", bb.Server)

	bb.conn, err = net.Dial("tcp", bb.Server+":"+bb.Port)
	if err != nil {
		log.Printf("Cannont cont to %s, retrying.\n", bb.Server)
		bb.connect()
	}

	log.Printf("connected to %s", bb.Server)
	bb.startTime = time.Now()
}

func (bb *TwitchBot) disconnect() {
	bb.conn.Close()
	upTime := time.Now().Sub(bb.startTime).Seconds()
	log.Printf("Closed connected from %s | Live for: %fs\n", bb.Server, upTime)
}

func (bb *TwitchBot) readCredentials() error {
	credFile, err := ioutil.ReadFile(bb.PrivatePath)
	if err != nil {
		return err
	}

	bb.Credentials = &types.OAuthCred{}

	dec := json.NewDecoder(strings.NewReader(string(credFile)))

	if err = dec.Decode(bb.Credentials); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (bb *TwitchBot) joinChannel() {
	log.Printf("Joining #%s...\n", bb.Channel)
	bb.conn.Write([]byte("PASS " + bb.Credentials.Password + "\r\n"))
	bb.conn.Write([]byte("NICK " + bb.Name + "\r\n"))
	bb.conn.Write([]byte("JOIN #" + bb.Channel + "\r\n"))

	log.Printf("Joined #%s as @%s!\n", bb.Channel, bb.Name)
}

func (bb *TwitchBot) Say(msg string) error {
	if msg == "" {
		return errors.New("TwitchBot.say: msg was empty")
	}

	log.Printf("DEBUG PRIVMSG #%s :%s\r\n", bb.Channel, msg)

	_, err := bb.conn.Write([]byte(fmt.Sprintf("PRIVMSG #%s :%s\r\n", bb.Channel, msg)))
	if err != nil {
		log.Print("Failed to write to output connection")
		log.Print(err)
		return err
	}
	return nil
}

func (bb *TwitchBot) handleChat() error {
	log.Printf("Watching #%s...\n", bb.Channel)

	tp := textproto.NewReader(bufio.NewReader(bb.conn))

	for {
		line, err := tp.ReadLine()

		if nil != err {
			bb.disconnect()
			return errors.New("bb.Bot.handleChat: Failed to read line from channel. disconnected")
		}

		log.Printf("%s\n", line)

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
					handlePrivateMessage(matches[3], userName, bb.Say, bb.ChatConfig)
				}
			}
		}
		time.Sleep(bb.MsgRate)
	}
}

func handlePrivateMessage(message string, userName string, say func(message string) error, chatConfig types.ChatConfig) {

	log.Printf("%s: %s\n", userName, message)

	cmdMatches := cmdRegex.FindStringSubmatch(message)

	if cmdMatches != nil {
		log.Printf("%s sent a command %s\n", userName, cmdMatches)

		cmd := cmdMatches[1]

		switch cmd {
		case "hello":
			log.Printf("%s said hello!", userName)
			say("Hello!")
		case "project":
			log.Printf("%s sent the project command!\n", userName)
			say(chatConfig.ProjectDescription)
		}
	}
}

func (bb *TwitchBot) getChatConfig() types.ChatConfig {
	return bb.ChatConfig
}

func (bb *TwitchBot) startWebServer() {

	webserver := web.ServerConfig{
		BotConfig: &(bb.ChatConfig),
		Port:      bb.ServerPort,
	}

	webserver.StartWebServer()
}

// Start initializes and runs the twitchbot with the provided configuration
func (bb *TwitchBot) Start() {
	err := bb.readCredentials()
	if err != nil {
		log.Println(err)
		log.Println("Aborting...")
		return
	}

	for {
		bb.startWebServer()
		bb.connect()
		bb.joinChannel()
		err = bb.handleChat()
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			log.Println(err)
			log.Println("Starting bot again...")
		} else {
			return
		}
	}
}
