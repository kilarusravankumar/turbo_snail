package track

import (
	"sync"
	"turbo_snail/broker/priority_queue"
	"turbo_snail/broker/message"
	"container/heap"
)

type Track struct {
	Name  string
	queue *priority_queue.MagicQueue
	// wal *os.File
	mu sync.Mutex
}

func New(trackName string) *Track {
	t := &Track{}
	t.Name = trackName
	t.queue = priority_queue.New()
	return t
}


func(t *Track) AppendMessage(msg *message.Message) {
	t.mu.Lock()
	defer t.mu.Unlock()
	heap.Push(t.queue, msg)
}


func (t *Track) PopMessage() *message.Message {
	t.mu.Lock()
    defer t.mu.Unlock()
	if t.queue.Len() < 1 {
		return nil
	}
	return heap.Pop(t.queue).(*message.Message)
}

func (t *Track) isEmpty() bool {
	return t.queue.Len() == 0
}

func (t *Track) Len() int {
	return t.queue.Len()
}
