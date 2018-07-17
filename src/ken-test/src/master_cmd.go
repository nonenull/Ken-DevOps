package main

import (
	"ken-common/src/ken-tcpclient"
	"fmt"
	"strconv"
	"time"
	"strings"
)

func main() {
	for i := 0; i < 2000; i++ {
		go testCase(i)
	}
	for {
		;;
	}
}

func testCase(i int) {
	time.Sleep(1)
	client, clientErr := ken_tcpclient.NewClient(
		fmt.Sprint("127.0.0.1:6577"),
		false,
		nil,
		false,
	)
	if clientErr != nil {
		fmt.Println("连接 Master-Server 发生错误: ", clientErr)
		return
	}
	args := "-i eth" + strconv.Itoa(i)
	hostname := "DESKTOP-UUE2QDH"
	funcArgs := hostname + " network.getip " + args
	result, resultErr := client.Send(
		ken_tcpclient.NewTcpClientPack(
			"cmd.request",
			fmt.Sprint("true ", funcArgs),

		),
	)
	if resultErr != nil {
		fmt.Println("连接 Servant-Server 发生错误: ", resultErr)
		return
	}
	aa := string(result)
	con := strings.Contains(aa, args)
	if !con {
		fmt.Println(aa, "====", args, strings.Contains(aa, args))
	}
}
