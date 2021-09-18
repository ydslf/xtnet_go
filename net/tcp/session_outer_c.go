package tcp

import (
	"net"
	myNet "xtnet/net"
)

type SessionOuterClient struct {
	Session
}

func (session *SessionOuterClient) start() {
	//TODO 接收key, 开接读程接收, 完了再开写协程
	session.receiveKey()
	session.doStart()
}

func (session *SessionOuterClient) receiveKey() {
	//TODO
}

type SessionOcCreator struct {
	sendChanSize int
}

func NewSessionOcCreator(sendChanSize int) ISessionCreator {
	return &SessionOcCreator{
		sendChanSize: sendChanSize,
	}
}

func (c *SessionOcCreator) CreateSession(netBase myNet.INetBase, conn net.Conn) ISession {
	return &SessionOuterClient{
		Session{
			netBase:   netBase,
			conn:      conn,
			closed:    0,
			sendChan:  make(chan []byte, c.sendChanSize),
			closeChan: make(chan int, 1),
		},
	}
}
