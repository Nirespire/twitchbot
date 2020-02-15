package main

import (
	"github.com/Nirespire/twitchbot/bot"
	"time"
)

func main() {
	myBot := bot.BasicBot{
		Channel:     "nirespire",
		MsgRate:     time.Duration(20/30) * time.Millisecond,
		Name:        "NirespireBot",
		Port:        "6667",
		PrivatePath: "oauth.json",
		Server:      "irc.chat.twitch.tv",
		ServerPort:  ":8080",
	}
	myBot.Start()
}
