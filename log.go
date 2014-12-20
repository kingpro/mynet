package mynet

import (
	"fmt"
	"log"
	"time"
)

type SimpleLog struct {
	logger *log.Logger
}

func NewLogger(logger *log.Logger) *SimpleLog {
	return &SimpleLog{
		logger: logger,
	}
}

func (lg *SimpleLog) Info(s string, v ...interface{}) {
	msg := "[info] " + time.Now().Format("2006-01-02 15:04:05") + " " + fmt.Sprintf(s, v...)
	lg.logger.Println(msg)
}

func (lg *SimpleLog) Error(s string, v ...interface{}) {
	msg := "[error] " + time.Now().Format("2006-01-02 15:04:05") + " " + fmt.Sprintf(s, v...)
	lg.logger.Println(msg)
}
