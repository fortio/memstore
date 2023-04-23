// First In First Out (FIFO) [Queue] with fixed capacity.
// Circular Buffer implementation in Go with both
// pub/sub thread safe blocking API and pure FIFO queue with set capacity
// unsynchronized base. Two versions of the same [Queue] interface:
// one actually a circular buffer [CircularBuffer], the other using a channel
// [CircularBufferChan], created using [cb.New] or [cb.NewC] respectively.
package cb

import "sync"

type Queue[T any] interface {
	Empty() bool
	Full() bool
	Size() int
	Capacity() int
	Push(item T) bool
	Pop() (value T, ok bool)
	PushBlocking(item T)
	PopBlocking() (value T)
}

// FIFO [Queue] with fixed capacity. Fixed array implementation.
type CircularBuffer[T any] struct {
	buffer []T
	head   int
	tail   int
	size   int
	mu     sync.Mutex
	full   sync.Cond
	empty  sync.Cond
}

// NewC returns the fixed array version of 0 alloc fixed capacity (optionally blocking) [Queue].
func New[T any](capacity int) *CircularBuffer[T] {
	cb := &CircularBuffer[T]{
		buffer: make([]T, capacity),
		head:   0,
		tail:   0,
		size:   0,
	}
	cb.full.L = &cb.mu
	cb.empty.L = &cb.mu
	return cb
}

func (cb *CircularBuffer[T]) Empty() bool {
	return cb.size == 0
}

func (cb *CircularBuffer[T]) Full() bool {
	return cb.size == len(cb.buffer)
}

func (cb *CircularBuffer[T]) Size() int {
	return cb.size
}

func (cb *CircularBuffer[T]) Capacity() int {
	return len(cb.buffer)
}

// Push adds an item to the queue. returns false if queue is full.
func (cb *CircularBuffer[T]) Push(item T) bool {
	if cb.Full() {
		return false
	}
	cb.buffer[cb.tail] = item
	cb.tail++
	if cb.tail == len(cb.buffer) { // variable somehow not faster than len(buffer)
		cb.tail = 0
	}
	// cb.tail = (cb.tail + 1) % cb.capacity // classic but slower
	cb.size++
	return true
}

// Pop removes an item from the queue. returns false if queue is empty.
func (cb *CircularBuffer[T]) Pop() (T, bool) {
	var zero T
	if cb.Empty() {
		return zero, false
	}
	item := cb.buffer[cb.head]
	// cb.head = (cb.head + 1) % cb.capacity // slower
	cb.head++
	if cb.head == len(cb.buffer) { // not slower than cb.capacity variable
		cb.head = 0
	}
	cb.size--
	return item, true
}

// Thread safe blocking versions:

// Push adds an item to the queue. blocks if queue is full.
func (cb *CircularBuffer[T]) PushBlocking(item T) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	for cb.Full() {
		cb.full.Wait()
	}
	cb.Push(item)
	/*
		if cb.Size() == 1 {
			cb.empty.Broadcast()
		}
	*/
	// or:
	cb.empty.Signal()
}

// Pop removes an item from the queue. blocks if queue is empty.
func (cb *CircularBuffer[T]) PopBlocking() T {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	for cb.Empty() {
		cb.empty.Wait()
	}
	item, _ := cb.Pop()
	/*
		if cb.Size() == len(cb.buffer)-1 {
			cb.full.Broadcast()
		}
	*/
	// or:
	cb.full.Signal()
	return item
}
