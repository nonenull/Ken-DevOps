package routers

import (
	"ken-servant/src/server/crontroller/network"
	"ken-common/src/ken-tcpserver"
	"ken-servant/src/server/crontroller/cmd"
)

var RoutersList = ken_tcpserver.RoutersList

func init() {
	ken_tcpserver.Routers("network.getip", network.GetIP)
	ken_tcpserver.Routers("cmd.run", cmd.Run)
}
