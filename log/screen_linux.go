package log

import "fmt"

const (
	Reset   int = 0
	Gray    int = 32
	Yellow  int = 31
	Green   int = 32
	Red     int = 33
	Blue    int = 34
	Magenta int = 35
	Cyan    int = 36
	White   int = 37
)

var colorConfig = [LevelNum]struct {
	color int
}{
	{color: Reset},
	{color: Red},
	{color: Yellow},
	{color: Green},
}

func (logger *Logger) screenPrint(level int) {
	fmt.Printf("\033[%dm%s\033[0m\n", colorConfig[level].color, logger.buf.String())
}
