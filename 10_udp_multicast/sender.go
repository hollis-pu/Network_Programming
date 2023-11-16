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
* @Create 2023-11-15 20:14
 */

// UDPSenderMulticast 多播发送端
func UDPSenderMulticast() {
	// 1.建立UDP多播组连接
	address := "192.168.1.255:6789"
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalln(err)
	}

	udpConn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}

	// 2.发送内容
	// 循环发送
	for {
		data := fmt.Sprintf("[%s]: %s",
			time.Now().Format("03:04:05.000"), "hello!")
		wn, err := udpConn.Write([]byte(data))
		if err != nil {
			log.Println(err)
		}
		log.Printf("send \"%s\"(%d) to %s\n", data, wn, raddr.String())
		time.Sleep(time.Second)
	}
}
