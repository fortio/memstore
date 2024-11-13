package cb_test

import (
	"sync"
	"testing"

	"fortio.org/memstore/cb"
)

const capacity = 5

func TestBoundaryConditions(t *testing.T) {
	testCases := []struct {
		name string
		cb   cb.Queue[int]
	}{
		{name: "CircBuffer", cb: cb.New[int](capacity)},
		{name: "Channel", cb: cb.NewC[int](capacity)},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			testBoundaryConditions(tt, tc.cb)
		})
	}
}

func testBoundaryConditions(t *testing.T, buffer cb.Queue[int]) {
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
	for i := range capacity {
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
	for i := range capacity {
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
		for i := range 100 {
			buffer.PushBlocking(i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := range 100 {
			item := buffer.PopBlocking()
			if item != i {
				t.Errorf("Popped item should be %d, got %d", i, item)
			}
		}
	}()

	wg.Wait()
}

func BenchmarkCircularBuffer(b *testing.B) {
	capacity := 100
	testCases := []struct {
		name string
		cb   cb.Queue[int]
	}{
		{name: "CircBuffer", cb: cb.New[int](capacity)},
		{name: "Channel", cb: cb.NewC[int](capacity)},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(bb *testing.B) {
			benchmarkCircularBuffer(bb, tc.cb)
		})
	}
}

func benchmarkCircularBuffer(b *testing.B, c cb.Queue[int]) {
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
	testCases := []struct {
		name string
		cb   cb.Queue[int]
	}{
		{name: "CircBuffer", cb: cb.New[int](10)},
		{name: "Channel", cb: cb.NewC[int](10)},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			testProducerConsumerScenario(tt, tc.cb)
		})
	}
}

func testProducerConsumerScenario(t *testing.T, buffer cb.Queue[int]) {
	var wg sync.WaitGroup
	wg.Add(11) // 10 producers + 1 consumer

	producer := func(id int) {
		defer wg.Done()
		for i := range 2 {
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
	for i := range 10 {
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
	capacity := 100
	testCases := []struct {
		name string
		cb   cb.Queue[int]
	}{
		{name: "CircBuffer", cb: cb.New[int](capacity)},
		{name: "Channel", cb: cb.NewC[int](capacity)},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(bb *testing.B) {
			benchmarkCircularBufferBlocking(bb, tc.cb)
		})
	}
}

func benchmarkCircularBufferBlocking(b *testing.B, c cb.Queue[int]) {
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

func benchmarkPushBlocking(b *testing.B, buffer cb.Queue[int], numProducers, numConsumers int) {
	var wg sync.WaitGroup
	wg.Add(numProducers + numConsumers)
	prodN := b.N * numConsumers
	consN := b.N * numProducers
	producer := func() {
		defer wg.Done()
		for i := range prodN {
			buffer.PushBlocking(i)
		}
	}

	consumer := func() {
		defer wg.Done()
		for range consN {
			buffer.PopBlocking()
		}
	}

	b.ResetTimer()

	// Start producer goroutines
	for range numProducers {
		go producer()
	}

	// Start consumer goroutines
	for range numConsumers {
		go consumer()
	}

	wg.Wait()
}

func BenchmarkHighContention(b *testing.B) {
	numProducers := 13
	numConsumers := 7

	testCases := []struct {
		name string
		cb   cb.Queue[int]
	}{
		{name: "CircBuffer", cb: cb.New[int](capacity)},
		{name: "Channel", cb: cb.NewC[int](capacity)},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(bb *testing.B) {
			benchmarkPushBlocking(bb, tc.cb, numProducers, numConsumers)
		})
	}
}
