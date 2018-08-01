package server

import (
	"ken-common/src/ken-tcpserver"
	"strings"
	"ken-common/src/ken-config"
	"fmt"
	"errors"
	"ken-master/src/logger"
	"ken-master/src/server/routers"
)

type Parse struct {}


/*
*	大致格式:
*	lineNum			packageContent
*
*		0		masterFunction\r\n
*		1		isKeepAlive\r\n
*		2		hostname\r\n
*		3		slaverFunction\r\n
*		4		kwargs\r\n
*		5		args\r\n
*/
func (self *Parse) Start(curPack string) (isKeepAlive bool, request *ken_tcpserver.Request, parseErr error) {
	var errText string
	request = &ken_tcpserver.Request{}
	lineSplit := strings.Split(curPack, ken_config.LineTag)
	if len(lineSplit) < 3 {
		errText = fmt.Sprint("报文格式不正确: ", curPack)
		logger.Warning(errText)
		parseErr = errors.New(errText)
		return
	}
	actionName := strings.ToLower(lineSplit[1])
	request.ActionName = actionName
	action, ok := routers.RoutersList[actionName]
	if !ok {
		errText = fmt.Sprint("没有获取到对应的函数: ", actionName)
		logger.Warning(errText)
		parseErr = errors.New(errText)
		return
	}
	//logger.Debug("routers.RoutersList ==", routers.RoutersList)
	request.Args = lineSplit[2:]
	request.Action = action
	isKeepAlive = strings.EqualFold(lineSplit[0], ken_config.KeepAliveTag)
	return
}
