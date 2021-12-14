package client

import (
	"encoding/binary"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/eventhandler"
	"xtnet/net/packet"
	"xtnet/net/rpc"
)

type Internal struct {
	loop         *frame.Loop
	byteOrder    binary.ByteOrder
	eventHandler *eventhandler.Client
	netRpc       rpc.IRpc
}

func NewInternal(loop *frame.Loop, byteOrder binary.ByteOrder) *Internal {
	return &Internal{
		loop:      loop,
		byteOrder: byteOrder,
	}
}

func (agent *Internal) SetEventHandler(eventHandler *eventhandler.Client) {
	agent.eventHandler = eventHandler
}

func (agent *Internal) SetNetRpc(netRpc rpc.IRpc) {
	agent.netRpc = netRpc
}

func (agent *Internal) HandleConnect(client net.IClient) {
	agent.loop.Post(func() {
		agent.eventHandler.OnConnectSuccess(client)
	})
}

func (agent *Internal) HandleConnectFailed(client net.IClient) {
	agent.loop.Post(func() {
		agent.eventHandler.OnConnectFailed(client)
	})
}

func (agent *Internal) HandleClientData(client net.IClient, data []byte) {
	rpk := packet.NewReadPacket(data, agent.byteOrder, 0, len(data))
	agent.netRpc.HandleSessionPacket(client.GetSession(), rpk)
}

func (agent *Internal) HandleConnectBreak(client net.IClient) {
	agent.loop.Post(func() {
		agent.eventHandler.OnConnectionBroken(client)
	})
}
