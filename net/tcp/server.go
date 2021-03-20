package tcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	mynet "xtnet/net"
)

type Server struct {
	addr         string
	maxPkgSize   uint32
	order        binary.ByteOrder
	sendBuffSize uint32
	close        bool
	listener     net.Listener
	wgClose      sync.WaitGroup

	onAccept       mynet.OnAccept
	onSessionData  mynet.OnSessionData
	onSessionClose mynet.OnSessionClose
	agent          mynet.IAgent
}

func NewServer(addr string, maxPkgSize uint32, order binary.ByteOrder, sendBuffSize uint32) *Server {
	return &Server{
		addr:         addr,
		maxPkgSize:   maxPkgSize,
		order:        order,
		sendBuffSize: sendBuffSize,
		close:        false,
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
		//TODO
		fmt.Printf("tcp.Server.Start: %v", err)
		return false
	}

	server.listener = listener
	go server.listen()
	return true
}

func (server *Server) Close() {
	server.listener.Close()
	server.close = true
	server.wgClose.Wait()
}

func (server *Server) listen() {
	server.wgClose.Add(1)
	defer server.wgClose.Done()

	for server.close == false {
		conn, err := server.listener.Accept()
		if err != nil {
			//TODO
			continue
		}
		//TODO
		//不在这创建pktProcessor，应该由上层创建
		pktProcessor := NewPktProcessorDefault(server.maxPkgSize, server.order)
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
