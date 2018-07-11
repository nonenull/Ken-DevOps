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
	testMap = make(map[string]bool)
	Rlock   sync.RWMutex
)

type ProxyCMD struct {
	Hostname  string
	Function  string
	Args      string
	KeepAlive bool
}

func (self *ProxyCMD) Start() ([]byte, error) {
	//Rlock.Lock()
	//testMap[self.Hostname] = true
	//logger.Debug("testMap===", testMap)
	//Rlock.Unlock()
	//return nil, nil
	client, err := self.getClient()
	if err != nil {
		logger.Error("获取 proxy连接失败")
		return nil, errors.New(fmt.Sprint("获取 proxy连接失败", err.Error()))
	}
	return client.Send(
		ken_tcpclient.NewTcpClientPack(
			self.Function,
			self.Args,
		),
	)
}

/*
*  获取 tcp 连接
*  此处需要判断条件:
*		1.
*/
func (self *ProxyCMD) getClient() (client *ken_tcpclient.Client, err error) {
	certName, slaveIP, cert := self.searchSlave(self.Hostname)
	if slaveIP == "" || bytes.Equal(cert, []byte("")) {
		return nil, errors.New(fmt.Sprint("Slave: ", self.Hostname, " 不存在"))
	}
	// 如果是短链接, 则新建一个连接
	if !self.KeepAlive {
		client, err = ken_tcpclient.NewClient(
			fmt.Sprint(slaveIP, ":", config.Fields.SERVANT_LISTEN_PORT),
			true,
			cert,
			self.KeepAlive,
		)
		logger.Debug("获取短链接", certName, "=== ", fmt.Sprintf("%p", &client))
		return
	}
	Rlock.RLock()
	clientObj, ok := ConnMap[certName]
	Rlock.RUnlock()
	if ok {
		client = clientObj
		logger.Debug("获取到已存在连接地址", certName, "=== ", fmt.Sprintf("%p", &client))
	} else {
		Rlock.Lock()
		client, err = ken_tcpclient.NewClient(
			fmt.Sprint(slaveIP, ":", config.Fields.SERVANT_LISTEN_PORT),
			true,
			cert,
			self.KeepAlive,
		)
		logger.Debug("获取新长连接", certName, "=== ", fmt.Sprintf("%p", &client))
		ConnMap[certName] = client
		Rlock.Unlock()
	}
	return
}

/*
* 	获取指定目录下的所有文件，不进入下一级目录搜索，匹配后缀过滤。
*	return:
*	slaveIP	string  获取slave的ip
*	cert	string  获取tls证书的内容
*	ok	bool  是否成功获取ip和证书
*/
func (self *ProxyCMD) searchSlave(hostname string) (certName string, slaveIP string, certByte []byte) {
	certsPath := config.Fields.CERTS_PATH
	certDir, err := ioutil.ReadDir(certsPath)
	if err != nil {
		logger.Error("读取 ", certsPath, " 目录失败: ", err)
		return
	}
	for _, fi := range certDir {
		// 忽略目录
		if fi.IsDir() {
			continue
		}
		// 根据 hostname 匹配文件名
		certName = fi.Name()
		suffix := "@" + hostname
		if strings.HasSuffix(strings.ToUpper(certName), strings.ToUpper(suffix)) {
			//logger.Debug("ffff===", fi.Mode(), fi.Name())
			slaveIP = strings.Replace(certName, suffix, "", -1)
			//logger.Debug("slaveIP===", slaveIP)
			certByte, _ = ioutil.ReadFile(path.Join(certsPath, certName))
			break
		}
	}
	return
}
