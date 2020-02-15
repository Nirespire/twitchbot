package util

import (
	"time"
)

const estFormat = "Jan 8 10:28:00 EST"

func timeStamp(format string) string {
	return time.Now().Format(format)
}

func TimeStamp() string {
	return timeStamp(estFormat)
}
