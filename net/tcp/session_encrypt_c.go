package tcp

import (
	"net"
	"sync/atomic"
)

type SessionEncryptClient struct {
	Session
}

func (session *SessionEncryptClient) start() {
	//TODO 接收key, 开接读接收, 完了再开写协程
	if atomic.CompareAndSwapInt32(&session.status, sessionStatusNone, sessionStatusInit) {
		session.receiveKey()
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

func (c *SessionEncryptCCreator) CreateSession(conn net.Conn) ISession {
	return &SessionEncryptClient{
		Session{
			conn:      conn,
			status:    sessionStatusNone,
			sendChan:  make(chan []byte, c.sendChanSize),
			closeChan: make(chan int, 1),
		},
	}
}
