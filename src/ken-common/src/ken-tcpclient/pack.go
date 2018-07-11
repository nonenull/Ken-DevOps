package ken_tcpclient

import "ken-common/src/ken-config"
var(
	LineTag = ken_config.LineTag
	EndTag = ken_config.EndTag
)

type TcpClientPack struct {
	Function string
	Args string
}

func (self *TcpClientPack)Build() string{
	return self.Function + LineTag + self.Args + EndTag
}

func NewTcpClientPack(function, args string) string{
	pack := &TcpClientPack{
		function,
		args,
	}
	return pack.Build()
}