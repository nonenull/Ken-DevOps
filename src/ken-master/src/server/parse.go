package server

import (
	"bytes"
	"ken-common/src/ken-tcpserver"
	"ken-master/src/logger"
	"ken-master/src/server/routers"
	"errors"
	"fmt"
)

type Parse struct {
}

func (self *Parse) Start(curPack []byte) (isKeepAlive bool, parseMap map[string]interface{}, parseErr error) {
	//network.ip      parseMap["action"]
	//-i eth0		  parseMap["args"]
	parseMap = make(map[string]interface{})
	var errText string
	lineSplit := bytes.Split(curPack, ken_tcpserver.LineTag)
	if len(lineSplit) < 3 {
		errText = fmt.Sprint("报文格式不正确: ", string(curPack))
		logger.Warning(errText)
		parseErr = errors.New(errText)
		return
	}
	actionName := bytes.ToLower(lineSplit[1])
	action, ok := routers.RoutersList[string(actionName)]
	if !ok {
		errText = fmt.Sprint("没有获取到对应的函数: ", string(actionName))
		logger.Warning(errText)
		parseErr = errors.New(errText)
		return
	}
	for k, v := range lineSplit {
		logger.Debug("lineSplit ==", k, "===", string(v))
	}

	logger.Debug("routers.RoutersList ==", routers.RoutersList)

	parseMap["action"] = action
	parseMap["args"] = lineSplit[2:]
	isKeepAlive = bytes.Equal(bytes.ToLower(lineSplit[0]), ken_tcpserver.KeepAliveTag)
	return
}
