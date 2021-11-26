package eventhandler

import (
	"xtnet/net"
	"xtnet/net/packet"
)

type OnConnectSuccess func(net.IClient)
type OnConnectFailed func(net.IClient)
type OnClientPacket func(net.IClient, *packet.ReadPacket)
type OnConnectionBroken func(net.IClient)

type Client struct {
	OnConnectSuccess   OnConnectSuccess
	OnConnectFailed    OnConnectFailed
	OnClientPacket     OnClientPacket
	OnConnectionBroken OnConnectionBroken
}

func NewClientEventHandler() *Client {
	return &Client{}
}
