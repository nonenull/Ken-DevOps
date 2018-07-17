package cmd

import (
	"ken-master/src/config"
	"strings"
	"fmt"
	"ken-common/src/ken-tcpclient"
	"ken-master/src/logger"
	"os"
	"ken-common/src/ken-tcpserver"
	"encoding/json"
)

/*
* 	命令行客户端，获取命令行参数，发给对应slave，获取执行结果
*	格式： ./command [-S] [slaverID] [func] [args]
*/
func NewCMD() {
	//logger.Debug("os.Args===",os.Args)
	args := os.Args[1:]
	if len(args) < 1 {
		logger.Exception(`
			Usage: ./command [-S] [servantID] [func] [args]
			
			Example:
				./command -S nginxserver network.getip -i eth0
			Option:
				[-S]	(可选)指定处理的连接类型为短连接, 默认为长连接
				[servantID]	servant的主机名
				[func]	在servant主机上执行的函数
				[args]	传递给执行函数的参数
		`)
		return
	}
	isKeepAlive := !(args[0] == "-S")
	var funcArgsList []string
	if !isKeepAlive {
		funcArgsList = args[1:]
	} else {
		funcArgsList = args[:]
	}
	funcArgs := strings.Join(funcArgsList, " ")
	client, clientErr := ken_tcpclient.NewClient(
		fmt.Sprint("127.0.0.1:", config.Fields.MASTER_LISTEN_PORT),
		false,
		nil,
		false,
	)
	if clientErr != nil {
		logger.Error("连接 Master-Server 发生错误: \n", clientErr)
		return
	}
	result, resultErr := client.Send(
		ken_tcpclient.NewTcpClientPack(
			"cmd.request",
			fmt.Sprint(isKeepAlive, " ", funcArgs),
		),
	)
	if resultErr != nil {
		logger.Error("连接 Servant-Server 发生错误: \n", resultErr)
		return
	}
	pretyCMD(result)
}

func pretyCMD(responseData []byte) {
	var response ken_tcpserver.Response
	if jsonErr := json.Unmarshal(responseData, &response); jsonErr == nil {
		pretyResult := fmt.Sprintf(`
	状态:
		%t
	执行结果:
		%s
	错误:
		%s
			`,response.IsOK, response.Result, response.Error)

		logger.Info(pretyResult)
	} else {
		logger.Exception("解析结果发生错误: \n", jsonErr)
	}
}
