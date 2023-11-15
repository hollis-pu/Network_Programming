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
* @Create 2023-11-08 9:48
 */

func TCPClientDecoder() {
	host := "localhost"
	port := 8899

	// 与服务器端建立连接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Println("client dial err=", err)
		return
	}
	defer conn.Close()

	log.Printf("conn 成功=%v time=%s\n", conn, time.Now().Format("2006-01-02 15:04:05"))

	data := ""
	decoder := NewDecoder(conn) // 创建解码器
	count := 0
	for {
		if err := decoder.Decode(&data); err != nil { // 解码操作
			log.Println(err)
			break
		}
		log.Println(count, "received data:", data)
		count++
	}
}
