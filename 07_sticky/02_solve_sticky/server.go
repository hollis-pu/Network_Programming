package main

import (
	"fmt"
	"log"
	"net"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-08 9:43
 */

func TCPServerEncoder() {
	log.Println("服务器端开始监听8899端口...")

	host := "localhost"
	port := 8899

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Println("listener err=", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go HandleConnEncoder(conn)
	}
}

func HandleConnEncoder(conn net.Conn) {
	log.Printf("accept conn from %s\n", conn.RemoteAddr())
	defer func() {
		log.Println("conn be closed")
		conn.Close()
	}()

	// 连续发送数据
	data := []string{
		"package data1: hello",
		"package data2: user name is tom",
		"package data3: password is 123456",
		"package data4: over",
	}
	encoder := NewEncoder(conn) // 创建自定义编码器
	index := 0
	for i := 0; i < 50; i++ {
		if err := encoder.Encode(data[index]); err != nil { // 利用编码器进行编码
			log.Println(err)
		}
		index = (index + 1) % 4
	}
}
