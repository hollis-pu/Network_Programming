package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-10-28 12:06
 */

func main() {
	TCPServerFormat()
}

func TCPServerFormat() {
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
		go HandlerConnFormat(conn) // 进行连接的读写操作
	}
}

func HandlerConnFormat(conn net.Conn) {
	defer conn.Close()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go SerWriteFormat(conn, &wg)

	wg.Wait()
}

func SerWriteFormat(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	for {
		count++
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

		// 1.json，文本编码
		// 数据编码后发送
		//encoder := json.NewEncoder(conn)                // 创建编码器
		//if err := encoder.Encode(message); err != nil { // 利用编码器进行编码，encode成功后，会写入到conn，也就是完成了conn.Write()
		//	if err == io.EOF {
		//		fmt.Println("客户端已关闭，服务器端也将停止...")
		//		break
		//	}
		//	log.Println(err)
		//	continue
		//}
		//fmt.Printf("message %d was send!\n", count)

		// 2.GOB，二进制编码
		encoder := gob.NewEncoder(conn)                 // 创建编码器，将前面的json包改为gob包即可
		if err := encoder.Encode(message); err != nil { // 利用编码器进行编码
			if err == io.EOF {
				fmt.Println("客户端已关闭，服务器端也将停止...")
				break
			}
			log.Println(err)
			continue
		}
		fmt.Printf("message %d was send!\n", count)

		time.Sleep(time.Second * 1)
	}
}
