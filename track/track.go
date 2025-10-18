package track

import (
	"container/heap"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
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
	AckLog           *os.File
	AckEncoder       *gob.Encoder
}

func New(trackName, walDir string) *Track {
	t := &Track{}
	t.Name = trackName
	var err error

	// Ensure the WAL directory exists
	if err := os.MkdirAll(walDir, 0o755); err != nil {
		log.Printf("Ensuring %s , directory exists for writing Write ahead logs: %v", walDir, err)
	}

	t.wal, err = os.Create(fmt.Sprintf("%s/%s.log", walDir, trackName))
	if err != nil {
		log.Fatalf("Failed to create WAL file for track %s: %v", trackName, err)
	}

	ackFullFilePath := fmt.Sprintf("%s/ack/%s.ack.log", walDir, trackName)
	ackDirpath := filepath.Dir(ackFullFilePath)
	err = os.MkdirAll(ackDirpath, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create ack directory in %s \n", ackDirpath)
	}

	t.AckLog, err = os.Create(ackFullFilePath)
	if err != nil {
		log.Fatalf("Failed to create ACK log file, which will create persistent and replay issues after ACK; \n %s", err.Error())
	}
	t.WalEncoder = gob.NewEncoder(t.wal)
	t.AckEncoder = gob.NewEncoder(t.AckLog)
	t.queue = priority_queue.New()
	t.InFlightMessages = make(map[uuid.UUID]*message.InFlightMessage, 0)

	// requeuing expired messages
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
	t.InFlightMessages[msg.ID] = &message.InFlightMessage{Message: msg, ExpiresAt: time.Now().Add(30 * time.Minute)}
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
		t.mu.Lock()
		for _uuid, inFlightMsg := range t.InFlightMessages {
			if time.Now().After(inFlightMsg.ExpiresAt) {
				err := t.AddMessage(inFlightMsg.Message)
				if err != nil {
					log.Default().Fatalf("Error occured while requeuing the expired Message \n %s", err.Error())
				}
				delete(t.InFlightMessages, _uuid)
			}
		}
		t.mu.Unlock()
	}
}

func (t *Track) ACKMessage(msgID uuid.UUID) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.InFlightMessages[msgID]; !ok {
		return errors.New("uuid provided is not in inflight messages; msg is either requeued or msg id is invalid")
	}
	entry := log_entry.NewACKLogMsg(msgID)
	if err := t.AckEncoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode log entry to ACK for track %s: %w", t.Name, err)
	}

	if err := t.AckLog.Sync(); err != nil {
		return fmt.Errorf("failed to sync ACK Log to disk for track %s: %w", t.Name, err)
	}
	delete(t.InFlightMessages, msgID)
	return nil
}

func (t *Track) NACKMessage(msgID uuid.UUID) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	inflightMsg, ok := t.InFlightMessages[msgID]
	if !ok {
		return fmt.Errorf("msg not found for the msg Id : %s", msgID)
	}
	t.AddMessage(inflightMsg.Message)
	delete(t.InFlightMessages, msgID)
	return nil
}
