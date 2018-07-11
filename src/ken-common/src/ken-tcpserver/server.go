package ken_tcpserver

import (
	"net"
	"fmt"
	"log"
	"crypto/tls"
	"ken-common/src/ken-config"
)

var (
	LineTag        = []byte(ken_config.LineTag)
	EndTag         = []byte(ken_config.EndTag)
	KeepAliveTag   = []byte(ken_config.KeepAliveTag)
	NoKeepAliveTag = []byte(ken_config.NoKeepAliveTag)
)

type Response struct {
	Result string `json:"result"`
	IsOK   bool   `json:"ok"`
	Error  string `json:"error"`
}

type Server struct {
	Host         string
	Port         int
	TLS          bool
	CertFilePair [2]string
	Parse        IParse
}

func (self *Server) Start() {
	listener, err, ok := self.GetListener()
	if !ok {
		log.Fatal("服务监听失败: ", err)
	}
	for {
		// 接受新连接
		var conn, acceptErr = listener.Accept()
		//log.Println("获取到新连接: ", conn.RemoteAddr())
		if acceptErr != nil {
			log.Println(fmt.Errorf("接受连接失败：", acceptErr))
			break
		}
		go func() {
			connect := Connect{
				Conn:         conn,
				Parse:        self.Parse,
			}
			connect.Handle()
		}()
	}
}

func (self *Server) GetListener() (listener net.Listener, err error, ok bool) {
	address := fmt.Sprint(self.Host, ":", self.Port)
	if !self.TLS {
		listener, err = net.Listen("tcp", address)
		ok = (err == nil)
		return
	}
	cert, err := tls.LoadX509KeyPair(self.CertFilePair[0], self.CertFilePair[1])
	if err != nil {
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err = tls.Listen("tcp", address, config)
	return listener, err, (err == nil)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
