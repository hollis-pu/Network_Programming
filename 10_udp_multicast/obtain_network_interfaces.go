package main

import (
	"fmt"
	"log"
	"net"
)

/**
* Description:
	获取所有网络接口
* @Author Hollis
* @Create 2023-11-15 23:34
*/

func main() {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	// 遍历并打印每个接口的名称
	for _, iface := range interfaces {
		fmt.Println(iface.Name)
	}
}

//以太网
//LetsTAP
//本地连接* 1
//本地连接* 2
//VMware Network Adapter VMnet1
//VMware Network Adapter VMnet8
//WLAN
//Loopback Pseudo-Interface 1
