package server

import (
	"bytes"
	"ken-common/src/ken-tcpserver"
	"ken-master/src/logger"
	"ken-master/src/server/routers"
)

type Parse struct {
}

func (self *Parse) Start(curPack []byte) (isKeepAlive bool, parseMap map[string]interface{}, isOK bool) {
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
	parseMap["args"] = lineSplit[2:]
	isKeepAlive = bytes.Equal(bytes.ToLower(lineSplit[0]), ken_tcpserver.KeepAliveTag)
	isOK = true
	//logger.Debug("parse keepalive==",bytes.ToLower(lineSplit[0]), ken_tcpserver.KeepAliveTag )
	return
}
