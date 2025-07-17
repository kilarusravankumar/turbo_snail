package priority_queue


import (
	"container/heap"
	"turbo_snail/broker/message"
)

type MagicQueue []*message.Message


func (q MagicQueue) Len() int {
	return len(q)
}


func (q MagicQueue) Swap(i, j int) {
	q[i] , q[j] = q[j], q[i]
}

func (q MagicQueue) Less(i,j int) bool {
	if q[i].Priority != q[j].Priority {
		return q[i].Priority > q[j].Priority
	}
	return q[i].Timestamp > q[j].Timestamp
}


func (q *MagicQueue) Push(x any) {
	msg := x.(*message.Message)
	*q = append(*q, msg)
}

func (q *MagicQueue) Pop() any {
	oldQ := *q
	oldQ_len := len(oldQ)
	top := oldQ[oldQ_len - 1]
	oldQ = oldQ[:oldQ_len - 1]
	*q = oldQ
	return top
}

func New() *MagicQueue {
	q := &MagicQueue{}
	heap.Init(q)
	return q
}
