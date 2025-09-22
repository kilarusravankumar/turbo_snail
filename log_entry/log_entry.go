package log_entry

import (
	"turbo_snail/message"
)

type LogEntry struct {
	Message *message.Message
}

func NewMsg(msg *message.Message) *LogEntry {

	return &LogEntry{
		Message: msg,
	}
}
