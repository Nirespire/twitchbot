# Twitchbot GO

Livecoding a Twitchbot at https://twitch.tv/nirespire

Following the tutorial here: https://dev.to/foresthoffman/building-a-twitchtv-chat-bot-with-go---part-1-i3k

## Setup

1. Generate OAuth credentials for twitch IRC authentication
2. Copy OAuth token into `oauth.json.template` and rename to `oauth.json`
3. Set appropriate configs in `main.go`
```go

chatConfig := types.ChatConfig{
	ProjectDescription: "Currently working on a twitch chatbot using GOLANG.",
}

myBot := bot.BasicBot{
		Channel:     "your_channel_name(lowercase)",
		MsgRate:     time.Duration(20/30) * time.Millisecond,
		Name:        "SomeBotName",
		Port:        "6667",
		PrivatePath: "oauth.json",
		Server:      "irc.chat.twitch.tv",
		ServerPort:  ":8080",
		ChatConfig:  chatConfig,
	}
```
4. `go install`
5. run `twitchbot`