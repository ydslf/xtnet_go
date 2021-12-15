package websocket

import (
	"github.com/gorilla/websocket"
	"net"
	"net/http"
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
	addr       string
	status     int32
	listener   net.Listener
	httpServer *http.Server
	wgClose    sync.WaitGroup
	agent      xtnetNet.IServerAgent
}

func NewServer(addr string, agent xtnetNet.IServerAgent) *Server {
	return &Server{
		addr:   addr,
		status: serverStatusClosed,
		agent:  agent,
	}
}

type Handler struct {
	server *Server
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"avatar-fight"},
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		xtnet.GetLogger().LogError("ServerHttp upgrader.Upgrade: v%", err)
		return
	}

	session := NewSession(conn)
	session.SetSessionStartCb(func(session xtnetNet.ISession) {
		h.server.agent.HandleAccept(h.server, session)
	})
	session.SetSessionDataCb(func(session xtnetNet.ISession, data []byte) {
		h.server.agent.HandleSessionData(h.server, session, data)
	})
	session.SetSessionCloseCb(func(session xtnetNet.ISession) {
		h.server.agent.HandleSessionClose(h.server, session)
	})
	session.start()
}

func (server *Server) Start() bool {
	if !atomic.CompareAndSwapInt32(&server.status, serverStatusClosed, serverStatusStarting) {
		xtnet.GetLogger().LogWarn("websocket.Server.Start: server is not closed")
		return false
	}

	listener, err := net.Listen("tcp", server.addr)
	if err != nil {
		xtnet.GetLogger().LogError("websocket.Server.Start: %v", err)
		atomic.StoreInt32(&server.status, serverStatusClosed)
		return false
	}

	atomic.StoreInt32(&server.status, serverStatusStarted)
	server.listener = listener
	server.httpServer = &http.Server{
		Handler: new(Handler),
	}

	go func() {
		if err := server.httpServer.Serve(server.listener); err != nil {
			xtnet.GetLogger().LogError("websocket.Server.Start: %v", err)
			atomic.StoreInt32(&server.status, serverStatusClosed)
			return
		}
	}()

	return true
}

func (server *Server) Close() {
	if !atomic.CompareAndSwapInt32(&server.status, serverStatusStarted, serverStatusClosed) {
		xtnet.GetLogger().LogWarn("tcp.Server.Close: server is not started")
		return
	}

	atomic.StoreInt32(&server.status, serverStatusClosed)
	server.httpServer.Shutdown(nil)
}
