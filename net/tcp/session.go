package tcp

import (
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
	xtnet_go "xtnet"
	mynet "xtnet/net"
)

type Session struct {
	conn           net.Conn
	pktProcessor   IPktProcessor
	closed         int64
	wgClose        sync.WaitGroup
	sendChan       chan []byte
	userData       interface{}
	onSessionData  mynet.OnSessionData
	onSessionClose mynet.OnSessionClose
	agent          mynet.IAgent
}

func newSession(conn net.Conn, pktProcessor IPktProcessor, sendBuffSize int) *Session {
	session := &Session{
		conn:         conn,
		pktProcessor: pktProcessor,
		closed:       0,
		sendChan:     make(chan []byte, sendBuffSize),
	}

	return session
}

func (session *Session) setCallback(onSessionData mynet.OnSessionData, onSessionClose mynet.OnSessionClose) {
	session.onSessionData = onSessionData
	session.onSessionClose = onSessionClose
}

func (session *Session) setAgent(agent mynet.IAgent) {
	session.agent = agent
}

func (session *Session) Send(data []byte) {
	if atomic.LoadInt64(&session.closed) == 1 {
		xtnet_go.GetLogger().LogWarn("tcp.Session.Send: session is closed")
		return
	}

	if len(session.sendChan) == cap(session.sendChan) {
		xtnet_go.GetLogger().LogWarn("tcp.Session.Send: sendChan is full")
	}

	session.sendChan <- data
}

//关闭session
func (session *Session) Close() {
	session.doClose(true)
}

//关闭session，等待关闭完成
func (session *Session) CloseWait() {
	if session.doClose(true) {
		session.wgClose.Wait()
	}
}

func (session *Session) doClose(active bool) bool {
	if !atomic.CompareAndSwapInt64(&session.closed, 0, 1) {
		if active {
			xtnet_go.GetLogger().LogWarn("tcp.Session.Close: session has been closed")
		}
		return false
	}

	if active {
		session.conn.SetReadDeadline(time.Now())
		close(session.sendChan)
	} else {
		session.conn.Close()
		close(session.sendChan)

		if session.agent != nil {
			session.agent.HandlerSessionClose(session)
		} else {
			session.onSessionClose(session)
		}
	}

	return true
}

func (session *Session) SetUserData(userData interface{}) {
	session.userData = userData
}

func (session *Session) GetUserData() interface{} {
	return session.userData
}

func (session *Session) start() {
	session.wgClose.Add(2)
	go session.readRoutine()
	go session.writeRoutine()
}

func (session *Session) readRoutine() {
	defer session.wgClose.Done()

	for {
		data, err := session.pktProcessor.UnPack(session)
		if err != nil {
			if err != io.EOF && !os.IsTimeout(err) {
				xtnet_go.GetLogger().LogError("session.readRoutine: err=%v", err)
			}
			session.doClose(false)
			return
		}

		if session.agent != nil {
			session.agent.HandlerSessionData(session, data)
		} else {
			session.onSessionData(session, data)
		}
	}
}

func (session *Session) writeRoutine() {
	defer session.wgClose.Done()

	for data := range session.sendChan {
		pktData := session.pktProcessor.Pack(data)
		_, err := session.conn.Write(pktData)
		if err != nil {
			xtnet_go.GetLogger().LogError("session.writeRoutine: err=%v", err)
			session.doClose(false)
			return
		}
	}

	session.conn.Close()
}
