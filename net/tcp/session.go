package tcp

import (
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
	xt "xtnet"
	xtNet "xtnet/net"
)

type closeType int

const (
	active closeType = iota
	byRead
	byWrite
)

const (
	sendChanSizeDefault int = 1024
)

const (
	sessionStatusNone int32 = iota
	sessionStatusInit
	sessionStatusRunning
	sessionStatusClosed
)

type Session struct {
	netBase   xtNet.INetBase
	conn      net.Conn
	pktProc   IPktProc
	status    int32
	wgClose   sync.WaitGroup
	sendChan  chan []byte
	closeChan chan int
	userData  interface{}
	agent     xtNet.IAgent
}

//TODO closed 改为状态
func (session *Session) setPktProc(pktProc IPktProc) {
	session.pktProc = pktProc
}

func (session *Session) SetAgent(agent xtNet.IAgent) {
	session.agent = agent
}

func (session *Session) Send(data []byte) {
	if atomic.LoadInt32(&session.status) == sessionStatusRunning {
		xt.GetLogger().LogWarn("tcp.Session.Send: session is closed")
		return
	}

	if len(session.sendChan) == cap(session.sendChan) {
		xt.GetLogger().LogWarn("tcp.Session.Send: sendChan is full")
	}

	session.sendChan <- data
}

//关闭session
func (session *Session) Close(waitWrite bool) {
	session.doClose(active, waitWrite)
}

//关闭session，阻塞等待关闭完成
func (session *Session) CloseBlock(waitWrite bool) {
	if session.doClose(active, waitWrite) {
		session.wgClose.Wait()
	}
}

func (session *Session) doClose(ct closeType, waitWrite bool) bool {
	if !atomic.CompareAndSwapInt32(&session.status, sessionStatusRunning, sessionStatusClosed) {
		if ct == active {
			xt.GetLogger().LogWarn("tcp.Session.Close: session has been closed")
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
	if atomic.CompareAndSwapInt32(&session.status, sessionStatusNone, sessionStatusInit) {
		session.doStart()
	}
}

func (session *Session) doStart() {
	session.wgClose.Add(2)
	go session.readRoutine()
	go session.writeRoutine()

	session.status = sessionStatusRunning
	session.netBase.OnSessionStarted(session)
}

func (session *Session) readRoutine() {
	defer session.wgClose.Done()

	for {
		data, err := session.pktProc.UnPack(session)
		if err != nil {
			if err != io.EOF && !os.IsTimeout(err) {
				xt.GetLogger().LogError("session.readRoutine: err=%v", err)
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
			pktData := session.pktProc.Pack(data)
			_, err := session.conn.Write(pktData)
			if err != nil {
				xt.GetLogger().LogError("session.writeRoutine: err=%v", err)
				session.doClose(byWrite, false)
				return
			}

		case <-session.closeChan:
			break FOR
		}
	}
	session.conn.Close()
}

type SessionCreator struct {
	sendChanSize int
}

func NewSessionCreator(sendChanSize int) ISessionCreator {
	return &SessionCreator{
		sendChanSize: sendChanSize,
	}
}

func (c *SessionCreator) CreateSession(netBase xtNet.INetBase, conn net.Conn) ISession {
	return &Session{
		netBase:   netBase,
		conn:      conn,
		status:    sessionStatusNone,
		sendChan:  make(chan []byte, c.sendChanSize),
		closeChan: make(chan int, 1),
	}
}
