package tcp

import (
	"net"
	"sync"
	"sync/atomic"
	xt "xtnet"
	xtNet "xtnet/net"
)

type Server struct {
	addr           string
	closed         int32
	listener       net.Listener
	wgClose        sync.WaitGroup
	sessionCreator ISessionCreator
	pktProcCreator IPktProcCreator
	agent          xtNet.IAgent
}

func NewServer(addr string, agent xtNet.IAgent) *Server {
	return &Server{
		addr:   addr,
		closed: 0,
		agent:  agent,
	}
}

func (server *Server) Start() bool {
	listener, err := net.Listen("tcp", server.addr)
	if err != nil {
		xt.GetLogger().LogError("tcp.Server.Start: %v", err)
		return false
	}

	server.listener = listener
	server.wgClose.Add(1)
	go server.listen()
	return true
}

func (server *Server) Close() {
	if atomic.CompareAndSwapInt32(&server.closed, 0, 1) {
		server.listener.Close()
		server.wgClose.Wait()
		//TODO 使用context关闭所有session
	}
}

func (server *Server) listen() {
	defer server.wgClose.Done()

	for atomic.LoadInt32(&server.closed) == 0 {
		conn, err := server.listener.Accept()
		if err != nil {
			xt.GetLogger().LogError("tcp.Server.listen: %v", err)
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
		session := server.sessionCreator.CreateSession(server, conn)
		pktProc := server.pktProcCreator.CreatePktProc()
		session.setPktProc(pktProc)
		session.SetAgent(server.agent)
		session.start()
	}
}

func (server *Server) OnSessionStarted(session xtNet.ISession) {
	if atomic.LoadInt32(&server.closed) == 0 {
		server.agent.HandlerAccept(session)
	}
}
