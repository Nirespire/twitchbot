package main

import (
	"time"

	"github.com/Nirespire/twitchbot/bot"
	"github.com/Nirespire/twitchbot/types"
)

func main() {

	chatConfig := types.ChatConfig{
		ProjectDescription: "Currently working on a twitch chatbot using GOLANG.",
	}

	myBot := bot.BasicBot{
		Channel:     "nirespire",
		MsgRate:     time.Duration(20/30) * time.Millisecond,
		Name:        "NirespireBot",
		Port:        "6667",
		PrivatePath: "oauth.json",
		Server:      "irc.chat.twitch.tv",
		ServerPort:  ":8080",
		ChatConfig:  chatConfig,
	}
	myBot.Start()
}
