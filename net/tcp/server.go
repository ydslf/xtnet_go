package tcp

import (
	"net"
	"sync"
	"sync/atomic"
	xtnet "xtnet"
	xtnetNet "xtnet/net"
)

const (
	serverStatusClosed int32 = iota
	serverStatusStarting
	serverStatusStarted
)

type Server struct {
	addr           string
	status         int32
	listener       net.Listener
	wgClose        sync.WaitGroup
	sessionCreator ISessionCreator
	pktProcCreator IPktProcCreator
	agent          xtnetNet.IServerAgent
}

func NewServer(addr string, agent xtnetNet.IServerAgent) *Server {
	return &Server{
		addr:   addr,
		status: serverStatusClosed,
		agent:  agent,
	}
}

func (server *Server) Start() bool {
	if !atomic.CompareAndSwapInt32(&server.status, serverStatusClosed, serverStatusStarting) {
		xtnet.GetLogger().LogWarn("tcp.Server.Start: server is not closed")
		return false
	}

	listener, err := net.Listen("tcp", server.addr)
	if err != nil {
		xtnet.GetLogger().LogError("tcp.Server.Start: %v", err)
		atomic.StoreInt32(&server.status, serverStatusClosed)
		return false
	}

	atomic.StoreInt32(&server.status, serverStatusStarted)
	server.listener = listener
	server.wgClose.Add(1)
	go server.listen()
	return true
}

func (server *Server) Close() {
	if !atomic.CompareAndSwapInt32(&server.status, serverStatusStarted, serverStatusClosed) {
		xtnet.GetLogger().LogWarn("tcp.Server.Close: server is not started")
		return
	}

	atomic.StoreInt32(&server.status, serverStatusClosed)
	server.listener.Close()
	server.wgClose.Wait()
	//TODO 使用context关闭所有session
}

func (server *Server) listen() {
	defer server.wgClose.Done()

	for atomic.LoadInt32(&server.status) == serverStatusStarted {
		conn, err := server.listener.Accept()
		if err != nil {
			xtnet.GetLogger().LogError("tcp.Server.listen: %v", err)
			continue
		}

		if server.sessionCreator == nil {
			server.sessionCreator = NewSessionCreator(sendChanSizeDefault)
		}

		if server.pktProcCreator == nil {
			server.pktProcCreator = NewPktProcCreator(maxPkgSizeDefault, orderDefault)
		}

		//连接个数限制，session列表等交给上层维护，因为是多协程的，在net.server中维护这些信息，需要加锁；
		//上层可能是单协程的，维护这些可以根据情况加锁或不加锁
		session := server.sessionCreator.CreateSession(conn)
		pktProc := server.pktProcCreator.CreatePktProc()
		session.setPktProc(pktProc)
		session.SetSessionStartCb(func(session xtnetNet.ISession) {
			server.agent.HandleAccept(server, session)
		})
		session.SetSessionDataCb(func(session xtnetNet.ISession, data []byte) {
			server.agent.HandleSessionData(server, session, data)
		})
		session.SetSessionCloseCb(func(session xtnetNet.ISession) {
			server.agent.HandleSessionClose(server, session)
		})
		session.start()
	}
}
