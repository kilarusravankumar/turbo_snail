package message

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	Data      map[string]interface{}
	Timestamp int64
	Priority  int8
}

func New(data map[string]interface{}, priority int8) *Message {
	return &Message{
		ID:        uuid.New(),
		Data:      data,
		Timestamp: time.Now().Unix(),
		Priority:  priority,
	}
}
