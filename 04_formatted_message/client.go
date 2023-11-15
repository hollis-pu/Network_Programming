package main

import (
	"encoding/gob"
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
	TCPClientFormat()
}

func TCPClientFormat() {
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
	go CliReadFormat(conn, &wg)
	wg.Wait()
}

func CliReadFormat(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {

		type Message struct {
			ID      uint   `json:"id,omitempty"`
			Code    string `json:"code,omitempty"`
			Content string `json:"content,omitempty"`
		}
		message := Message{}

		// 接收数据，接收到数据后，先解码
		// 1.json解码
		//decoder := json.NewDecoder(conn)                 // 创建解码器
		//if err := decoder.Decode(&message); err != nil { // 解码操作，从conn中读取内容，将解码结果放入message中
		//	log.Println(err)
		//	if err == io.EOF {
		//		fmt.Println("服务器端已关闭，客户端也将停止...")
		//		break
		//	}
		//	continue
		//}
		//log.Println(message)

		// 2.GOB解码
		decoder := gob.NewDecoder(conn)                  // 创建解码器，将前面的json包改为gob包即可
		if err := decoder.Decode(&message); err != nil { // 解码操作
			log.Println(err)
			continue
		}
		log.Println(message)
	}
}
