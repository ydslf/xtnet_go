package log

import (
	"fmt"
	"syscall"
)

const (
	Black       int = 0
	Blue        int = 1
	Green       int = 2
	Cyan        int = 3
	Red         int = 4
	Purple      int = 5
	Yellow      int = 6
	LightGray   int = 7
	Gray        int = 8
	LightBlue   int = 9
	LightGreen  int = 10
	LightCyan   int = 11
	LightRed    int = 12
	LightPurple int = 13
	LightYellow int = 14
	White       int = 15
)

var colorConfig = [LevelNum]struct {
	color int
}{
	{color: LightGray},
	{color: LightRed},
	{color: LightYellow},
	{color: LightCyan},
}

var (
	kernel32                = syscall.NewLazyDLL(`kernel32.dll`)
	setConsoleTextAttribute = kernel32.NewProc(`SetConsoleTextAttribute`)
	closeHandle             = kernel32.NewProc(`CloseHandle`)
)

func (logger *Logger) screenPrint(level int) {
	handle, _, _ := setConsoleTextAttribute.Call(uintptr(syscall.Stdout), uintptr(colorConfig[level].color))
	fmt.Print(logger.buf.String())
	handle, _, _ = setConsoleTextAttribute.Call(uintptr(syscall.Stdout), uintptr(LightGray))
	closeHandle.Call(handle)
}
