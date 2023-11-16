package main

import (
	"log"
	"net"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-16 10:23
 */

// UDPReceiverBroadcast 广播的接收端
func UDPReceiverBroadcast() {
	// 1.广播监听地址
	laddr, err := net.ResolveUDPAddr("udp", ":6789")
	if err != nil {
		log.Fatalln(err)
	}

	// 2.广播监听
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}

	// 3.接收数据
	// 4.处理数据
	buf := make([]byte, 1024)
	for {
		n, raddr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
		}
		log.Printf("received \"%s\" from %s\n", string(buf[:n]), raddr.String())
	}
}
