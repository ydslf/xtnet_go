package xtnet_go

import "xtnet/log"

var logger *log.Logger

func SetLogger(l *log.Logger) {
	logger = l
}

func GetLogger() *log.Logger {
	return logger
}
