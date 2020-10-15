package main

import (
	"encoding/binary"
	"fmt"
	"time"
	"xtnet/net/tcp"
)

func main() {
	testServer := tcp.NewServer("127.0.0.1:7001", 1024, binary.BigEndian, 1024)
	testServer.SetCallback(func(session *tcp.Session) {
		fmt.Println("OnAccept")
	}, func(session *tcp.Session, data []byte) {
		fmt.Println("OnSessionData: ", data)
	}, func(session *tcp.Session) {
		fmt.Println("OnSessionClose")
	})
	fmt.Println(testServer)
	testServer.Start()

	time.Sleep(time.Minute * 10)
}
