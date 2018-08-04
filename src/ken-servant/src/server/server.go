package server

import (
	"ken-common/src/ken-tcpserver"
	"ken-servant/src/config"
	"ken-servant/src/logger"
	"path"
)

func NewServer() {
	host := config.Fields.SERVANT_LISTEN_HOST
	port := config.Fields.SERVANT_LISTEN_PORT
	certPath := config.Fields.CERT_PATH
	logger.Info(host, ":", port, " start listen")

	ken_tcpserver.SetTcpServerLogger(logger.Logger)
	server := ken_tcpserver.Server{
		Host: host,
		Port: port,
		TLS:  true,
		CertFilePair: [2]string{
			path.Join(certPath, config.Fields.CERT_PUBLIC_NAME),
			path.Join(certPath, config.Fields.CERT_PRIVATE_NAME),
		},
		Parse:        &Parse{},
	}
	server.Start()
}
