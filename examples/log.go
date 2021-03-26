package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"
	xtlog "xtnet/log"
)

func test1() {
	log.SetPrefix("xt: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.Println("a")
	log.Fatalln("b")
	log.Println("aa")
	log.Panicln("c")
}

func test2() {
	file, err := os.OpenFile("errors.txt",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Trace := log.New(ioutil.Discard,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info := log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning := log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error := log.New(io.MultiWriter(file, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Trace.Println("I have something standard to say")
	Info.Println("Special Information")
	Warning.Println("There is something you need to know about")
	Error.Println("Something has failed")
}

var (
	kernel32    *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
	proc        *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
	CloseHandle *syscall.LazyProc = kernel32.NewProc(`CloseHandle`)

	// 给字体颜色对象赋值
	FontColor Color = Color{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
)

type Color struct {
	black        int // 黑色
	blue         int // 蓝色
	green        int // 绿色
	cyan         int // 青色
	red          int // 红色
	purple       int // 紫色
	yellow       int // 黄色
	light_gray   int // 淡灰色（系统默认值）
	gray         int // 灰色
	light_blue   int // 亮蓝色
	light_green  int // 亮绿色
	light_cyan   int // 亮青色
	light_red    int // 亮红色
	light_purple int // 亮紫色
	light_yellow int // 亮黄色
	white        int // 白色
}

func testColor() {
	if runtime.GOOS == "windows" {
		//handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.red))
		//print("windows1")
		//fmt.Println("windows2")
		//fmt.Println("windows3")
		//CloseHandle.Call(handle)
	} else {
		fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor Linux", 0x1B)
		fmt.Printf("linux ")
	}

	fmt.Printf("\033[33m")
	fmt.Println("linux1")
	fmt.Println("linux2")
	fmt.Printf("\033[0m")
	fmt.Println("linux3")

	fmt.Printf("\033[1;31m%s\033[0m\n", "Red.")
	fmt.Printf("\033[1;37m%s\033[0m\n", "Red.")
}

func testMyLog() {
	logger := xtlog.NewLogger("output/log", 1024, false)
	logger.SetLogLevel(xtlog.LevelDebug)
	logger.LogDebug("我是debug")
	logger.LogWarn("我是warn")
	logger.LogError("我是error")

	time.Sleep(time.Second * 3)
}

func main() {
	//test1()
	//test2()
	//testColor()
	testMyLog()
}
