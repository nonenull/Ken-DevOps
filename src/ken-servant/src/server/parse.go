package server

import (
	"bytes"
	"ken-common/src/ken-tcpserver"
	"ken-master/src/logger"
	"ken-servant/src/server/routers"
	"strings"
	"fmt"
)

type Parse struct {
}

func (self *Parse) Start(curPack []byte)  (isKeepAlive bool, parseMap map[string]interface{}, isOK bool) {
	//network.ip      parseMap["action"]
	//-i eth0		  parseMap["args"]
	parseMap = make(map[string]interface{})
	lineSplit := bytes.Split(curPack, ken_tcpserver.LineTag)
	if len(lineSplit) < 3 {
		logger.Warning("报文格式不正确！\n", string(curPack))
		return
	}
	actionName := bytes.ToLower(lineSplit[1])
	action, ok := routers.RoutersList[string(actionName)]
	if !ok {
		logger.Warning("没有获取到对应的函数！\n", string(actionName))
		return
	}
	parseMap["action"] = action
	parseMap["args"] = self.ParseArgs(string(lineSplit[2]))
	return
}

/*
*	解析参数
*	-i eth0 -t fuck 转化为 map[string]string
*/
func (self *Parse) ParseArgs(argsStr string) map[string]string {
	argsMap := make(map[string]string)
	var argsSplit []string
	argsSplit = strings.Split(argsStr, "-")
	//logger.Debug("argsSplit===", argsSplit)
	if len(argsSplit) < 2 {
		logger.Warning("发现不合格参数: ", argsSplit)
		return argsMap
	}
	// 参数是  -i xxx -b xxx 类型的
	for _, value := range argsSplit {
		if value == "" {
			continue
		}
		valueSplit := strings.Split(value, " ")
		argsMap[fmt.Sprint("-", valueSplit[0])] = valueSplit[1]
	}
	return argsMap
}
