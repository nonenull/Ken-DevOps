package server

import (
	"ken-common/src/ken-tcpserver"
	"ken-master/src/config"
	"ken-master/src/logger"
)

func NewServer() {
	host := config.Fields.MASTER_LISTEN_HOST
	port := config.Fields.MASTER_LISTEN_PORT
	logger.Info(host, ":", port, " start listen")
	server := ken_tcpserver.Server{
		Host:         host,
		Port:         port,
		Parse:        &Parse{},
	}
	server.Start()
}
