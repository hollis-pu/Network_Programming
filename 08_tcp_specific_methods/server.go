package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-08 12:12
 */

func TcpServerSpecial() {
	log.Println("服务器端开始监听8888端口...")

	// 服务器端开始监听端口
	laddr, err := net.ResolveTCPAddr("tcp", ":8888")
	if err != nil {
		log.Fatalln("listen err=", err)
	}
	tcpListener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatalln("listen err=", err)
	}
	defer tcpListener.Close()
	log.Println("listen=", tcpListener)

	// 循环等待客户端的连接
	for {
		fmt.Println("等待客户端连接...")
		conn, err := tcpListener.AcceptTCP() // 等待客户端连接
		if err != nil {
			log.Println("Accept() err=", err)
		} else {
			log.Printf("Accept() suc= %v 客户端ip=%v time=%s\n", conn, conn.RemoteAddr().String(), time.Now().Format("2006-01-02 15:04:05"))
		}
		go HandleTcpConnSpecial(conn)
	}
}
func HandleTcpConnSpecial(tcpConn *net.TCPConn) {
	// 循环接收客户端发送的数据
	defer tcpConn.Close() // 关闭conn

	// 设置连接属性
	tcpConn.SetKeepAlive(true)

	for {
		// 创建新的切片
		buf := make([]byte, 1024)
		log.Printf("服务器在等待客户端 %s 发送信息\n", tcpConn.RemoteAddr().String())
		n, err := tcpConn.Read(buf) // 从conn读取，如果客户端没发送任何数据，则一直在此等待
		if err == io.EOF {
			//fmt.Println("服务器端Read err=", err)
			log.Println("客户端已关闭，其process服务也将关闭！")
			return // 如果客户端已经关闭了，则本协程也关闭掉
		}
		// 显示客户端发送到服务器端的内容
		log.Printf(string(buf[:n])) // 注意，这里只取长度为n的内容，才是实际读到的内容
	}
}
