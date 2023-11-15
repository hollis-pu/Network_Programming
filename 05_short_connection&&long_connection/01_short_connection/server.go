package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
)

/**
* Description:
	短连接的使用，发送一次数据之后，关闭连接。
* @Author Hollis
* @Create 2023-11-01 22:06
*/

func main() {
	TCPServer()
}

func TCPServer() {
	fmt.Println("服务器端开始监听8888端口...")

	host := "localhost"
	port := 8888

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("listener err=", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go HandlerConn(conn) // 进行连接的读写操作
	}
}

func HandlerConn(conn net.Conn) {
	defer conn.Close()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go ServerWrite(conn, &wg)

	wg.Wait()
}

func ServerWrite(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	// 向客户端发送数据
	// 定义需要发送的数据
	type Message struct {
		ID      uint   `json:"id,omitempty"`
		Code    string `json:"code,omitempty"`
		Content string `json:"content,omitempty"`
	}

	message := Message{
		ID:      uint(rand.Int()),
		Code:    "SERVER-STANDARD",
		Content: "message from server",
	}

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(message); err != nil { // 利用编码器进行编码
		log.Println(err)
		return
	}
	fmt.Printf("message was send!\n")
	return
}
