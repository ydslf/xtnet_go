package main

import (
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
	"xtnet/net/rpc"
	"xtnet/net/tcp"
	xtnetTimer "xtnet/timer"
)

var logger *log.Logger

func main() {
	time.After(1)
	logger = log.NewLogger("output/log", 1024*1024, true, true)
	logger.SetLogLevel(log.LevelDebug)
	xtnet.SetLogger(logger)

	serviceMain := frame.NewService()
	netRpc := rpc.NewRpc(serviceMain)
	serverAgent := server.NewNormal(serviceMain)
	eventHandler := eventhandler.NewServerEventHandler()
	testServer := tcp.NewServer("127.0.0.1:7001", serverAgent)

	timerManager := xtnetTimer.NewManager(serviceMain)
	timer := timerManager.NewTimer(xtnetTimer.System)

	serverAgent.SetEventHandler(eventHandler)
	serverAgent.SetNetRpc(netRpc)
	eventHandler.OnAccept = func(server xtnetNet.IServer, session xtnetNet.ISession) {
		fmt.Println("OnAccept: ", session)

		timer.Start(time.Second*2, false, func() {
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
			//session.CloseBlock(false)
		})
	}
	eventHandler.OnSessionPacket = func(server xtnetNet.IServer, session xtnetNet.ISession, rpk *packet.ReadPacket) {
		fmt.Println("OnSessionData: ", rpk)
	}
	eventHandler.OnSessionClose = func(server xtnetNet.IServer, session xtnetNet.ISession) {
		fmt.Println("OnSessionClose")
	}

	fmt.Println(testServer)
	testServer.Start()
	serviceMain.Run()

	testServer.Close()
	logger.Close()
}
