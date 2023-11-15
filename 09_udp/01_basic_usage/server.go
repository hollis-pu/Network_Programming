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

func UDPServerBasic() {
	// 1.解析地址
	laddr, err := net.ResolveUDPAddr("udp", ":9876")
	if err != nil {
		log.Fatalln(err)
	}

	// 2.监听
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%s server is listening on %s", "UDP", udpConn.LocalAddr().String())
	defer udpConn.Close()

	// 3.读
	buf := make([]byte, 1024)
	readLen, raddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("received  %s from %s\n", string(buf[:readLen]), raddr.String())

	// 4.写
	data := []byte("received:" + string(buf[:readLen]))
	writeLen, err := udpConn.WriteToUDP(data, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("send  %s(%d) to %s\n", string(data), writeLen, raddr.String())
}
