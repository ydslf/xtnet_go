package frame

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestLoopQueueGrows(t *testing.T) {
	loop := NewLoop(LoopSizeMin, false)
	var handled atomic.Int32

	for range LoopSizeMin + 1 {
		loop.Post(func() {
			handled.Add(1)
		})
	}
	loop.Post(func() {
		loop.Close(true)
	})

	loop.Run()
	if got := handled.Load(); got != int32(LoopSizeMin+1) {
		t.Fatalf("handled %d functions, want %d", got, LoopSizeMin+1)
	}
}

func TestLoopCloseWaitHandle(t *testing.T) {
	loop := NewLoop(LoopSizeMin, false)
	var handled atomic.Int32

	loop.Post(func() {
		handled.Add(1)
		loop.Close(true)
	})
	loop.Post(func() {
		handled.Add(1)
	})

	loop.Run()
	if got := handled.Load(); got != 2 {
		t.Fatalf("handled %d functions, want 2", got)
	}
}

func TestLoopCloseWithoutWait(t *testing.T) {
	loop := NewLoop(LoopSizeMin, false)
	var handled atomic.Int32

	loop.Post(func() {
		handled.Add(1)
		loop.Close(false)
	})
	loop.Post(func() {
		handled.Add(1)
	})

	loop.Run()
	if got := handled.Load(); got != 1 {
		t.Fatalf("handled %d functions, want 1", got)
	}
}

func TestLoopPostDoesNotOverwriteRunningBatch(t *testing.T) {
	const batchSize = 128

	loop := NewLoop(LoopSizeMin, false)
	handled := make([]int, batchSize*2)
	batchStarted := make(chan struct{})
	continueBatch := make(chan struct{})
	runDone := make(chan struct{})

	loop.Post(func() {
		handled[0]++
		close(batchStarted)
		<-continueBatch
	})
	for i := 1; i < batchSize; i++ {
		index := i
		loop.Post(func() {
			handled[index]++
		})
	}

	go func() {
		loop.Run()
		close(runDone)
	}()

	select {
	case <-batchStarted:
	case <-time.After(time.Second):
		t.Fatal("loop did not start the first batch")
	}

	for i := batchSize; i < batchSize*2; i++ {
		index := i
		loop.Post(func() {
			handled[index]++
		})
	}
	loop.Post(func() {
		loop.Close(true)
	})
	close(continueBatch)

	select {
	case <-runDone:
	case <-time.After(time.Second):
		t.Fatal("loop did not finish")
	}

	for i, count := range handled {
		if count != 1 {
			t.Fatalf("function %d handled %d times, want 1", i, count)
		}
	}
}

func TestLoopTakeAllPreservesCapacity(t *testing.T) {
	loop := NewLoop(LoopSizeMin, false)
	initialCap := cap(loop.functions)

	loop.Post(func() {})
	loop.Post(func() {})
	loop.RunOnce()
	functions := loop.takeAll()
	clear(functions)

	if got := cap(loop.functions) + cap(loop.running); got < initialCap {
		t.Fatalf("queue capacity shrank to %d, want at least %d", got, initialCap)
	}
}
