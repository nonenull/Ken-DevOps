package master

import (
	"ken-master/src/logger"
	"strings"
	"strconv"
	"bytes"
	"ken-master/src/config"
	"io/ioutil"
	"path"
	"ken-common/src/ken-tcpserver"
)

func AddCert(v map[string]interface{}) (string, bool, error) {
	conn := v["conn"].(*ken_tcpserver.Connect)
	remoteAddr := conn.Conn.RemoteAddr().String()
	args := v["args"].([][]byte)
	hostname := args[0]
	ipAndHost := getIP(remoteAddr) + "@" + string(hostname)
	logger.Debug("host==", ipAndHost)
	randInt, cert := departCert(args[1])
	decryptCertText := decryptCert(randInt, cert)
	// 写入证书
	certSavePath := path.Join(config.Fields.CERTS_PATH, ipAndHost)
	saveErr := ioutil.WriteFile(certSavePath, decryptCertText, 0600)
	return "", saveErr == nil, saveErr
}

/*
*	将证书解除混淆
*/
func decryptCert(randInt int, cert []byte) []byte {
	lineTag := []byte("\n")
	logger.Debug("randInt==", randInt)
	logger.Debug("cert===", cert)
	splitCert := bytes.Split(cert, lineTag)
	tmpCert := splitCert[1]
	splitCert[1] = splitCert[randInt]
	splitCert[5] = tmpCert
	return bytes.Join(splitCert, lineTag)
}

/*
*	将接收到的证书分离， 得到随机数和混淆证书
*	param:
*		cert	[]byte 	接收到的混淆证书
*	return:
*		randInt	int 	随机数
*		cert 	[]byte	混淆的证书
*/
func departCert(cert []byte) (int, []byte) {
	certLen := len(cert)
	certText := cert[:certLen-1]
	randInt, _ := strconv.ParseInt(string(cert[certLen-1:]), 10, 64)
	return int(randInt), certText
}

func getIP(remoteAddr string) string {
	splitRemoteAddr := strings.Split(remoteAddr, ":")
	return splitRemoteAddr[0]
}
