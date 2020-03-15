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

	mockChatConfig := types.ChatConfig {
		ProjectDescription: "Initial Project Description",
	}

	handlePrivateMessage("Some message", "username1", mockSay, mockChatConfig)

	if(sayCalled) {
		t.Fail()
	}
}

func Test_HandlePrivateMessage_invalidCommandInput(t *testing.T) {

	sayCalled := false

	mockSay := func(message string) error {
		sayCalled = true
		return nil
	}

	mockChatConfig := types.ChatConfig {
		ProjectDescription: "Initial Project Description",
	}

	handlePrivateMessage("!invalidcommand", "username1", mockSay, mockChatConfig)

	if(sayCalled) {
		t.Fail()
	}
}

func Test_HandlePrivateMessage_helloCommandInput(t *testing.T) {

	sayCalled := false

	mockSay := func(message string) error {
		sayCalled = true
		if(message != "Hello!") {
			t.Fail()
		}
		return nil
	}

	mockChatConfig := types.ChatConfig {
		ProjectDescription: "Initial Project Description",
	}

	handlePrivateMessage("!hello", "username1", mockSay, mockChatConfig)

	if(!sayCalled) {
		t.Fail()
	}
}

func Test_HandlePrivateMessage_projectCommandInput(t *testing.T) {

	sayCalled := false

	mockChatConfig := types.ChatConfig {
		ProjectDescription: "Initial Project Description",
	}

	mockSay := func(message string) error {
		sayCalled = true
		if(message != mockChatConfig.ProjectDescription) {
			t.Fail()
		}
		return nil
	}

	handlePrivateMessage("!project", "username1", mockSay, mockChatConfig)

	if(!sayCalled) {
		t.Fail()
	}
}