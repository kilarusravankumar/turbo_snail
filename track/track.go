package track

import (
	"container/heap"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"turbo_snail/config"
	"turbo_snail/log_entry"
	"turbo_snail/message"
	"turbo_snail/priority_queue"

	"github.com/google/uuid"
)

type Track struct {
	Name             string
	queue            *priority_queue.MagicQueue
	wal              *os.File
	WalEncoder       *gob.Encoder
	mu               sync.Mutex
	InFlightMessages map[uuid.UUID]*message.InFlightMessage
}

func New(trackName string) *Track {
	t := &Track{}
	t.Name = trackName
	var err error

	// Ensure the WAL directory exists
	if err := os.MkdirAll(config.WAL_DIR, 0755); err != nil {
		log.Printf("Ensuring %s , directory exists for writing Write ahead logs: %v", config.WAL_DIR, err)
	}

	t.wal, err = os.Create(fmt.Sprintf("%s/%s.log", config.WAL_DIR, trackName))
	if err != nil {
		log.Fatalf("Failed to create WAL file for track %s: %v", trackName, err)
	}
	t.WalEncoder = gob.NewEncoder(t.wal)
	t.queue = priority_queue.New()
	t.InFlightMessages = make(map[uuid.UUID]*message.InFlightMessage, 0)
	go t.RequeueExpiredMessages()
	return t
}

func (t *Track) AddMessage(msg *message.Message) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	entry := log_entry.NewMsg(msg)

	if err := t.WalEncoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode log entry to WAL for track %s: %w", t.Name, err)
	}

	if err := t.wal.Sync(); err != nil {
		return fmt.Errorf("failed to sync WAL to disk for track %s: %w", t.Name, err)
	}

	heap.Push(t.queue, msg)

	return nil
}

func (t *Track) PopMessage() *message.Message {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.queue.Len() < 1 {
		return nil
	}
	msg := heap.Pop(t.queue).(*message.Message)

	return msg
}

func (t *Track) IsEmpty() bool {
	return t.queue.Len() == 0
}

func (t *Track) Len() int {
	return t.queue.Len()
}

func (t *Track) RequeueExpiredMessages() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for range ticker.C {
		for _uuid, inFlightMsg := range t.InFlightMessages {
			if time.Now().After(inFlightMsg.ExpiresAt) {
				err := t.AddMessage(inFlightMsg.Message)
				if err != nil {
					log.Default().Fatalf("Error occured while requeuing the expired Message \n %s", err.Error())
				}
				delete(t.InFlightMessages, _uuid)
			}
		}
	}

}
