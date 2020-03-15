package util

import (
	"regexp"
	"strings"
)

// @badge-info=;badges=staff/1;color=#0D4200;display-name=ronni;emote-sets=0,33,50,237,793,2126,3517,4578,5569,9400,10337,12239;mod=1;subscriber=1;turbo=1;user-type=staff :tmi.twitch.tv USERSTATE #dallas


var tagRegex *regexp.Regexp = regexp.MustCompile(`@(.*?)[ ]`)

func ParseUserState(chatMessage string) map[string]string {

	tagString := tagRegex.FindStringSubmatch(chatMessage)
	tagsArray := strings.Split(tagString[1], ";")
	tagsMap := make(map[string]string)

	for _, tag := range tagsArray {
		keyValue := strings.Split(tag, "=")
		tagsMap[keyValue[0]] = keyValue[1]
    }

	return tagsMap
}