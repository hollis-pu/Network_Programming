package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-08 12:12
 */

func TcpClientSpecial() {
	log.Println("客户端开始连接...")

	// 与服务器端建立连接
	laddr, err := net.ResolveTCPAddr("tcp", ":8888")
	conn, err := net.DialTCP("tcp", laddr, laddr)
	if err != nil {
		log.Println("client dial err=", err)
		return
	}
	defer conn.Close()
	log.Printf("conn 成功=%v time=%s\n", conn, time.Now().Format("2006-01-02 15:04:05"))

	reader := bufio.NewReader(os.Stdin)

	// 客户端可以一直向服务器端写入数据，每次写入一行。当输入exit时，关闭连接。
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("readString err=", err)
		}

		if line == "exit\r\n" { // 注意，这里要加上\r\n
			log.Println("客户端退出！")
			break
		}
		n, err := conn.Write([]byte(line)) // 通过连接向服务器端写入数据
		if err != nil {
			log.Println("conn.Write err=", err)
		}
		log.Printf("%s 客户端发送了 %d 字节的数据\n", time.Now().Format("2006-01-02 15:04:05"), n)
	}
}
