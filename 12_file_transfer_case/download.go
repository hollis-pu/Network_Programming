package main

import (
	"log"
	"net"
	"os"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-16 10:59
 */

// UDPFileDownload 文件下载
func UDPFileDownload() {
	// 1.建立UDP连接
	lAddress := ":5678"
	laddr, err := net.ResolveUDPAddr("udp", lAddress)
	if err != nil {
		log.Fatalln(err)
	}
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer udpConn.Close()
	log.Printf("%s server is listening on %s", "UDP", udpConn.LocalAddr().String())

	// 2.接收文件名，并确认
	buf := make([]byte, 16*1024)
	rn, raddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalln(err)
	}
	Filename := string(buf[:rn])
	if _, err := udpConn.WriteToUDP([]byte("filename ok"), raddr); err != nil {
		log.Fatalln(err)
	}

	// 3.接收文件内容，并写入文件
	// 创建文件
	file, err := os.Create(Filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// 读取数据
	i := 0
	for {
		rn, _, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("read length:", rn)
		if _, err := file.Write(buf[:rn]); err != nil {
			log.Fatalln(err)
		}
		i++
		log.Println(i, "file write some content.")
	}
}
