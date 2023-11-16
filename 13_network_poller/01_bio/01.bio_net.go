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

// BIONet 网络IO的阻塞
func BIONet() {
	addr := "127.0.0.1:5678"
	wg := sync.WaitGroup{}

	// 1.模拟读，体会读的阻塞状态
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		conn, _ := net.Dial("tcp", addr)
		defer conn.Close()
		buf := make([]byte, 1024)
		log.Println(time.Now().Format("03:04:05.000"), "start read.")
		n, _ := conn.Read(buf) // 当发送端没有发送内容到buf中时，Read()操作就处于阻塞状态
		log.Println(time.Now().Format("03:04:05.000"), "content:", string(buf[:n]))
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
				conn.Write([]byte("Blocking I/O"))
			}(conn)
		}
	}(&wg)

	wg.Wait()
}
