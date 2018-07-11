package main

import (
	"net"
	"time"
	"fmt"
)

var data = []byte("1Keep-Alive\r\nnetwork.getip\r\n-i eth0\r\n\r\n")

func test(i int) {
	count := 0
	conn, err := net.Dial("tcp", "127.0.0.1:6578")
	if err != nil {
		//handle error
	}
	var nData []byte
	for _ = range time.Tick(time.Second*3) {
		ff := fmt.Sprint("[", i, "-", count, "]")
		nData = append([]byte(ff), data...)
		conn.Write(nData)
		fmt.Println("i====", i, "-------", count)
		count++
	}
}

func long() {
	for i := 0; i < 10000; i++ {
		go test(i)
	}
}

func short() {
	for i := 0; i < 1; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:6578")
		defer conn.Close()
		if err != nil {
			fmt.Println("fucker error")
			//handle error
		}
		conn.Write(data)
		readBuf := make([]byte, 2)
		var datas []byte
		for {
			_, err := conn.Read(readBuf)
			// 当有错误时间发生时，跳出循环，将断开连接
			// 短链接在此触发io.EOF,跳出循环，断开连接
			//fmt.Println("err===", err)
			if err != nil {
				break
			}
			datas = append(datas, readBuf...)
			//fmt.Println("hhhh===", string(datas))
		}
		fmt.Println("结束===", string(datas))
	}
}

func main() {
	//long()
	go short()
	mchan := make(chan int)
	mchan <- 1
}
