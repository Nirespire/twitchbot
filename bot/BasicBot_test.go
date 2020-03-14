package bot

import (
	"testing"
	"github.com/Nirespire/twitchbot/types"
)


func Test_HandlePrivateMessage_nonCommandInput(t *testing.T) {

	sayCalled := false

	mockSay := func(message string) error {
		sayCalled = true
		return nil
	}

	mockChatConfig := types.ChatConfig{
		ProjectDescription: "Initial Project Description",
	}

	handlePrivateMessage("Some message", "username1", mockSay, mockChatConfig)

	if(sayCalled) {
		t.Fail()
	}

	if(mockChatConfig.ProjectDescription != "Initial Project Description") {
		t.Fail()
	}
}