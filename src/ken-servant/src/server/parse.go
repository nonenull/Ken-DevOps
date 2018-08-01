package server

import (
	"ken-servant/src/logger"
	"strings"
	"fmt"
	"errors"
	"ken-common/src/ken-config"
	"ken-common/src/ken-tcpserver"
	"ken-servant/src/server/routers"
)

type Parse struct{}

/*
*	大致格式:
*	lineNum			packageContent
*
*		0		isKeepAlive\r\n
*		1		function\r\n
*		2		kwargs\r\n
*		3		args\r\n
*/
func (self *Parse) Start(curPack string) (isKeepAlive bool, request *ken_tcpserver.Request, parseErr error) {
	//logger.Debug(fmt.Sprintf("%q", curPack))
	var (
		errText string
	)
	request = &ken_tcpserver.Request{}
	lineSplit := strings.Split(curPack, ken_config.LineTag)
	lineSplitLen := len(lineSplit)
	// 最短参数长度
	if lineSplitLen < 2 {
		errText = fmt.Sprint("报文格式不正确: ", curPack)
		logger.Warning(errText)
		parseErr = errors.New(errText)
		return
	}
	if lineSplitLen > 0 {
		isKeepAliveStr := lineSplit[0]
		isKeepAlive = strings.EqualFold(isKeepAliveStr, ken_config.KeepAliveTag)
	}
	if lineSplitLen > 1 {
		funcStr := lineSplit[1]
		actionName := strings.ToLower(funcStr)
		action, ok := routers.RoutersList[string(actionName)]
		if !ok {
			errText = fmt.Sprint("没有获取到对应的函数: ", actionName)
			logger.Warning(errText)
			parseErr = errors.New(errText)
			return
		}
		request.Action = action
	}
	if lineSplitLen > 2 {
		argsArray := lineSplit[2:]
		request.KWargs, request.Args = self.ParseArgs(argsArray)
	}
	return
}

/*
*	解析参数
*	-i eth0 -t fuck 转化为 map[string]string
*   "xxx xxx xxx" 类型的参数转为[]string
*/
func (self *Parse) ParseArgs(argsArray []string) (kwargs map[string]string, args []string) {
	kwargs = make(map[string]string)
	argsArrayLen := len(argsArray)
	if argsArrayLen > 0 {
		kwargsStr := argsArray[0]
		kwargsSplit := strings.Split(kwargsStr, "-")
		logger.Debug("kwargsSplit===", kwargsSplit)
		for _, value := range kwargsSplit {
			if value == "" {
				continue
			}
			valueSplit := strings.Split(value, " ")
			logger.Debug("valueSplit===", valueSplit)
			key := fmt.Sprint("-", valueSplit[0])
			if len(valueSplit) > 1 {
				kwargs[key] = valueSplit[1]
			}else{
				kwargs[key] = ""
			}
		}
	}
	if argsArrayLen > 1 {
		args = argsArray[1:]
	}
	return
}
