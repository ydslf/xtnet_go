package tcp

import (
	"net"
	myNet "xtnet/net"
)

type SessionOuterServer struct {
	Session
}

func (session *SessionOuterServer) start() {
	//TODO 发送key, 开接写协程写, 完了再开读协程
	session.sendKey()
	session.doStart()
}

func (session *SessionOuterServer) sendKey() {
	//TODO
}

type SessionOsCreator struct {
	sendChanSize int
}

func NewSessionOsCreator(sendChanSize int) ISessionCreator {
	return &SessionOsCreator{
		sendChanSize: sendChanSize,
	}
}

func (c *SessionOsCreator) CreateSession(netBase myNet.INetBase, conn net.Conn) ISession {
	return &SessionOuterServer{
		Session{
			netBase:   netBase,
			conn:      conn,
			closed:    0,
			sendChan:  make(chan []byte, c.sendChanSize),
			closeChan: make(chan int, 1),
		},
	}
}
