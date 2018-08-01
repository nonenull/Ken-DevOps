package cmd

import (
	"bytes"
	"fmt"
	"ken-common/src/ken-tcpclient"
	"io/ioutil"
	"strings"
	"path"
	"ken-master/src/logger"
	"ken-master/src/config"
	"errors"
	"sync"
)

var (
	ConnMap = make(map[string]*ken_tcpclient.Client)
	Rlock   sync.RWMutex
)

type ProxyCMD struct {
	Hostname  string
	Function  string
	Args      string
	KeepAlive bool
	// servant 相关
	certName string
	slaveIP  string
	cert     []byte
}

func (self *ProxyCMD) Start() ([]byte, error) {
	client, err := self.getClient()
	if err != nil {
		return nil, errors.New(fmt.Sprint("连接 Servant-Server 失败: ", err.Error()))
	}
	return client.Send(
		ken_tcpclient.NewTcpClientPack(
			self.Function,
			self.Args,
		),
	)
}

func (self *ProxyCMD) ReStart() ([]byte, error) {
	// 发生异常之后, 代表当前缓存的连接对象已经不可用,需要销毁
	Rlock.Lock()
	delete(ConnMap, self.certName)
	Rlock.Unlock()
	return self.Start()
}

/*
*   获取 tcp 连接
*	当设置以短链接模式时, 每一次请求都单独开辟一个新连接
*	长连接模式下, 新创建一个长连接, 此连接将被重用
*/
func (self *ProxyCMD) getClient() (*ken_tcpclient.Client, error) {
	addrStr := fmt.Sprint(self.slaveIP, ":", config.Fields.SERVANT_LISTEN_PORT)
	// 如果是短链接, 则新建一个连接
	if !self.KeepAlive {
		client, err := ken_tcpclient.NewClient(
			addrStr,
			true,
			self.cert,
			self.KeepAlive,
		)
		//logger.Debug("获取短链接", certName, "=== ", fmt.Sprintf("%p", &client.Conn))
		return client, err
	}
	// 长连接的处理逻辑
	// 此处加锁为了处理map并发不安全的问题
	Rlock.RLock()
	client, ok := ConnMap[self.certName]
	Rlock.RUnlock()
	if ok {
		logger.Debug("获取到已存在连接地址", self.certName, "=== ", fmt.Sprintf("%p", &client.Conn))
		//// 检查连接是否可用
		//if  self.isClientAlive(client) {
		return client, nil
		//}
		// 获取到的缓存连接已经断开的情况
		// 重新创建一个新的连接,并缓存
		//goto NewClient
	}
	goto NewClient

NewClient:
// 长连接模式下
// 连接库中没有缓存连接, 则新建一个,并缓存到MAP中
	Rlock.Lock()
	client, err := ken_tcpclient.NewClient(
		addrStr,
		true,
		self.cert,
		self.KeepAlive,
	)
	logger.Debug("获取新长连接", self.certName, "=== ", fmt.Sprintf("%p", &client))
	ConnMap[self.certName] = client
	Rlock.Unlock()
	return client, err
}

/*
* 	获取指定目录下的所有文件，不进入下一级目录搜索，匹配后缀过滤。
*	return:
*	slaveIP	string  获取slave的ip
*	cert	string  获取tls证书的内容
*	ok	bool  是否成功获取ip和证书
*/
func (self *ProxyCMD) searchSlave() error {
	certsPath := config.Fields.CERTS_PATH
	certDir, err := ioutil.ReadDir(certsPath)
	if err != nil {
		logger.Error("读取 ", certsPath, " 目录失败: ", err)
	}
	for _, fi := range certDir {
		// 忽略目录
		if fi.IsDir() {
			continue
		}
		// 根据 hostname 匹配文件名
		self.certName = fi.Name()
		suffix := "@" + self.Hostname
		if strings.HasSuffix(strings.ToUpper(self.certName), strings.ToUpper(suffix)) {
			//logger.Debug("ffff===", fi.Mode(), fi.Name())
			self.slaveIP = strings.Replace(self.certName, suffix, "", -1)
			//logger.Debug("slaveIP===", slaveIP)
			self.cert, _ = ioutil.ReadFile(path.Join(certsPath, self.certName))
			break
		}
	}
	if self.slaveIP == "" || bytes.Equal(self.cert, []byte("")) {
		return errors.New(fmt.Sprint("Slave: ", self.Hostname, " 不存在"))
	}
	return nil
}

func NewProxyCMD(hostname string, function string, args string, keepAlive bool) (responseData []byte, responseErr error) {
	proxyCMD := new(ProxyCMD)
	proxyCMD.Hostname = hostname
	proxyCMD.Function = function
	proxyCMD.Args = args
	proxyCMD.KeepAlive = keepAlive
	searchSlaveErr := proxyCMD.searchSlave()
	if searchSlaveErr != nil {
		err := errors.New(searchSlaveErr.Error())
		return nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Error("连接servant-server 发生异常, 命令将重新执行: ", err)
			responseData, responseErr = proxyCMD.ReStart()
		}
	}()
	responseData, responseErr = proxyCMD.Start()
	if responseErr != nil {
		logger.Error("代理连接发生错误:", responseErr)
	}
	return responseData, responseErr
}
