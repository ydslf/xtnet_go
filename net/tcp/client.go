package tcp

import (
	"sync/atomic"
	myNet "xtnet/net"
)

type Client struct {
	closed int64
	agent  myNet.IAgent
}

func (client *Client) OnSessionStarted(session myNet.ISession) {
	if atomic.LoadInt64(&client.closed) == 0 {
		client.agent.HandlerConnected(session)
	}
}
