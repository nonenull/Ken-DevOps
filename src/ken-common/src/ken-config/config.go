package ken_config

import (
	"os"
	"log"
	"bufio"
	"io"
	"strings"
	"reflect"
	"strconv"
)

type Config struct {
	confFd    *os.File
	ConfPath  string
	fieldsMap map[string]string
}

func (self *Config) OpenConf() {
	confFd, err := os.Open(self.ConfPath)
	if err != nil {
		log.Fatal("读取配置发生错误: ", err)
	}
	self.confFd = confFd
}

func (self *Config) ParseConf(fieldStruct IFields) interface{} {
	self.fieldsMap = make(map[string]string)
	defer self.confFd.Close()
	br := bufio.NewReader(self.confFd)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		strLine := string(line)
		trimLine := strings.Replace(strLine, " ", "", -1)
		// log.Println("trimLine", trimLine)
		// 忽略注释行和空行
		if trimLine == "" || trimLine[0] == '#' {
			continue
		}
		splitLine := strings.Split(trimLine, "=")
		self.fieldsMap[splitLine[0]] = splitLine[1]
	}
	// 将 MAP 映射到结构体中
	//fieldStruct.Mapping(self.fieldsMap)
	self.Mapping(fieldStruct, self.fieldsMap)
	return fieldStruct
}

/*
*	将配置信息映射到结构体里面
*	如果 配置文件里面没有配置设置项，将设置为defaultFields 里面对应的值
*/
func (self *Config) Mapping(fieldStruct interface{}, fieldsMap map[string]string) {
	obj := reflect.ValueOf(fieldStruct).Elem()
	typeObj := reflect.TypeOf(fieldStruct).Elem()
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		typeField := typeObj.Field(i)
		name := typeField.Name
		value, isKeyExist := fieldsMap[name]
		// 设置默认值
		if !isKeyExist {
			value = string(typeField.Tag)
			//log.Println("value==",value)
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int:
			//log.Println("value~~~~==",value)
			int64Val, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Fatal("配置项的值非法：", name, " = ", value)
			}
			field.SetInt(int64Val)
		default:
			log.Fatal("存在非法配置项：", name, " = ", value)
			continue
		}
	}
}

/*
*	自定义的Fields接口
*/
type IFields interface {
}