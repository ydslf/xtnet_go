package tcp

import (
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
	xtnet_go "xtnet"
	myNet "xtnet/net"
)

type Session struct {
	netBase      myNet.INetBase
	conn         net.Conn
	pktProcessor IPktProcessor
	closed       int64
	wgClose      sync.WaitGroup
	sendChan     chan []byte
	closeChan    chan int
	userData     interface{}
	agent        myNet.IAgent
}

func newSession(netBase myNet.INetBase, conn net.Conn, pktProcessor IPktProcessor, sendBuffSize int) *Session {
	session := &Session{
		netBase:      netBase,
		conn:         conn,
		pktProcessor: pktProcessor,
		closed:       0,
		sendChan:     make(chan []byte, sendBuffSize),
		closeChan:    make(chan int, 1),
	}

	return session
}

func (session *Session) setAgent(agent myNet.IAgent) {
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

type closeType uint

const (
	active  closeType = 0
	byRead  closeType = 1
	byWrite closeType = 2
)

//关闭session
func (session *Session) Close(waitWrite bool) {
	session.doClose(active, waitWrite)
}

//关闭session，等待关闭完成
func (session *Session) CloseBlock(waitWrite bool) {
	if session.doClose(active, waitWrite) {
		session.wgClose.Wait()
	}
}

func (session *Session) doClose(ct closeType, waitWrite bool) bool {
	if !atomic.CompareAndSwapInt64(&session.closed, 0, 1) {
		if ct == active {
			xtnet_go.GetLogger().LogWarn("tcp.Session.Close: session has been closed")
		}
		return false
	}

	if ct == active {
		session.conn.SetReadDeadline(time.Now())
		if waitWrite {
			close(session.sendChan)
		} else {
			session.closeChan <- 1
		}
	} else {
		if ct == byRead {
			session.closeChan <- 1
		} else {
			session.conn.Close()
		}

		session.agent.HandlerSessionClose(session)
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
	session.doStart()
}

func (session *Session) doStart() {
	session.wgClose.Add(2)
	go session.readRoutine()
	go session.writeRoutine()

	session.netBase.OnSessionStarted(session)
}

func (session *Session) readRoutine() {
	defer session.wgClose.Done()

	for {
		data, err := session.pktProcessor.UnPack(session)
		if err != nil {
			if err != io.EOF && !os.IsTimeout(err) {
				xtnet_go.GetLogger().LogError("session.readRoutine: err=%v", err)
			}
			session.doClose(byRead, false)
			return
		}

		session.agent.HandlerSessionData(session, data)
	}
}

func (session *Session) writeRoutine() {
	defer session.wgClose.Done()

FOR:
	for {
		select {
		case data := <-session.sendChan:
			if data == nil {
				break FOR
			}
			pktData := session.pktProcessor.Pack(data)
			_, err := session.conn.Write(pktData)
			if err != nil {
				xtnet_go.GetLogger().LogError("session.writeRoutine: err=%v", err)
				session.doClose(byWrite, false)
				return
			}

		case <-session.closeChan:
			break FOR
		}
	}
	session.conn.Close()
}
