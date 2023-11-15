package main

import (
	"fmt"
	"log"
	"sync"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-03 8:35
 */

func TCPClientPool() {
	log.Println("客户端开始连接...")
	host := "localhost"
	port := 8888

	// 1.建立连接池
	pool, err := NewTcpPool(fmt.Sprintf("%s:%d", host, port), PoolConfig{
		InitConnNum: 4,
		MaxIdleNum:  20,
		Factory:     &TcpConnFactory{},
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(pool, len(pool.idleList))

	// 2.复用连接池中的连接
	wg := sync.WaitGroup{}
	clientNum := 18
	wg.Add(clientNum)
	for i := 0; i < clientNum; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			conn, err := pool.Get() // 获取连接
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(conn)
			err = pool.Put(conn) // 放回连接
			if err != nil {
				log.Println(err)
				return
			}
		}(&wg)
	}
	wg.Wait()

	// 释放连接池
	err = pool.Release()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(pool, pool.idleList)
}
