package main

import (
	"ken-master/src/cmd"
	"ken-master/src/logger"
)

func main(){
	logger.Debug("CMD MODE STARTED ...")
	cmd.NewCMD()
}