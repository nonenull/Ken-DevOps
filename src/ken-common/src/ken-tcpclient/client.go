package ken_tcpclient

import (
	"net"
	"log"
	"strings"
	"ken-common/src/ken-config"
	"crypto/tls"
	"time"
	"crypto/x509"
	"errors"
	"fmt"
	"bytes"
)

type Client struct {
	Addr      string
	Conn      net.Conn
	TLS       bool
	TLSFile   []byte
	KeepAlive bool
}

func (self *Client) Send(requestData string) ([]byte, error) {
	if !self.Check(requestData) {
		self.Conn.Close()
		return nil, errors.New("数据包格式不正确")
	}
	// 判断keepalive 如果为true, 则在requestData前添加keepalive头
	if self.KeepAlive {
		requestData = ken_config.KeepAliveTag + ken_config.LineTag + requestData
	} else {
		requestData = ken_config.NoKeepAliveTag + ken_config.LineTag + requestData
	}
	self.Conn.Write([]byte(requestData))
	readBuf := make([]byte, ken_config.ReadBuffSize)
	endTagByte := []byte(ken_config.EndTag)
	var response []byte

	if !self.KeepAlive {
		for {
			readLen, readErr := self.Conn.Read(readBuf)
			if readErr != nil {
				break
			}
			response = append(response, readBuf[:readLen]...)
		}
	} else {
		// 长连接因为没有IO.EOF错误,需要自己解包
		for {
			readLen, readErr := self.Conn.Read(readBuf)
			if readErr != nil {
				break
			}
			response = append(response, readBuf[:readLen]...)
			if bytes.Contains(response, endTagByte) {
				break
			}
		}
	}
	return bytes.TrimSuffix(response, endTagByte), nil
}

/*
*	获取 TCP 连接
*	根据self.TLS设置来决定获取加密连接或者非加密连接
*/
func (self *Client) GetConn() (conn net.Conn, err error) {
	if self.TLS {
		cert, getCertErr := self.GetCert(self.TLSFile)
		if getCertErr != nil {
			//log.Println("getCertErr 发生错误:",getCertErr)
			return nil, errors.New(fmt.Sprint("读取证书发生错误:", getCertErr.Error()))
		}
		conf := &tls.Config{
			Time:    time.Now,
			RootCAs: cert,
		}
		conn, err = tls.Dial("tcp", self.Addr, conf)
		//log.Println("获取tls conn:",conn)
		return
	}
	conn, err = net.Dial("tcp", self.Addr)
	self.Conn = conn
	return
}

func (self *Client) GetCert(certFile []byte) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certFile)
	return certPool, nil
}

/*
*	检查数据格式是否符合标准
*/
func (self *Client) Check(data string) bool {
	return strings.Contains(data, ken_config.EndTag)
}

func NewClient(Addr string, TLS bool, TLSFile []byte, KeepAlive bool) (client *Client, err error) {
	client = &Client{
		Addr:      Addr,
		TLS:       TLS,
		TLSFile:   TLSFile,
		KeepAlive: KeepAlive,
	}
	client.Conn, err = client.GetConn()
	return
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
