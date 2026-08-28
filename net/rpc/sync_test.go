package rpc

import (
	"encoding/binary"
	"errors"
	"sync"
	"testing"
	"xtnet/frame"
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

func TestSyncHandlerAsyncResponseCallback(t *testing.T) {
	loop := frame.NewLoop(0, false)
	rpc := &Sync{loop: loop}
	want := packet.NewReadPacket([]byte{1}, binary.BigEndian, 0, 1)
	var gotPacket *packet.ReadPacket
	var gotErr error
	context := &ContextSync{
		contextID: 1,
		t:         rqtAsync,
		cb: func(rpk *packet.ReadPacket, err error) {
			gotPacket = rpk
			gotErr = err
		},
	}
	rpc.contexts.Store(context.contextID, context)

	rpc.handlerResponse(context.contextID, want)
	loop.RunOnce()

	if gotPacket != want {
		t.Fatalf("callback packet = %p, want %p", gotPacket, want)
	}
	if gotErr != nil {
		t.Fatalf("callback error = %v, want nil", gotErr)
	}
}

func TestSyncHandlerAsyncTimeoutCallbackOnce(t *testing.T) {
	rpc := &Sync{}
	callbackCount := 0
	var gotPacket *packet.ReadPacket
	var gotErr error
	context := &ContextSync{
		contextID: 1,
		t:         rqtAsync,
		cb: func(rpk *packet.ReadPacket, err error) {
			callbackCount++
			gotPacket = rpk
			gotErr = err
		},
	}
	rpc.contexts.Store(context.contextID, context)

	rpc.handlerAsyncTimeout(context.contextID, context)
	rpc.handlerAsyncTimeout(context.contextID, context)

	if callbackCount != 1 {
		t.Fatalf("timeout callback count = %d, want 1", callbackCount)
	}
	if gotPacket != nil {
		t.Fatalf("timeout callback packet = %p, want nil", gotPacket)
	}
	if !errors.Is(gotErr, ErrRequestTimeout) {
		t.Fatalf("timeout callback error = %v, want %v", gotErr, ErrRequestTimeout)
	}
}
