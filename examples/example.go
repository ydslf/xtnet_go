package main

import (
	"fmt"
	"time"
	xtnet_go "xtnet"
	"xtnet/encoding"
	"xtnet/frame"
	"xtnet/log"
	xtnet "xtnet/net"
	"xtnet/net/packet"
	"xtnet/net/tcp"
	xttimer "xtnet/timer"
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
	xtnet_go.SetLogger(logger)

	serviceMain := frame.NewService()
	loop := serviceMain.GetLoop()
	netRpc := xtnet.NewRpc(loop)
	netAgent := xtnet.NewAgent(loop)
	testServer := tcp.NewServer("127.0.0.1:7001", netAgent)

	timerManager := xttimer.NewManager(loop)
	timer := timerManager.NewTimer(xttimer.System)

	netAgent.SetNetRpc(netRpc)
	netAgent.SetOnAccept(func(session xtnet.ISession) {
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
	})
	netAgent.SetOnSessionPacket(func(session xtnet.ISession, rpk *packet.ReadPacket) {
		fmt.Println("OnSessionData: ", rpk)
	})
	netAgent.SetOnSessionClose(func(session xtnet.ISession) {
		fmt.Println("OnSessionClose")
	})

	fmt.Println(testServer)
	testServer.Start()
	loop.Run()

	testServer.Close()
	logger.Close()
}
