package cmd

import "ken-common/src/ken-tcpserver"

func Run(request *ken_tcpserver.Request) *ken_tcpserver.Response {
	return &ken_tcpserver.Response{
		"this is function Run",
		true,
		nil,
	}
}
