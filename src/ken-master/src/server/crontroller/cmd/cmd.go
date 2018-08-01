package cmd

import (
	"ken-common/src/ken-tcpserver"
	"strings"
	"encoding/json"
	"ken-master/src/logger"
	"ken-common/src/ken-config"
)

/*
*	代理 master-cmd 的请求, 转发到 相应的servant上去
*/
func Request(request *ken_tcpserver.Request) (response *ken_tcpserver.Response) {
	keepAlive := request.Args[0]
	hostname := request.Args[1]
	function := request.Args[2]
	args := strings.Join(request.Args[3:], ken_config.LineTag)
	responseData, responseErr := NewProxyCMD(
		hostname,
		function,
		args,
		keepAlive == "true",
	)
	//json str 转struct
	response  = &ken_tcpserver.Response{}
	if jsonErr := json.Unmarshal(responseData, response); jsonErr != nil {
		logger.Debug("获取结果发生错误jsonErr : ", jsonErr.Error())
		logger.Debug("获取结果发生错误responseData : ", string(responseData))
		logger.Debug("获取结果发生错误responseErr : ", responseErr)
		response.Error = responseErr.Error()
	}
	//logger.Debug("responsedata==", string(responseData))
	return
}
