package main

import (
	"ken-servant/src/cert"
	"ken-servant/src/server"
)

func main() {
	cert.CreateAndSend()
	server.NewServer()
}