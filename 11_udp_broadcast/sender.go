package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-16 10:24
 */

// UDPSenderBroadcast 广播的发送端
func UDPSenderBroadcast() {
	// 1.处理监听地址
	laddr, err := net.ResolveUDPAddr("udp", ":9876")
	if err != nil {
		log.Fatalln(err)
	}

	// 2.建立连接
	udpConn, err := net.ListenUDP("udp", laddr)

	// 3.发送数据
	// 广播地址
	rAddress := "192.168.1.255:6789"
	raddr, err := net.ResolveUDPAddr("udp", rAddress)
	if err != nil {
		log.Fatalln(err)
	}

	// 循环发送
	for {
		data := fmt.Sprintf("[%s]: %s",
			time.Now().Format("03:04:05.000"), "hello!")
		wn, err := udpConn.WriteToUDP([]byte(data), raddr)
		if err != nil {
			log.Println(err)
		}
		log.Printf("send \"%s\"(%d) to %s\n", data, wn, raddr)
		time.Sleep(time.Second)
	}
}
