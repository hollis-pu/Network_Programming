package main

import (
	"log"
	"net"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-15 20:04
 */

// UDPReceiverMulticast 多播接收端
func UDPReceiverMulticast() {
	// 1.多播监听地址
	address := "224.1.1.2:6789"
	gaddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalln(err)
	}

	// 2.多播监听

	// 获取本地网络接口
	ifi, err := net.InterfaceByName("WLAN")
	if err != nil {
		log.Fatalln(err)
	}

	udpConn, err := net.ListenMulticastUDP("udp", ifi, gaddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%s server is listening on %s\n", "UDP", udpConn.LocalAddr().String())
	//defer udpConn.Close()

	// 3.接收数据
	// 循环接收
	buf := make([]byte, 1024)
	for {
		n, raddr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
		}
		log.Printf("received \"%s\" from %s\n", string(buf[:n]), raddr.String())
	}
}
