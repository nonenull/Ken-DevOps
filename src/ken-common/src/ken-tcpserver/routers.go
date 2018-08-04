package ken_tcpserver

var RoutersList = make(map[string]interface{})

func Routers(actionName string, actionFunc interface{}) {
	// 判断key不存在，如果存在则报错
	if _, ok := RoutersList[actionName]; ok {
		logger.Exception("routs 有重复项, 名：", actionName)
	}
	RoutersList[actionName] = actionFunc
}
