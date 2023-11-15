package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-10-27 21:48
 */
func main() {
	fmt.Println("客户端开始连接...")
	host := "localhost"
	port := 8888

	// 与服务器端建立连接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("client dial err=", err)
		return
	}
	defer conn.Close()
	fmt.Printf("conn 成功=%v time=%s\n", conn, time.Now().Format("2006-01-02 15:04:05"))

	reader := bufio.NewReader(os.Stdin)

	// 客户端可以一直向服务器端写入数据，每次写入一行。当输入exit时，关闭连接。
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("readString err=", err)
		}

		if line == "exit\r\n" { // 注意，这里要加上\r\n
			fmt.Println("客户端退出！")
			break
		}
		n, err := conn.Write([]byte(line)) // 通过连接向服务器端写入数据
		if err != nil {
			fmt.Println("conn.Write err=", err)
		}
		fmt.Printf("%s 客户端发送了 %d 字节的数据\n", time.Now().Format("2006-01-02 15:04:05"), n)
	}
}
