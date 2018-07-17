package network

import (
	"errors"
	"ken-master/src/logger"
)

func GetIP(v interface{}) (string, bool, error) {
	//t := rand.Intn(5)
	//time.Sleep(time.Duration(t) * time.Second)
	request := v.(map[string]interface{})
	args := request["args"].(map[string]string)
	logger.Debug("args...==", args)
	var str string
	for k, v := range args {
		str += k + " " + v + " "
	}
	return str, true, errors.New("fucker")
}

func GetIP2(v interface{}) (string, bool, error) {
	return "1.1.1.1", true, nil
}
