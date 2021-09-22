package tcp

import (
	"net"
	"sync/atomic"
	myNet "xtnet/net"
)

type SessionEncryptClient struct {
	Session
}

func (session *SessionEncryptClient) start() {
	//TODO 接收key, 开接读程接收, 完了再开写协程
	if atomic.CompareAndSwapInt32(&session.status, statusNone, statusInit) {
		session.receiveKey()
		session.doStart()
	}
}

func (session *SessionEncryptClient) receiveKey() {
	//TODO
}

type SessionEncryptCCreator struct {
	sendChanSize int
}

func NewSessionEncryptCreator(sendChanSize int) ISessionCreator {
	return &SessionEncryptCCreator{
		sendChanSize: sendChanSize,
	}
}

func (c *SessionEncryptCCreator) CreateSession(netBase myNet.INetBase, conn net.Conn) ISession {
	return &SessionEncryptClient{
		Session{
			netBase:   netBase,
			conn:      conn,
			status:    statusNone,
			sendChan:  make(chan []byte, c.sendChanSize),
			closeChan: make(chan int, 1),
		},
	}
}
