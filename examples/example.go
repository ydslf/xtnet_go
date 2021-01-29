package main

import (
	"encoding/binary"
	"fmt"
	"time"
	"xtnet/frame"
	xtnet "xtnet/net"
	"xtnet/net/packet"
	"xtnet/net/tcp"
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

func main() {
	loop := frame.NewLoop()
	netRpc := xtnet.NewRpc(loop)
	netAgent := xtnet.NewAgent(loop)
	testServer := tcp.NewServer("127.0.0.1:7001", 1024, binary.BigEndian, 1024)

	netAgent.SetNetRpc(netRpc)
	netAgent.SetOnAccept(func(session xtnet.ISession) {
		fmt.Println("OnAccept")
	})
	netAgent.SetOnSessionPacket(func(session xtnet.ISession, rpk *packet.ReadPacket) {
		fmt.Println("OnSessionData: ", rpk)
	})
	netAgent.SetOnSessionClose(func(session xtnet.ISession) {
		fmt.Println("OnSessionClose")
	})
	testServer.SetAgent(netAgent)

	fmt.Println(testServer)
	testServer.Start()
	loop.Run()

	time.Sleep(time.Second * 1)
}
