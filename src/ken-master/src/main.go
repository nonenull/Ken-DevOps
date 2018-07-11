package main

import (
	"ken-master/src/server"
	"ken-master/src/logger"
	"ken-master/src/config"
)

func main() {
	logger.Debug(config.Fields.LOG_LEVEL)
	server.NewServer()
}
