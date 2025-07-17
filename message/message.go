package message

import (
	"github.com/google/uuid"
	"time"
)


type Message struct {
	ID uuid.UUID
	Data []byte
	Timestamp int64 
	Priority int8
}


func New(data []byte, priority int8) *Message {
	return &Message{
		ID: uuid.New(),
		Data: data,
		Timestamp: time.Now().Unix(),
		Priority: priority, 
	}
}

