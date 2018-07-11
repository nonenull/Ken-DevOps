package server

import (
	"bytes"
	"ken-common/src/ken-tcpserver"
	"ken-master/src/logger"
	"ken-master/src/server/routers"
)

type Parse struct {
}

func (self *Parse) Start(curPack []byte) (bool, map[string]interface{}, bool) {
	//network.ip      parseMap["action"]
	//-i eth0		  parseMap["args"]
	var parseMap = make(map[string]interface{})
	lineSplit := bytes.Split(curPack, ken_tcpserver.LineTag)
	if len(lineSplit) < 3 {
		logger.Warning("报文格式不正确！\n", string(curPack))
		return false, nil, false
	}
	actionName := bytes.ToLower(lineSplit[1])
	action, ok := routers.RoutersList[string(actionName)]
	if !ok {
		logger.Warning("没有获取到对应的函数！\n", string(actionName))
		return false, nil, false
	}
	parseMap["action"] = action
	parseMap["args"] = lineSplit[2:]
	//logger.Debug("parse keepalive==",bytes.ToLower(lineSplit[0]), ken_tcpserver.KeepAliveTag )
	return bytes.Equal(bytes.ToLower(lineSplit[0]), ken_tcpserver.KeepAliveTag), parseMap, true
}
