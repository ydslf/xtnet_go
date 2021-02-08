package main

import (
	"fmt"
	"time"
	xttimer "xtnet/timer"
)

func main() {
	/*
		for i:=0; i<100; i++{
			time.AfterFunc(time.Millisecond*1, func(){
				fmt.Println("timer1")
			})
		}
	*/
	for i := 0; i < 100; i++ {
		xttimer.AfterFunc(time.Millisecond*1, func() {
			fmt.Println("hahah")
		})
	}
	time.Sleep(time.Hour * 1)
}
