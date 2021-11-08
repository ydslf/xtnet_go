package tcp

import (
	"sync/atomic"
	xtNet "xtnet/net"
)

type Client struct {
	closed int64
	agent  xtNet.IAgent
}

func (client *Client) OnSessionStarted(session xtNet.ISession) {
	if atomic.LoadInt64(&client.closed) == 0 {
		client.agent.HandlerConnected(session)
	}
}

func (client *Client) Start() {

}
