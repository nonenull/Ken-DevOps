package ken_tcpserver

import (
	"net"
	"bytes"
	"reflect"
	"encoding/json"
	"ken-common/src/ken-config"
	"fmt"
	"runtime/debug"
)

type Request struct {
	Conn 		*Connect
	ActionName 	string
	Action		interface{}
	Args		[]string
	KWargs 		map[string] string
}

type Response struct {
	Result	string	`json:"result"`
	IsOK	bool	`json:"ok"`
	Error	string 	`json:"error"`
}

type IParse interface {
	Start(curPack string) (bool, *Request, error)
}

type Connect struct {
	Conn      net.Conn
	Routers   map[string]interface{}
	Parse     IParse
	KeepAlive bool
}

func (self *Connect) Handle() {
	defer self.Conn.Close()
	// 防止发生代码逻辑引发的panic错误导致进程退出
	defer func() {
		if err := recover(); err != nil{
			errText := fmt.Sprint("代码发生错误, 请在log文件里查看详细信息: ", err)
			self.Conn.Write(
				[]byte(ErrResponse(errText)),
			)
			logger.Error(fmt.Sprintf("%v\n%s", err, debug.Stack()))
		}
	}()
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
			request *Request
			parseErr error
		)
		self.KeepAlive, request, parseErr = self.Parse.Start(string(packData))
		// 传输过来的包有问题的情况下, 返回错误, 并且断开连接
		if parseErr != nil {
			self.Conn.Write([]byte(ErrResponse(parseErr.Error())))
			break
		}
		result := self.RunFunc(request)
		//fmt.Println("result==", result)
		self.Conn.Write([]byte(result))
		// 此处可hack, 在RunFunc中再修改keepalive的值
		// 如果为短连接则在发送结果后断开连接
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
	data = data[index + EndTagLen:]
	return data, packData, true
}

/*
* 	处理解好的包
*	parseMap["conn"]	net.Conn		TCP连接对象
*	parseMap["action"]	func			函数对象
*	parseMap["args"]	[][]byte  		参数
*/
func (self *Connect) RunFunc(request *Request) string {
	request.Conn = self
	funcVal := reflect.ValueOf(request.Action)
	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(request)
	response := funcVal.Call(params)[0]
	if response.Type().String() != "*ken_tcpserver.Response" {
		return ErrResponse("获取Servant返回值错误,原因：返回值类型错误, 非*ken_tcpserver.Response") + ken_config.EndTag
	}
	jsonResponse, jsonErr := json.Marshal(response.Interface())

	// 函数正常转为json并返回
	if jsonErr == nil {
		return string(jsonResponse) + ken_config.EndTag
	}
	// 转为json失败, 则生成error response 返回
	return ErrResponse("转换执行结果发生错误,原因："+jsonErr.Error()) + ken_config.EndTag
}

func ErrResponse(errText string) string {
	errorResponse := &Response{}
	errorResponse.IsOK = false
	errorResponse.Error = errText
	errorJsonResponse, _ := json.Marshal(errorResponse)
	return string(errorJsonResponse)
}
