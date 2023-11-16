package main

import (
	"log"
	"net"
	"sync"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-16 22:34
 */

// NIONetChannel 网络IO的非阻塞
func NIONetChannel() {
	addr := "127.0.0.1:5678"
	wg := sync.WaitGroup{}

	// 1.模拟读，体会读的阻塞状态
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		conn, _ := net.Dial("tcp", addr)
		defer conn.Close()
		log.Println(time.Now().Format("03:04:05.000"), "start read.")

		// 独立的goroutine，完成Read操作，将结果Send到channel中
		wgwg := sync.WaitGroup{}
		chRead := make(chan []byte)
		wgwg.Add(1)
		go func() {
			defer wgwg.Done()
			buf := make([]byte, 1024)
			n, _ := conn.Read(buf)
			chRead <- buf[:n]
		}()

		// 使用select+default实现非阻塞操作
		var data []byte
		select {
		case data = <-chRead:
		default:
		}

		log.Println(time.Now().Format("03:04:05.000"), "content:", string(data))
		wgwg.Wait()
	}(&wg)

	// 2.模拟写
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		l, _ := net.Listen("tcp", addr)
		defer l.Close()

		for {
			conn, _ := l.Accept()
			go func(conn net.Conn) {
				defer conn.Close()
				log.Println("connected.")

				// 阻塞时长
				time.Sleep(2 * time.Second)
				conn.Write([]byte("Non-Blocking I/O"))
			}(conn)
		}
	}(&wg)

	wg.Wait()
}
