package main

import (
	"ken-master/src/logger"
	"ken-master/src/cmder"
)

func main(){
	logger.Debug("CMD MODE STARTED ...")
	cmder.NewCMD()
}