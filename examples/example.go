package main

import (
	"fmt"
	"time"
	xt "xtnet"
	"xtnet/encoding"
	"xtnet/frame"
	"xtnet/log"
	xtNet "xtnet/net"
	"xtnet/net/agent/session"
	"xtnet/net/eventhandler"
	"xtnet/net/packet"
	"xtnet/net/rpc"
	"xtnet/net/tcp"
	xtTimer "xtnet/timer"
)

/*
	type struct MsgHello{
		content: string
	}

	agent.setNetRpc(netRpc)
	agent.SetCallback()
	agent.registMsgHandler(msgID, func(session, msg *MsgHello));

	testServer.SetCallback(func(session *tcp.Session) {
		agent.OnConnect()
	}, func(session *tcp.Session, data []byte) {
		agent.OnData()
	}, func(session *tcp.Session) {
		agent.OnClose()
	})
*/

var logger *log.Logger

func main() {
	time.After(1)
	logger = log.NewLogger("output/log", 1024*1024, true, true)
	logger.SetLogLevel(log.LevelDebug)
	xt.SetLogger(logger)

	serviceMain := frame.NewService()
	netRpc := rpc.NewRpc(serviceMain)
	netAgent := session.NewAgent(serviceMain)
	eventHandler := eventhandler.NewSessionEventHandler()
	testServer := tcp.NewServer("127.0.0.1:7001", netAgent)

	timerManager := xtTimer.NewManager(serviceMain)
	timer := timerManager.NewTimer(xtTimer.System)

	netAgent.SetNetRpc(netRpc)
	eventHandler.OnAccept = func(session xtNet.ISession) {
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
	eventHandler.OnSessionPacket = func(session xtNet.ISession, rpk *packet.ReadPacket) {
		fmt.Println("OnSessionData: ", rpk)
	}
	eventHandler.OnSessionClose = func(session xtNet.ISession) {
		fmt.Println("OnSessionClose")
	}

	fmt.Println(testServer)
	testServer.Start()
	serviceMain.Run()

	testServer.Close()
	logger.Close()
}
