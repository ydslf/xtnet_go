package log

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

const (
	LevelNone    = 0
	LevelError   = 1
	LevelWarning = 2
	LevelDebug   = 3
)

const (
	MsgChanSize int32 = 0xFFFF
	FileSizeMin       = 1 * 1024 * 1024
	FileSizeMax       = 64 * 1024 * 1024
)

type Logger struct {
	dir           string
	fileSizeLimit int
	screenPrint   bool
	logLevel      int
	async         bool
	close         bool
	msgChan       chan *loggerMsg

	dirReal      string
	checkBaseDir bool
	curFileDay   int
	file         *os.File
	buf          bytes.Buffer
}

type loggerMsg struct {
	content string
	level   int
	time    time.Time
}

func NewLogger(dir string, fileSizeLimit int, screenPrint bool) *Logger {
	logger := new(Logger)
	logger.dir = dir
	logger.dirReal = dir
	logger.fileSizeLimit = fileSizeLimit
	logger.screenPrint = screenPrint
	logger.logLevel = LevelNone
	logger.close = false
	logger.msgChan = make(chan *loggerMsg, MsgChanSize)

	if logger.fileSizeLimit < FileSizeMin {
		logger.fileSizeLimit = FileSizeMin
	} else if logger.fileSizeLimit > FileSizeMax {
		logger.fileSizeLimit = FileSizeMax
	}

	go logger.worker()

	return logger
}

func (logger *Logger) SetLogLevel(level int) {
	logger.logLevel = level
}

func (logger *Logger) SetScreenPrint(screenPrint bool) {
	logger.screenPrint = screenPrint
}

func (logger *Logger) LogDebug(format string, v ...interface{}) {
	if logger.logLevel < LevelDebug {
		return
	}
	logger.pushLog(LevelDebug, format, v...)
}

func (logger *Logger) LogWarning(format string, v ...interface{}) {
	if logger.logLevel < LevelWarning {
		return
	}
	logger.pushLog(LevelWarning, format, v...)
}

func (logger *Logger) LogError(format string, v ...interface{}) {
	if logger.logLevel < LevelError {
		return
	}
	logger.pushLog(LevelError, format, v...)
}

func (logger *Logger) pushLog(level int, format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	msg := &loggerMsg{content, level, time.Now()}
	logger.msgChan <- msg
}

func (logger *Logger) worker() {
	for msg := range logger.msgChan {
		logger.showLog(msg.content, msg.level, msg.time)
	}
}

func (logger *Logger) showLog(content string, level int, time time.Time) {
	if !logger.createDir() {
		return
	}

	if !logger.createFile() {
		return
	}

	logger.buf.Reset()
	logger.buf.WriteString(time.Format("2006_01_02 15:04:05"))
	logger.buf.WriteString(content)
	logger.buf.WriteString("\n")
	logger.file.Write(logger.buf.Bytes())
}

func (logger *Logger) createDir() bool {
	if !logger.checkBaseDir {
		err := os.MkdirAll(logger.dir, os.ModePerm)
		if err != nil {
			fmt.Println("log create dir failed, err: ", err)
			return false
		}

		logger.dirReal = logger.dir + "/"
		logger.checkBaseDir = true

		if logger.file != nil {
			logger.file.Close()
			logger.file = nil
		}
	}
	return true
}

func (logger *Logger) createFile() bool {
	if logger.file == nil {
		return logger._createFile()
	}

	if logger.curFileDay != time.Now().YearDay() {
		return logger._createFile()
	}

	fs, _ := logger.file.Stat()
	if fs.Size() > int64(logger.fileSizeLimit) {
		return logger._createFile()
	}

	return true
}

func (logger *Logger) _createFile() bool {
	if logger.file != nil {
		logger.file.Close()
		logger.file = nil
	}

	fileName := fmt.Sprintf("%sLog_%s.txt", logger.dirReal, time.Now().Format("20060102_150405"))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("_createFile failed, filename=", fileName)
		return false
	}

	logger.curFileDay = time.Now().YearDay()
	logger.file = file
	return true
}
