package main

import (
	"encoding/gob"
	"errors"
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
* @Create 2023-11-01 22:33
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
	go ClientReceivePing(conn, &wg)

	wg.Wait()
}

func ClientReceivePing(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	type Message struct {
		ID   uint      `json:"id,omitempty"`
		Code string    `json:"code,omitempty"`
		Time time.Time `json:"time,omitempty"`
	}
	message := Message{}

	for {
		decoder := gob.NewDecoder(conn)                  // 创建解码器，将前面的json包改为gob包即可
		if err := decoder.Decode(&message); err != nil { // 解码操作
			log.Println(err)
			if errors.Is(err, io.EOF) {
				log.Println("连接已关闭，读取结束！")
				break
			}
		}
		log.Println(message)
		if message.Code == "PING-SERVER" {
			ClientSendPong(conn)
		}
	}
}

func ClientSendPong(conn net.Conn) {

	type Message struct {
		ID   uint      `json:"id,omitempty"`
		Code string    `json:"code,omitempty"`
		Time time.Time `json:"time,omitempty"`
	}
	pongMsg := Message{
		ID:   uint(rand.Int()),
		Code: "PONG-CLIENT",
		Time: time.Now(),
	}

	encoder := gob.NewEncoder(conn)                 // 创建编码器，将前面的json包改为gob包即可
	if err := encoder.Encode(pongMsg); err != nil { // 利用编码器进行编码
		log.Println(err)
		return
	}
	log.Println("pong was send to", conn.RemoteAddr())
	return
}
