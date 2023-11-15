package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-01 22:06
 */

func main() {
	TCPClient()
}

func TCPClient() {
	log.Println("客户端开始连接...")
	host := "localhost"
	port := 8888

	// 与服务器端建立连接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Println("client dial err=", err)
		return
	}
	defer conn.Close()
	log.Printf("conn 成功=%v time=%s\n", conn, time.Now().Format("2006-01-02 15:04:05"))
	wg := sync.WaitGroup{}
	wg.Add(1)
	go ClientRead(conn, &wg)
	wg.Wait()
}

func ClientRead(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	type Message struct {
		ID      uint   `json:"id,omitempty"`
		Code    string `json:"code,omitempty"`
		Content string `json:"content,omitempty"`
	}
	message := Message{}

	for {
		decoder := gob.NewDecoder(conn)
		if err := decoder.Decode(&message); err != nil { // 解码操作
			log.Println(err)
			if errors.Is(err, io.EOF) {
				log.Println("连接已关闭，读取结束！")
				break
			}
		}
		log.Println(message)
	}
}
