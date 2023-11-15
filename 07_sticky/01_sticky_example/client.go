package _1_sticky_example

import (
	"fmt"
	"log"
	"net"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-08 9:48
 */

func TCPClientSticky() {
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

	//tcpConn, ok := conn.(*net.TCPConn)
	//if !ok {
	//	log.Println("not tcp conn")
	//}
	//err = tcpConn.SetNoDelay(false)
	//if err != nil {
	//	log.Println(err)
	//}

	buf := make([]byte, 512)
	for {
		readLen, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("received data: %s", string(buf[:readLen]))
	}
}
