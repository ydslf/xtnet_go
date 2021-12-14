package main

import (
	"encoding/binary"
	"fmt"
	"time"
	xtnet "xtnet"
	"xtnet/encoding"
	"xtnet/frame"
	"xtnet/log"
	xtnetNet "xtnet/net"
	"xtnet/net/agent/server"
	"xtnet/net/eventhandler"
	"xtnet/net/packet"
	"xtnet/net/tcp"
	xtnetTimer "xtnet/timer"
)

var logger *log.Logger

func main() {
	time.After(1)
	logger = log.NewLogger("output/log", 1024*1024, true, true)
	logger.SetLogLevel(log.LevelDebug)
	xtnet.SetLogger(logger)
	loop := frame.NewLoop(10000, true)
	timerManager := xtnetTimer.NewManager(loop)

	timer1 := timerManager.NewTimer(xtnetTimer.System)
	eventHandler := eventhandler.NewServerEventHandler()
	eventHandler.OnAccept = func(server xtnetNet.IServer, session xtnetNet.ISession) {
		fmt.Println("OnAccept: ", session)
		timer1.Start(time.Second*5, false, func() {
			msg1, _ := encoding.Encode("abc")
			session.Send(msg1)
			msg2, _ := encoding.Encode("def")
			session.Send(msg2)
			msg3, _ := encoding.Encode("ghi")
			session.Send(msg3)
			msg4, _ := encoding.Encode("111")
			session.Send(msg4)
			msg5, _ := encoding.Encode("222")
			session.Send(msg5)
		})
	}
	eventHandler.OnSessionPacket = func(server xtnetNet.IServer, session xtnetNet.ISession, rpk *packet.ReadPacket) {
		fmt.Println("OnSessionData: ", rpk)
	}
	eventHandler.OnSessionClose = func(server xtnetNet.IServer, session xtnetNet.ISession) {
		fmt.Println("OnSessionClose")
		timer1.Stop()
	}
	serverAgent := server.NewNormal(loop, binary.BigEndian)
	serverAgent.SetEventHandler(eventHandler)
	testServer := tcp.NewServer("127.0.0.1:7001", serverAgent)
	testServer.Start()

	loop.Run()
	testServer.Close()
	logger.Close()
}
