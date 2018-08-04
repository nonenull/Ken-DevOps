package network

import (
	"ken-servant/src/logger"
	"ken-common/src/ken-tcpserver"
)

func GetIP(request *ken_tcpserver.Request) *ken_tcpserver.Response {
	if _, ok := request.KWargs["-h"]; ok {
		return  &ken_tcpserver.Response{
			`sufppp get ip`,
			true,
			"",
		}
	}

	logger.Debug("kwargs...==", request.KWargs)
	logger.Debug("args...==", request.Args)
	var str string
	for k, v := range request.KWargs {
		str += k + " " + v + " "
	}
	var astr string
	for _, b := range request.Args {
		astr += b + "@"
	}
	return  &ken_tcpserver.Response{
		"this is test \n" + str + "\n" + astr,
		true,
		"fucker",
	}
}
