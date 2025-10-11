// Package log_entry all wal log and ack log entry structs are stored here.
package log_entry

import (
	"github.com/google/uuid"
)

type ACKLogEntry struct {
	MsgID uuid.UUID
}

func NewACKLogMsg(msgID uuid.UUID) *ACKLogEntry {
	return &ACKLogEntry{
		MsgID: msgID,
	}
}
