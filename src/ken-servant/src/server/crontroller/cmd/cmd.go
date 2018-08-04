package cmd

import (
	"ken-common/src/ken-tcpserver"
	"os/exec"
	"bytes"
	"ken-master/src/logger"
)

func Run(request *ken_tcpserver.Request) *ken_tcpserver.Response {
	if _, ok := request.KWargs["-h"]; ok {
		return &ken_tcpserver.Response{
			`
			Usage: ./master-cmd cmd.run [-s shell] [shell code]
			
			Example:
				./master-cmd cmd.run -s shell "route -n"
			Option:
				[-s shell]		(可选)shell 类型
				[shell code]	shell 代码
			`,
			true,
			"",
		}
	}
	logger.Debug("request.Args====", request.Args)
	logger.Debug("request.KWargs====", request.KWargs)
	if len(request.Args) < 1 {
		return &ken_tcpserver.Response{
			"",
			false,
			"没有输入命令",
		}
	}
	cmd := exec.Command("/bin/bash", "-c", request.Args[0])
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	return &ken_tcpserver.Response{
		out.String(),
		err == nil,
		err.Error(),
	}
}
