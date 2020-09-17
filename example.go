package main

import (
	"fmt"
	"xtnet/net/tcp"
)

func main(){
	testServer := tcp.NewServer("127.0.0.1", func(session *tcp.Session){

	})
	fmt.Println(testServer)
}