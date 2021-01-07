package main

import (
	"encoding/binary"
	"fmt"
	"time"
	"xtnet/frame"
	mynet "xtnet/net"
	"xtnet/net/tcp"
)

/*
	type struct MsgHello{
		content: string
	}

	agent.SetCallback()
	agent.registMsgHandler(func(session, msg *MsgHello));

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
	testServer := tcp.NewServer("127.0.0.1:7001", 1024, binary.BigEndian, 1024)
	netAgent := frame.NewNetAgent(loop)
	netAgent.SetCbOnAccept(func(session mynet.Session) {
		fmt.Println("OnAccept")
	})
	netAgent.SetCbOnSessionData(func(session mynet.Session, data []byte) {
		msgID := data[0:4]
		fmt.Println("OnSessionData: ", msgID, data[4])
	})
	netAgent.SetCbOnSessionClose(func(session mynet.Session) {
		fmt.Println("OnSessionClose")
	})
	testServer.SetCallback(netAgent.OnAccept, netAgent.OnSessionData, netAgent.OnSessionClose)

	/*
		testServer.SetCallback(func(session mynet.Session) {
			fmt.Println("OnAccept")
		}, func(session mynet.Session, data []byte) {
			msgID := data[0:4]
			fmt.Println("OnSessionData: ", msgID, data[4])
		}, func(session mynet.Session) {
			fmt.Println("OnSessionClose")
		})
	*/

	fmt.Println(testServer)
	testServer.Start()

	time.Sleep(time.Minute * 10)
}
