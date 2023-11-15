package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-10-27 21:42
 */
func main() {
	fmt.Println("服务器端开始监听8888端口...")

	host := "localhost"
	port := 8888

	// 服务器端开始监听端口
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("listen err=", err)
		return
	}
	defer listen.Close()
	fmt.Println("listen=", listen)

	// 循环等待客户端的连接
	for {
		fmt.Println("等待客户端连接...")
		conn, err := listen.Accept() // 等待客户端连接
		if err != nil {
			fmt.Println("Accept() err=", err)
		} else {
			fmt.Printf("Accept() suc= %v 客户端ip=%v time=%s\n", conn, conn.RemoteAddr().String(), time.Now().Format("2006-01-02 15:04:05"))
		}
		go process(conn)
	}
}
func process(conn net.Conn) {
	// 循环接收客户端发送的数据
	defer conn.Close() // 关闭conn

	for {
		// 创建新的切片
		buf := make([]byte, 1024)
		fmt.Printf("服务器在等待客户端 %s 发送信息\n", conn.RemoteAddr().String())
		n, err := conn.Read(buf) // 从conn读取，如果客户端没发送任何数据，则一直在此等待
		if err == io.EOF {
			//fmt.Println("服务器端Read err=", err)
			fmt.Println("客户端已关闭，其process服务也将关闭！")
			return // 如果客户端已经关闭了，则本协程也关闭掉
		}
		// 显示客户端发送到服务器端的内容
		fmt.Printf(string(buf[:n])) // 注意，这里只取长度为n的内容，才是实际读到的内容
	}
}
