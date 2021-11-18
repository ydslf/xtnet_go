package tcp

import (
	"sync/atomic"
	xtNet "xtnet/net"
)

type Client struct {
	ip      string
	port    int
	closed  int32
	session Session
	agent   xtNet.IAgent
}

func (client *Client) OnSessionStarted(session xtNet.ISession) {
	if atomic.LoadInt32(&client.closed) == 0 {
		client.agent.HandlerConnected(session)
	}
}

func (client *Client) Start() bool {
	return true
}

func (client *Client) Close() {

}
