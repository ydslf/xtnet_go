package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func testFile() {
	logFile, err := os.OpenFile("log/123.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("读取日志文件失败", err)
	}
	defer logFile.Close()
	logFile.WriteString("test\n")

	fs, _ := logFile.Stat()
	fmt.Println("size=", fs.Size())

	logFile.WriteString("test\n")
	fs, _ = logFile.Stat()
	fmt.Println("size=", fs.Size())
}

func testDir() {
	err := os.MkdirAll("output/log", os.ModePerm)
	if err != nil {
		fmt.Println("err: ", err)
	}
}

func testTime() {
	fmt.Println(time.Now().Format("20060102_150405"))
}

func main() {
	//testFile()
	//testDir()
	testTime()
}
