package logger

import (
	"os"
	"ken-common/src/ken-logger"
	"ken-master/src/config"
)

func isDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	panic("not reached")
}

func checkPath(filaPath string) {
	if !isDirExists(filaPath) {
		os.MkdirAll(filaPath, 0666)
	}
}

var Logger *ken_logger.Logger

// 一点准备工作
func init() {
	// 检查日志目录，如果目录不存在，创建
	logPath := config.Fields.LOG_PATH
	logLevel := config.Fields.LOG_LEVEL
	logName	 := config.Fields.LOG_NAME
	checkPath(logPath)
	Logger = ken_logger.NewLogger(logLevel, logPath, logName)
}

func Debug(v ... interface{}) {
	Logger.OutPut("DEBUG", v...)
}
func Warning(v ... interface{}) {
	Logger.OutPut("WARNING", v...)
}

func Info(v ... interface{}) {
	Logger.OutPut("INFO", v...)
}

func Error(v ... interface{}) {
	Logger.OutPut("ERROR", v...)
}

func Exception(v ... interface{}) {
	Error(v...)
	os.Exit(1)
}
