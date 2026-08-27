package rpc

import (
	"encoding/binary"
	"sync"
	"testing"
	"xtnet/net/packet"
)

func TestSyncGenContextIDConcurrent(t *testing.T) {
	const (
		goroutines = 32
		perRoutine = 1000
	)

	rpc := &Sync{}
	ids := make(chan int32, goroutines*perRoutine)
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range perRoutine {
				ids <- rpc.GenContextID()
			}
		}()
	}
	wg.Wait()
	close(ids)

	seen := make(map[int32]struct{}, goroutines*perRoutine)
	for id := range ids {
		if _, ok := seen[id]; ok {
			t.Fatalf("GenContextID() generated duplicate ID %d", id)
		}
		seen[id] = struct{}{}
	}
}

func TestSyncGenContextIDWrap(t *testing.T) {
	rpc := &Sync{}
	rpc.contextID.Store(maxContextID - 1)

	if got := rpc.GenContextID(); got != maxContextID {
		t.Fatalf("first GenContextID() = %d, want %d", got, maxContextID)
	}
	if got := rpc.GenContextID(); got != 1 {
		t.Fatalf("wrapped GenContextID() = %d, want 1", got)
	}
}

func TestSyncHandlerResponseRoutesByContextID(t *testing.T) {
	rpc := &Sync{}
	first := &ContextSync{
		contextID: 1,
		t:         rqtSync,
		response:  make(chan *packet.ReadPacket, 1),
	}
	second := &ContextSync{
		contextID: 2,
		t:         rqtSync,
		response:  make(chan *packet.ReadPacket, 1),
	}
	rpc.contexts.Store(first.contextID, first)
	rpc.contexts.Store(second.contextID, second)

	want := packet.NewReadPacket([]byte{2}, binary.BigEndian, 0, 1)
	rpc.handlerResponse(second.contextID, want)

	select {
	case got := <-second.response:
		if got != want {
			t.Fatalf("handlerResponse() response = %p, want %p", got, want)
		}
	default:
		t.Fatal("handlerResponse() did not route response to matching context")
	}

	select {
	case <-first.response:
		t.Fatal("handlerResponse() routed response to the wrong context")
	default:
	}

	if _, ok := rpc.contexts.Load(second.contextID); ok {
		t.Fatal("handlerResponse() did not remove completed context")
	}
	if _, ok := rpc.contexts.Load(first.contextID); !ok {
		t.Fatal("handlerResponse() removed unrelated context")
	}
}
