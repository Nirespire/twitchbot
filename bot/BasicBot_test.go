package bot

import (
	"testing"
	"net"
)

// 1
type mockNet struct {
	net.Conn
}

// 2
func (m *mockNet) Dial(protocol string, address string) (net.Conn, error) {
	// 3
	return nil, nil
}


func TestConnect(t *testing.T) {
	// 4
	client := &mockNet{}

	myBot := bot.BasicBot {
		Channel:     "nirespire",
		MsgRate:     time.Duration(20/30) * time.Millisecond,
		Name:        "NirespireBot",
		Port:        "6667",
		PrivatePath: "oauth.json",
		Server:      "irc.chat.twitch.tv",
		ServerPort:  ":8080",
	}
	
	// 5
	myBot.connect()

	if(myBot.conn == nil) {
		t.Error("Test failed")
	}
}