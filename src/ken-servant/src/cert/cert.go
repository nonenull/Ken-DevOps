package cert

import (
	"path/filepath"
	"os"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	rd "math/rand"
	"math/big"
	"encoding/pem"
	"crypto/rand"
	"fmt"
	"ken-servant/src/config"
	"ken-servant/src/logger"
	"strings"
	"io/ioutil"
	"ken-common/src/ken-tcpclient"
	"time"
	"net"
)

const (
	LineTag string = "\r\n"
	EndTag  string = "\r\n\r\n"
)

var certPath = config.Fields.CERT_PATH
var privateCertPath string = filepath.Join(certPath, config.Fields.CERT_PRIVATE_NAME)
var publicCertPath string = filepath.Join(certPath, config.Fields.CERT_PUBLIC_NAME)

type Cert struct {
}

func (self *Cert) CreateCert() {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       []string{"Ken DevOps Tools."},
			OrganizationalUnit: nil,
			CommonName:         "Ken DevOps Tools.",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(100, 0, 0),
		IsCA:      false,
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		// 证书用途(客户端认证，数据加密)
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           self.GetMasterIP(),
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	if err != nil {
		logger.Exception("创建证书链失败: ", err)
	}
	certOut, err := os.Create(publicCertPath)
	if err != nil {
		logger.Exception("无法写入公钥: ", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(privateCertPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logger.Exception("打开私钥失败:", err)
		return
	}
	pem.Encode(
		keyOut,
		&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)},
	)
	keyOut.Close()
}

func (self *Cert) ReadCert(certPath string) (string, error) {
	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return "", err
	}
	return string(cert), nil
}

func (self *Cert) GetMasterIP() (ips []net.IP) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Exception("获取Master IP 发生错误:", err)
	}
	for _, ip := range addrs {
		if ipnet, ok := ip.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	logger.Debug("本机IP： ", ips)
	return
}

func (self *Cert) GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return hostname
}

/*
* 将证书做一个简单加密
* 获取一个随机数，将随机数所对应的行与证书第一行交换
*/
func (self *Cert) EncryptCert() string {
	pubCertData, err := self.ReadCert(publicCertPath)
	if err != nil {
		logger.Error("读取证书失败")
	}
	splitCert := strings.Split(pubCertData, "\n")
	certLen := len(splitCert)
	randInt := rd.Intn(certLen - 3)
	// 将随机行数与第一行交换
	tmpCert := splitCert[1]
	splitCert[1] = splitCert[randInt]
	splitCert[5] = tmpCert
	// 交换之后重新合并
	newPubCertData := strings.Join(splitCert, "\n")
	return fmt.Sprintf("%s%d", newPubCertData, randInt)
}

/*
*	将证书传给master
*/
func (self *Cert) SendCertToMaster() {
	client, clientErr := ken_tcpclient.NewClient(
		fmt.Sprint(config.Fields.MASTER_HOST, ":", config.Fields.MASTER_PORT),
		false,
		nil,
		false,
	)
	if clientErr != nil {
		logger.Error("连接Master-Server失败:", clientErr)
		return
	}
	_, err := client.Send(
		ken_tcpclient.NewTcpClientPack(
			"master.addCert",
			self.GetHostname()+LineTag+self.EncryptCert(),
		),
	)
	if err != nil {
		logger.Error("发送证书至Master-Server失败", err)
		return
	}
	logger.Info("证书发送完毕")
}

func isExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

/*
*	创建并发送证书
*/
func CreateAndSend() {
	privateOK, _ := isExist(privateCertPath)
	publicOK, _ := isExist(publicCertPath)
	logger.Debug(privateOK, publicOK)
	if !privateOK || ! publicOK {
		cert := &Cert{}
		cert.CreateCert()
		cert.SendCertToMaster()
	}
}
