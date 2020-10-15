package tcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

type OnAccept func(*Session)
type OnSessionData func(*Session, []byte)
type OnSessionClose func(*Session)

type Server struct {
	addr         string
	maxPkgSize   uint32
	order        binary.ByteOrder
	sendBuffSize uint32
	close        bool
	listener     net.Listener
	wgClose      sync.WaitGroup

	onAccept       OnAccept
	onSessionData  OnSessionData
	onSessionClose OnSessionClose
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

func (server *Server) SetCallback(onAccept OnAccept, onSessionData OnSessionData, onSessionClose OnSessionClose) {
	server.onAccept = onAccept
	server.onSessionData = onSessionData
	server.onSessionClose = onSessionClose
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
		pkgProcessor := NewPktProcessorDefault(server.maxPkgSize, server.order)
		session := newSession(conn, pkgProcessor, server.sendBuffSize)
		session.setCallback(server.onSessionData, server.onSessionClose)
		session.start()
		server.onAccept(session)
	}
}
