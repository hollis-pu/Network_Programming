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
	for {
		ConcurrentRW(conn)
		time.Sleep(time.Second * 2)
	}
}

func ConcurrentRW(conn net.Conn) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go Write(conn, &wg)

	wg.Add(1)
	go Read(conn, &wg)
}

func Write(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		writeLen, err := conn.Write([]byte("send some data from client\n"))
		if err != nil {
			log.Println(err)
		}
		log.Printf("client write len is %d\n", writeLen)
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
		log.Printf("received from server data is: %s", string(buf[:readLen]))
	}
}
