package cmd

import (
	"strings"
	"encoding/json"
	"errors"
	"ken-common/src/ken-tcpserver"
	"ken-master/src/logger"
	"fmt"
)

func Request(v map[string]interface{}) (result string, isOK bool, err error) {
	argsSplit := strings.Split(string(v["args"].([][]byte)[0]), " ")
	//logger.Debug("argsSplit==", argsSplit)
	keepAlive := argsSplit[0]
	hostname := argsSplit[1]
	function := argsSplit[2]
	args := strings.Join(argsSplit[3:], " ")
	//conn := v["conn"].(*ken_tcpserver.Connect)
	proxyCMD := ProxyCMD{
		hostname,
		function,
		args,
		keepAlive == "true",
	}
	responseData, err := proxyCMD.Start()
	if err != nil {
		logger.Error("代理连接发生错误:", err)
	}

	//json str 转struct
	var response ken_tcpserver.Response
	if jsonErr := json.Unmarshal(responseData, &response); jsonErr == nil {
		result = response.Result
		isOK = response.IsOK
		err = errors.New(response.Error)
	} else {
		err = errors.New(fmt.Sprint("执行发生错误:", jsonErr.Error()))
	}
	return
}
