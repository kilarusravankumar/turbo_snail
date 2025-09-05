package log_entry

import (
	"time"
	"turbo_snail/message"
)

type LogEntry struct {
	deletedAt time.Time
	Message   *message.Message
}

func NewMsg(msg *message.Message) *LogEntry {

	return &LogEntry{
		Message: msg,
	}
}

func DeleteMsg(msg *message.Message) *LogEntry {
	return &LogEntry{
		deletedAt: time.Now(),
		Message:   msg,
	}
}

func (entry LogEntry) IsDeleted() bool {
	return !entry.deletedAt.IsZero()
}
