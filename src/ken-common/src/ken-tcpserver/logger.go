package ken_tcpserver

import (
	"ken-common/src/ken-logger"
	"os"
)

var logger *TcpServerLogger

type TcpServerLogger struct {
	*ken_logger.Logger
}

func (self TcpServerLogger) Debug(v ... interface{}) {
	self.OutPut("DEBUG", v...)
}
func (self TcpServerLogger) Warning(v ... interface{}) {
	self.OutPut("WARNING", v...)
}

func (self TcpServerLogger) Info(v ... interface{}) {
	self.OutPut("INFO", v...)
}

func (self TcpServerLogger) Error(v ... interface{}) {
	self.OutPut("ERROR", v...)
}

func (self TcpServerLogger) Exception(v ... interface{}) {
	self.Error(v...)
	os.Exit(1)
}

// 一点准备工作
func SetTcpServerLogger(l *ken_logger.Logger) {
	logger = &TcpServerLogger{
		l,
	}
}