package routers

import (
	"ken-servant/src/server/crontroller/network"
	"ken-common/src/ken-tcpserver"
)

var RoutersList = ken_tcpserver.RoutersList

func init() {
	ken_tcpserver.Routers("network.getip", network.GetIP)
	ken_tcpserver.Routers("network.getip2", network.GetIP2)
}
