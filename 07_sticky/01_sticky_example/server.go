package _1_sticky_example

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

func TCPServerSticky() {
	log.Println("服务器端开始监听8888端口...")

	host := "localhost"
	port := 8888

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
		go HandleConnSticky(conn)
	}
}

func HandleConnSticky(conn net.Conn) {
	log.Printf("accept conn from %s\n", conn.RemoteAddr())
	defer func() {
		log.Println("conn be closed")
		conn.Close()
	}()

	// 连续发送数据
	data := "package data."
	for i := 0; i < 50; i++ {
		_, err := conn.Write([]byte(data))
		if err != nil {
			log.Println(err)
		}
	}
}
