package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

/**
* Description:
	长连接的心跳检测，服务器端发送心跳，客户端响应心跳。
* @Author Hollis
* @Create 2023-11-01 22:33
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
	defer func() {
		fmt.Println("connection closed...")
		conn.Close()
	}()
	wg := sync.WaitGroup{}

	// 独立的goroutine，在建立连接后，周期性发送ping
	wg.Add(1)
	go ServerSendPing(conn, &wg)

	wg.Wait()
}

func ServerSendPing(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	// ping失败的最大次数
	const maxPingNum = 3
	pingErrCounter := 0

	ctx, cancel := context.WithCancel(context.Background())

	go ServerReceivePong(conn, ctx)

	type Message struct {
		ID   uint      `json:"id,omitempty"`
		Code string    `json:"code,omitempty"`
		Time time.Time `json:"time,omitempty"`
	}

	// 周期性地发送
	ticker := time.NewTicker(2 * time.Second) // 循环定时器，每隔2秒自动写入一个数据到管道中
	for t := range ticker.C {
		pingMsg := Message{
			ID:   uint(rand.Int()),
			Code: "PING-SERVER",
			Time: t,
		}

		log.Println("ping send to", conn.RemoteAddr())
		encoder := gob.NewEncoder(conn)
		if err := encoder.Encode(pingMsg); err != nil {
			log.Println(err)
			pingErrCounter++ // 累加错误计数器
			log.Printf("pingErrCounter:%d\n\n", pingErrCounter)
			if pingErrCounter == maxPingNum { // 判断是否达到次数上限
				cancel()
				log.Println("ping错误达到3次，自动关闭连接！")
				return // 心跳失败，同时，也需要终止pong的处理
			}
		}
	}
}

func ServerReceivePong(conn net.Conn, ctx context.Context) {
	type Message struct {
		ID   uint      `json:"id,omitempty"`
		Code string    `json:"code,omitempty"`
		Time time.Time `json:"time,omitempty"`
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			message := Message{}
			decoder := gob.NewDecoder(conn)
			if err := decoder.Decode(&message); err != nil { // 解码操作
				break
			}
			if message.Code == "PONG-CLIENT" {
				log.Println("receive pong from", conn.RemoteAddr())
			}
		}
	}
}
