package util

import (
	"testing"
)

func Test_parseUserState_singleTag(t *testing.T) {
	testString := "@color=#0D4200 :tmi.twitch.tv USERSTATE #dallas"

	result := parseUserState(testString)

	if(len(result) != 1) {
		t.Fail()
	}

	if(result["color"] != "#0D4200") {
		t.Fail()
	}
}

func Test_parseUserState_multiTag(t *testing.T) {
	testString := "@badge-info=;badges=staff/1;color=#0D4200;display-name=ronni;emote-sets=0,33,50,237,793,2126,3517,4578,5569,9400,10337,12239;mod=1;subscriber=1;turbo=1;user-type=staff :tmi.twitch.tv USERSTATE #dallas"

	result := parseUserState(testString)

	if(len(result) != 9) {
		t.Fail()
	}
}