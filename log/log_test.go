package log

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	logger := NewLogger("output/log", 1024, false, true)
	logger.SetLogLevel(LevelDebug)
	logger.LogDebug("我是debug")
	logger.LogWarn("我是warn")
	logger.LogError("我是error")

	time.Sleep(time.Hour * 1)
}
