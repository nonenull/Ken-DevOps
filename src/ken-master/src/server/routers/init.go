package routers

import (
	"ken-master/src/server/crontroller/master"
	"ken-master/src/server/crontroller/cmd"
	"ken-common/src/ken-tcpserver"
)

var RoutersList = ken_tcpserver.RoutersList

func init() {
	ken_tcpserver.Routers("master.addcert", master.AddCert)
	ken_tcpserver.Routers("cmd.request", cmd.Request)
}
