package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/url"
	"sync/atomic"
	"time"
	xtnet "xtnet"
	xtnetNet "xtnet/net"
)

const (
	clientStatusClosed int32 = iota
	clientStatusConnecting
	clientStatusConnected
)

var (
	ClientErrWrongStatus = errors.New("wrong status")
	ClientErrTimeout     = errors.New("connect time out")
)

type Client struct {
	addr    string
	status  int32
	session xtnetNet.ISession
	agent   xtnetNet.IClientAgent
}

func NewClient(addr string, agent xtnetNet.IClientAgent) *Client {
	return &Client{
		addr:   addr,
		status: clientStatusClosed,
		agent:  agent,
	}
}

func (client *Client) GetSession() xtnetNet.ISession {
	return client.session
}

func (client *Client) Connect() bool {
	if !atomic.CompareAndSwapInt32(&client.status, clientStatusClosed, clientStatusConnecting) {
		xtnet.GetLogger().LogWarn("websocket.Client.Connect: client is not closed")
		return false
	}

	go func() {
		u, errParse := url.Parse(client.addr)
		if errParse != nil {
			xtnet.GetLogger().LogError("tcp.Client.Connect: url.Parse error=%s", errParse.Error())
			return
		}

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			xtnet.GetLogger().LogError("tcp.Client.Connect: connect error=%s", err.Error())
			atomic.StoreInt32(&client.status, clientStatusClosed)
			client.agent.HandleConnectFailed(client)
			return
		}

		session := NewSession(conn)
		session.SetSessionStartCb(func(session xtnetNet.ISession) {
			atomic.StoreInt32(&client.status, clientStatusConnected)
			client.agent.HandleConnect(client)
		})
		session.SetSessionDataCb(func(session xtnetNet.ISession, data []byte) {
			client.agent.HandleClientData(client, data)
		})
		session.SetSessionCloseCb(func(session xtnetNet.ISession) {
			atomic.StoreInt32(&client.status, clientStatusClosed)
			client.session = nil
			client.agent.HandleConnectBreak(client)
		})
		client.session = session
		session.start()
	}()

	return true
}

func (client *Client) ConnectSync(TimeOutMS int) error {
	if !atomic.CompareAndSwapInt32(&client.status, clientStatusClosed, clientStatusConnecting) {
		xtnet.GetLogger().LogWarn("tcp.Client.Connect: client is not closed")
		return ClientErrWrongStatus
	}

	u, errParse := url.Parse(client.addr)
	if errParse != nil {
		xtnet.GetLogger().LogError("tcp.Client.Connect: url.Parse error=%s", errParse.Error())
		return errParse
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		xtnet.GetLogger().LogWarn("tcp.Client.Connect: connect error=%s", err.Error())
		atomic.StoreInt32(&client.status, clientStatusClosed)
		client.agent.HandleConnectFailed(client)
		return err
	}

	chanSign := make(chan int)

	session := NewSession(conn)
	session.SetSessionStartCb(func(session xtnetNet.ISession) {
		atomic.StoreInt32(&client.status, clientStatusConnected)
		chanSign <- 1
	})
	session.SetSessionDataCb(func(session xtnetNet.ISession, data []byte) {
		client.agent.HandleClientData(client, data)
	})
	session.SetSessionCloseCb(func(session xtnetNet.ISession) {
		atomic.StoreInt32(&client.status, clientStatusClosed)
		client.session = nil
		client.agent.HandleConnectBreak(client)
	})
	client.session = session
	session.start()

	select {
	case <-chanSign:
		return nil
	case <-time.After(time.Millisecond * time.Duration(TimeOutMS)):
		client.Close(false)
		return ClientErrTimeout
	}
}

func (client *Client) Send(data []byte) {
	if atomic.LoadInt32(&client.status) == clientStatusConnected {
		client.session.Send(data)
	}
}

func (client *Client) Close(waitWrite bool) {
	if atomic.LoadInt32(&client.status) == clientStatusConnecting || atomic.LoadInt32(&client.status) == clientStatusConnected {
		atomic.StoreInt32(&client.status, clientStatusClosed)
		client.session.Close(waitWrite)
		client.session = nil
	}
}
