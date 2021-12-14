package client

import (
	"encoding/binary"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/eventhandler"
	"xtnet/net/packet"
)

type Normal struct {
	loop         *frame.Loop
	byteOrder    binary.ByteOrder
	eventHandler *eventhandler.Client
}

func NewNormal(loop *frame.Loop, byteOrder binary.ByteOrder) *Normal {
	return &Normal{
		loop:      loop,
		byteOrder: byteOrder,
	}
}

func (agent *Normal) SetEventHandler(eventHandler *eventhandler.Client) {
	agent.eventHandler = eventHandler
}

func (agent *Normal) HandleConnect(client net.IClient) {
	agent.loop.Post(func() {
		agent.eventHandler.OnConnectSuccess(client)
	})
}

func (agent *Normal) HandleConnectFailed(client net.IClient) {
	agent.loop.Post(func() {
		agent.eventHandler.OnConnectFailed(client)
	})
}

func (agent *Normal) HandleClientData(client net.IClient, data []byte) {
	rpk := packet.NewReadPacket(data, agent.byteOrder, 0, len(data))
	agent.loop.Post(func() {
		agent.eventHandler.OnClientPacket(client, rpk)
	})
}

func (agent *Normal) HandleConnectBreak(client net.IClient) {
	agent.loop.Post(func() {
		agent.eventHandler.OnConnectionBroken(client)
	})
}
