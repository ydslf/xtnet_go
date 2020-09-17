package tcp

import (
	"fmt"
	"net"
	"sync"
)

type OnAccept func(*Session)
type OnSessionData func(*Session)
type OnSessionClose	func(*Session)

type Server struct {
	addr 			string
	onAccept		OnAccept
	onSessionData	OnSessionData
	onSessionClose	OnSessionClose
	close			bool
	listener   		net.Listener
	wgClose			sync.WaitGroup
}

func NewServer(addr string, onAccept OnAccept, onSessionData OnSessionData, onSessionClose OnSessionClose) *Server {
	return &Server{
		addr: addr,
		onAccept: onAccept,
		onSessionData: onSessionData,
		onSessionClose: onSessionClose,
		close: false,
	}
}

func (server *Server) Start() bool {
	listener, err := net.Listen("tcp", server.addr)
	if err != nil{
		//TODO
		fmt.Printf("tcp.Server.Start: %v", err)
		return false
	}

	server.listener = listener
	go server.listen()
	return true
}

func (server *Server) Close(){
	server.listener.Close()
	server.close = true
	server.wgClose.Wait()
}

func (server *Server) listen(){
	server.wgClose.Add(1)
	defer server.wgClose.Done()

	for server.close == false {
		conn, err := server.listener.Accept()
		if err != nil {
			//TODO
			continue
		}
		session := newSession(conn, server.onSessionData, server.onSessionClose)
		session.Start()
		server.onAccept(session)
	}
}
