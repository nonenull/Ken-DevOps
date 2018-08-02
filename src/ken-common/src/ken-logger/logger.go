package ken_logger

import (
	"log"
	"os"
	"fmt"
	"io"
	"path"
)

var LevelMap = map[string]int{
	"DEBUG":   0,
	"WARNING": 1,
	"INFO":    2,
	"ERROR":   3,
}

type Logger struct {
	logger  *log.Logger
	level   string
	logPath string
	logName string
}

func (self *Logger) CreateLog() {
	filePath :=  path.Join(self.logPath, self.logName )
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND, 0666)
	if (err != nil) {
		log.Fatal("打开日志文件错误! ", err)
	}
	self.logger = log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags|log.Lshortfile)
}

func (self *Logger) OutPut(curLevel string, v ... interface{}) {
	if LevelMap[curLevel] < LevelMap[self.level] {
		return
	}
	prefix := "[" + curLevel + "] "
	content := fmt.Sprint(append([]interface{}{prefix}, v...)...)
	self.logger.Output(3, content)
}

func NewLogger(level string, logPath string, logName string) *Logger {
	myLogger := Logger{
		level: level,
		logPath:logPath,
		logName:logName,
	}
	myLogger.CreateLog()
	return &myLogger
}
