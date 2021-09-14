package tcp

import (
	"net"
	"sync"
	"sync/atomic"
	xtnet_go "xtnet"
	mynet "xtnet/net"
)

type Server struct {
	addr         string
	sendBuffSize int
	closed       int64
	listener     net.Listener
	wgClose      sync.WaitGroup

	onAccept          mynet.OnAccept
	onSessionData     mynet.OnSessionData
	onSessionClose    mynet.OnSessionClose
	pktProcessorMaker IPktProcessorMaker
	agent             mynet.IAgent
}

func NewServer(addr string, sendBuffSize int) *Server {
	return &Server{
		addr:         addr,
		sendBuffSize: sendBuffSize,
		closed:       0,
	}
}

func (server *Server) SetCallback(onAccept mynet.OnAccept, onSessionData mynet.OnSessionData, onSessionClose mynet.OnSessionClose) {
	server.onAccept = onAccept
	server.onSessionData = onSessionData
	server.onSessionClose = onSessionClose
}

func (server *Server) SetAgent(agent mynet.IAgent) {
	server.agent = agent
}

func (server *Server) Start() bool {
	listener, err := net.Listen("tcp", server.addr)
	if err != nil {
		xtnet_go.GetLogger().LogError("tcp.Server.Start: %v", err)
		return false
	}

	server.listener = listener
	server.wgClose.Add(1)
	go server.listen()
	return true
}

func (server *Server) Close() {
	if atomic.CompareAndSwapInt64(&server.closed, 0, 1) {
		server.listener.Close()
		server.wgClose.Wait()
	}
}

func (server *Server) listen() {
	defer server.wgClose.Done()

	for atomic.LoadInt64(&server.closed) == 0 {
		conn, err := server.listener.Accept()
		if err != nil {
			xtnet_go.GetLogger().LogError("tcp.Server.listen: %v", err)
			continue
		}

		if server.pktProcessorMaker == nil {
			server.pktProcessorMaker = NewPktProcessorMaker(pktHeadSizeDefault, maxPkgSizeDefault, orderDefault)
		}

		//连接个数限制，session列表等交给上层维护，因为是多协程的，在net.server中维护这些信息，需要加锁；
		//上层可能是单协程的，维护这些可以根据情况加锁或不加锁
		pktProcessor := server.pktProcessorMaker.CreatePktProcessor()
		session := newSession(conn, pktProcessor, server.sendBuffSize)
		session.setCallback(server.onSessionData, server.onSessionClose)
		session.setAgent(server.agent)
		session.start()
		if server.agent != nil {
			server.agent.HandlerAccept(session)
		} else {
			server.onAccept(session)
		}
	}
}
