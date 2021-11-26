package tcp

import (
	"sync/atomic"
	xtNet "xtnet/net"
)

type Client struct {
	addr    string
	closed  int32
	session Session
	agent   xtNet.IClientAgent
}

func NewClient(addr string, agent xtNet.IClientAgent) *Client {
	return &Client{
		addr:   addr,
		closed: 0,
		agent:  agent,
	}
}

func (client *Client) OnSessionStarted(session xtNet.ISession) {
	if atomic.LoadInt32(&client.closed) == 0 {
		client.agent.HandlerConnected(client)
	}
}

func (client *Client) Connect() bool {
	//conn, err := net.Dial("tcp", client.addr)

	return true
}

func (client *Client) ConnectSync() bool {
	return true
}

func (client *Client) Close() {

}
