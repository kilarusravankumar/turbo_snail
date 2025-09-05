package track

import (
	"container/heap"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"
	"turbo_snail/log_entry"
	"turbo_snail/message"
	"turbo_snail/priority_queue"
)

type Track struct {
	Name       string
	queue      *priority_queue.MagicQueue
	wal        *os.File
	WalEncoder *gob.Encoder
	mu         sync.Mutex
}

func New(trackName string) *Track {
	t := &Track{}
	t.Name = trackName
	var err error

	// Ensure the WAL directory exists
	if err := os.MkdirAll("wal", 0755); err != nil {
		log.Fatalf("Failed to create WAL directory: %v", err)
	}

	t.wal, err = os.Create(fmt.Sprintf("wal/%s.log", trackName))
	if err != nil {
		log.Fatalf("Failed to create WAL file for track %s: %v", trackName, err)
	}
	t.WalEncoder = gob.NewEncoder(t.wal)
	t.queue = priority_queue.New()
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
	deleteEntry := log_entry.NewMsg(msg)

	if err := t.WalEncoder.Encode(deleteEntry); err != nil {
		log.Fatalf("Error occured while appending delete log the wal file.\n %s", err.Error())
		return nil
	}

	if err := t.wal.Sync(); err != nil {
		log.Fatalf("Error occured while appending delete log the wal file.\n %s", err.Error())
		return nil
	}
	return msg
}

func (t *Track) IsEmpty() bool {
	return t.queue.Len() == 0
}

func (t *Track) Len() int {
	return t.queue.Len()
}
