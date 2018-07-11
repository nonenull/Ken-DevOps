package ken_tcpserver

import (
	"net"
	"bytes"
	"reflect"
	"encoding/json"
	"ken-common/src/ken-config"
)

type IParse interface {
	Start(curPack []byte) (bool, map[string]interface{}, bool)
}

type Connect struct {
	Conn      net.Conn
	Routers   map[string]interface{}
	Parse     IParse
	KeepAlive bool
}

func (self *Connect) Handle() {
	defer self.Conn.Close()
	var data []byte
	readBuf := make([]byte, ken_config.ReadBuffSize)
	for {
		readLen, err := self.Conn.Read(readBuf)
		// 当有错误时间发生时，跳出循环，将断开连接
		// 短链接在此触发io.EOF,跳出循环，断开连接
		if err != nil {
			//log.Println("读取信息发生了错误:", err)
			break
		}
		data = append(data, readBuf[:readLen]...)
		unPackData, packData, unPackOK := self.UnPack(data)
		if !unPackOK {
			continue
		}
		data = unPackData
		var (
			parseMap    map[string]interface{}
			parseOK     bool
		)
		self.KeepAlive, parseMap, parseOK = self.Parse.Start(packData)
		if !parseOK {
			continue
		}
		result := self.RunFunc(parseMap)
		//fmt.Println("result==", result)
		self.Conn.Write([]byte(result))
		// 此处可hack, 在RunFunc中再修改keepalive的值
		if !self.KeepAlive {
			break
		}
	}
}

/*
* 	解包
*	param:
*		data []byte	连接目前接收到的所有数据
*	return:
*		data []byte		剔除解包成功数据的所有数据
*		packData []byte	解包出来的数据
*		hasPack bool	是否解包成功
*/
func (self *Connect) UnPack(data []byte) ([]byte, []byte, bool) {
	// 获取解析到的第一段数据索引
	index := bytes.Index(data, EndTag)
	// 判断是否有结束标志
	if index < 0 {
		return nil, nil, false
	}
	packData := data[0:index]
	data = data[index+len(EndTag):]
	return data, packData, true
}

/*
* 	处理解好的包
*	parseMap["conn"]	net.Conn		TCP连接对象
*	parseMap["action"]	func			函数对象
*	parseMap["args"]	[][]byte  		参数
*/
func (self *Connect) RunFunc(parseMap map[string]interface{}) string {
	parseMap["conn"] = self
	funcVal := reflect.ValueOf(parseMap["action"])
	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(parseMap)
	result := funcVal.Call(params)
	response := &Response{}
	response.Result = result[0].String()
	response.IsOK = result[1].Bool()
	if result[2].Interface() == nil {
		response.Error = ""
	} else {
		errFunc := result[2].MethodByName("Error")
		errText := errFunc.Call(make([]reflect.Value, 0))
		response.Error = errText[0].String()
	}
	jsonResponse, jsonErr := json.Marshal(response)
	if jsonErr == nil {
		return string(jsonResponse) + ken_config.EndTag
	}
	errorResponse := &Response{}
	errorResponse.IsOK = false
	errorResponse.Error = "转换执行结果发生错误,原因：" + jsonErr.Error()
	errorJsonResponse, _ := json.Marshal(errorResponse)
	return string(errorJsonResponse) + ken_config.EndTag
}
