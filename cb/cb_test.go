package cb_test

import (
	"sync"
	"testing"

	"fortio.org/memstore/cb"
)

func TestBoundaryConditions(t *testing.T) {
	const capacity = 5
	buffer := cb.New[int](capacity)

	// Test empty buffer
	if !buffer.Empty() {
		t.Error("Buffer should be empty")
	}

	v, ok := buffer.Pop()
	if ok {
		t.Errorf("Popping from an empty buffer should fail, got %d", v)
	}

	// Test buffer size and capacity
	if buffer.Size() != 0 || buffer.Capacity() != capacity {
		t.Errorf("Buffer should have size 0 and capacity %d", capacity)
	}

	// Test pushing and popping single item
	buffer.Push(1)
	if buffer.Empty() {
		t.Error("Buffer should not be empty")
	}
	if buffer.Size() != 1 {
		t.Error("Buffer should have size 1")
	}

	item, ok := buffer.Pop()
	if !ok || item != 1 {
		t.Error("Popped item should be 1")
	}
	if !buffer.Empty() {
		t.Error("Buffer should be empty after popping the only item")
	}

	// Test filling the buffer
	for i := 0; i < capacity; i++ {
		if !buffer.Push(i) {
			t.Errorf("Failed to push item %d to the buffer", i)
		}
	}
	if !buffer.Full() {
		t.Error("Buffer should be full")
	}

	ok = buffer.Push(-1) // should fail
	if ok {
		t.Error("Pushing to a full buffer should fail")
	}

	// Test popping items from the full buffer
	for i := 0; i < capacity; i++ {
		item, ok := buffer.Pop()
		if !ok || item != i {
			t.Errorf("Popped item should be %d, got %d", i, item)
		}
	}

	// Test concurrent PushBlocking and PopBlocking
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			buffer.PushBlocking(i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			item := buffer.PopBlocking()
			if item != i {
				t.Errorf("Popped item should be %d, got %d", i, item)
			}
		}
	}()

	wg.Wait()
}

func BenchmarkCircularBuffer(b *testing.B) {
	c := cb.New[int](100)
	var x int
	var ok bool
	for i := 0; i < b.N; i++ {
		c.Push(i)
		x, ok = c.Pop()
		if !ok {
			b.Fatal("Dequeue failed")
		}
	}
	b.Logf("x=%d", x)
}

func TestProducerConsumerScenario(t *testing.T) {
	buffer := cb.New[int](10)

	var wg sync.WaitGroup
	wg.Add(11) // 10 producers + 1 consumer

	producer := func(id int) {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			buffer.PushBlocking(id*10 + i)
		}
	}

	consumer := func() {
		defer wg.Done()
		itemsDequeued := 0
		for itemsDequeued < 20 { // Total items to dequeue: 10 producers * 2 items each
			v := buffer.PopBlocking()
			t.Logf("Dequeued value: %d", v)
			itemsDequeued++
		}
	}

	// Start 10 producer goroutines
	for i := 0; i < 10; i++ {
		go producer(i)
	}

	// Start 1 consumer goroutine
	go consumer()

	wg.Wait()

	if !buffer.Empty() {
		t.Error("Buffer should be empty after all items are dequeued")
	}
}

func BenchmarkCircularBufferBlocking(b *testing.B) {
	c := cb.New[int](100)
	var x int
	for i := 0; i < b.N; i++ {
		c.PushBlocking(i)
		x = c.PopBlocking()
		if x != i {
			b.Fatalf("Wrong values %d %d", x, i)
		}
	}
	b.Logf("x=%d", x)
}

func benchmarkPushBlocking(b *testing.B, numProducers, numConsumers int) {
	buffer := cb.New[int](20) // small queue, higher contention

	var wg sync.WaitGroup
	wg.Add(numProducers + numConsumers)
	prodN := b.N * numConsumers
	consN := b.N * numProducers
	producer := func() {
		defer wg.Done()
		for i := 0; i < prodN; i++ {
			buffer.PushBlocking(i)
		}
	}

	consumer := func() {
		defer wg.Done()
		for i := 0; i < consN; i++ {
			buffer.PopBlocking()
		}
	}

	b.ResetTimer()

	// Start producer goroutines
	for i := 0; i < numProducers; i++ {
		go producer()
	}

	// Start consumer goroutines
	for i := 0; i < numConsumers; i++ {
		go consumer()
	}

	wg.Wait()
}

func BenchmarkPushBlocking(b *testing.B) {
	numProducers := 13
	numConsumers := 7

	b.Run("BenchmarkPushBlocking 7,5", func(b *testing.B) {
		benchmarkPushBlocking(b, numProducers, numConsumers)
	})
}
