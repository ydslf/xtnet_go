package tcp

import (
	"errors"
	"net"
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
	addr           string
	status         int32
	session        xtnetNet.ISession
	sessionCreator ISessionCreator
	pktProcCreator IPktProcCreator
	agent          xtnetNet.IClientAgent
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
		xtnet.GetLogger().LogWarn("tcp.Client.Connect: client is not closed")
		return false
	}

	go func() {
		conn, err := net.Dial("tcp", client.addr)
		if err != nil {
			atomic.StoreInt32(&client.status, clientStatusClosed)
			client.agent.HandleConnectFailed(client)
			return
		}

		if client.sessionCreator == nil {
			client.sessionCreator = NewSessionCreator(sendChanSizeDefault)
		}
		if client.pktProcCreator == nil {
			client.pktProcCreator = NewPktProcCreator(maxPkgSizeDefault, orderDefault)
		}

		session := client.sessionCreator.CreateSession(conn)
		pktProc := client.pktProcCreator.CreatePktProc()
		session.setPktProc(pktProc)
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

	conn, err := net.Dial("tcp", client.addr)
	if err != nil {
		atomic.StoreInt32(&client.status, clientStatusClosed)
		return err
	}

	chanSign := make(chan int)

	if client.sessionCreator == nil {
		client.sessionCreator = NewSessionCreator(sendChanSizeDefault)
	}
	if client.pktProcCreator == nil {
		client.pktProcCreator = NewPktProcCreator(maxPkgSizeDefault, orderDefault)
	}

	session := client.sessionCreator.CreateSession(conn)
	pktProc := client.pktProcCreator.CreatePktProc()
	session.setPktProc(pktProc)
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
	session.start()
	client.session = session

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
