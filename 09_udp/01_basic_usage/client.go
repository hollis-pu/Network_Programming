package _1_basic_usage

import (
	"log"
	"net"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-12 10:14
 */

func UDPClientBasic() {
	// 1.建立连接
	raddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9876")
	if err != nil {
		log.Fatalln(err)
	}
	udpConn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}

	// 2.写
	data := []byte("Go UDP Program")
	writeLen, err := udpConn.Write(data)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("send  %s(%d) to %s\n", string(data), writeLen, raddr.String())

	// 3.读
	buf := make([]byte, 1024)
	readLen, raddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("received  %s from %s\n", string(buf[:readLen]), raddr.String())
}
