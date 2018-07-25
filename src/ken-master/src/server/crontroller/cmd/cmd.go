package cmd

import (
	"strings"
	"encoding/json"
	"errors"
	"ken-common/src/ken-tcpserver"
	"ken-master/src/logger"
)

func Request(v map[string]interface{}) (result string, isOK bool, err error) {
	argsSplit := strings.Split(string(v["args"].([][]byte)[0]), " ")
	//logger.Debug("argsSplit==", argsSplit)
	keepAlive := argsSplit[0]
	hostname := argsSplit[1]
	function := argsSplit[2]
	args := strings.Join(argsSplit[3:], " ")
	responseData, responseErr := NewRequest(
		hostname,
		function,
		args,
		keepAlive == "true",
	)
	//json str 转struct
	var response ken_tcpserver.Response
	if jsonErr := json.Unmarshal(responseData, &response); jsonErr == nil {
		result = response.Result
		isOK = response.IsOK
		err = errors.New(response.Error)
	} else {
		logger.Debug("获取结果发生错误jsonErr : ", jsonErr.Error())
		logger.Debug("获取结果发生错误responseData : ", string(responseData))
		logger.Debug("获取结果发生错误responseErr : ", responseErr)
		err = responseErr
	}
	return
}
