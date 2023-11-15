package main

import (
	"fmt"
	"log"
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
		ConcurrentRW(conn) // 进行连接的读写操作
	}

}
func ConcurrentRW(conn net.Conn) {
	// 循环接收客户端发送的数据
	defer conn.Close() // 关闭conn
	wg := sync.WaitGroup{}
	// 并发地写
	wg.Add(1)
	go Write(conn, &wg)
	// 并发地读
	wg.Add(1)
	go Read(conn, &wg)

	wg.Wait()
}

func Write(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		writeLen, err := conn.Write([]byte("send some data from server\n"))
		if err != nil {
			log.Println(err)
		}
		log.Printf("server write len is %d\n", writeLen)
		time.Sleep(time.Second * 2)
	}
}

func Read(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		buf := make([]byte, 1024)
		readLen, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
		}
		log.Printf("received from client data is: %s", string(buf[:readLen]))
	}
}
