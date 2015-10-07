package queue

import (
	"errors"
	"sync"
)

var (
	Empty = errors.New("Queue is empty")
)

type (
	fifo struct {
		lock sync.Mutex

		storage []string
	}
)

// Create a new fifo
func NewFifo() *fifo {
	queue := &fifo{
		storage: make([]string, 0),
		lock:    sync.Mutex{},
	}
	return queue
}

// Add an element to the back fo the fifo
func (queue *fifo) Push(item string) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	queue.storage = append(queue.storage, item)
}

// Remove an element from the fifo and return it
func (queue *fifo) Pop() (string, error) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	if len(queue.storage) == 0 {
		return "", Empty
	}
	value := queue.storage[0]
	queue.storage = queue.storage[1:]
	return value, nil
}

func (queue fifo) Length() int {
	return len(queue.storage)
}
