package message

import (
	"time"
)

type InFlightMessage struct {
	Message   *Message
	ExpiresAt time.Time
}
