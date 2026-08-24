package rpc

import (
	"encoding/binary"
	"testing"
	"xtnet/net/packet"
)

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
