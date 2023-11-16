package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-16 10:59
 */

// UDPFileUpload 文件上传
func UDPFileUpload() {
	// 1.获取文件信息
	// 打开文件
	filename := "./data/music.mp3"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close() // 关闭文件
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	// 主要用到的两个文件信息：fileInfo.Size()  fileInfo.Name()
	log.Println("send file size:", fileInfo.Size())

	// 2.连接服务器
	rAddress := "192.168.1.6:5678"
	raddr, err := net.ResolveUDPAddr("udp", rAddress)
	if err != nil {
		log.Fatalln(err)
	}
	udpConn, err := net.DialUDP("udp", nil, raddr)
	defer udpConn.Close()

	// 3.发送文件名
	if _, err := udpConn.Write([]byte(fileInfo.Name())); err != nil {
		log.Fatalln(err)
	}

	// 4.服务器确认
	buf := make([]byte, 16*1024)
	rn, err := udpConn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	if "filename ok" != string(buf[:rn]) {
		log.Fatalln(errors.New("server not ready"))
	}

	// 5.发送文件内容
	// 读取文件内容，利用连接发送到接收端
	// file.Read()
	i := 0
	for {
		// 读取文件内容
		rn, err := file.Read(buf)
		if err != nil {
			// io.EOF错误表示文件读取完毕
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}
		// 发送到接收端
		if _, err := udpConn.Write(buf[:rn]); err != nil {
			log.Fatalln(err)
		}
		i++
	}
	// 文件发送完成
	log.Println("transfer times:", i)
	log.Println("file send complete.")
	time.Sleep(time.Second) // 可能存在接收端还没有接收完，发送端提前退出的情况
}
