package config

import (
	"ken-common/src/ken-config"
)

var Fields *mFields

func init() {
	confPath := "../conf/agent.conf"
	mConfig := &ken_config.Config{
		ConfPath: confPath,
	}
	mConfig.OpenConf()
	field := &mFields{}
	Fields = mConfig.ParseConf(field).(*mFields)
}
