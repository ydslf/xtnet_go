package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// 希望捕获所有所有 panic
	defer func() {
		r := recover()
		fmt.Println("222", r)
	}()

	// 启动新协程
	go func() {
		defer func() {
			r := recover()
			fmt.Println(r)
			buf := make([]byte, 4096)
			runtime.Stack(buf, false)
			fmt.Println("11", string(buf))
		}()

		panic(123)
	}()
	// 等待一下，不然协程可能来不及执行
	time.Sleep(5 * time.Second)
	fmt.Println("这条消息打印不出来")
}
