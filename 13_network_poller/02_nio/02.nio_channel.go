package main

import (
	"log"
	"sync"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-16 22:52
 */

// NIOChannel channel的阻塞
func NIOChannel() {
	// 0.初始化数据
	wg := sync.WaitGroup{}
	ch := make(chan struct{ id uint }) // IO channel

	// 1.模拟读，体会读的阻塞状态
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println(time.Now().Format("03:04:05.000"), "start read.")
		var content struct{ id uint }

		// 使用select的default子句，完成非阻塞操作
		select {
		case content = <-ch: // IO Read
		default:
		}

		log.Println(time.Now().Format("03:04:05.000"), "content:", content)
	}(&wg)

	// 2.模拟写
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		time.Sleep(time.Second) // 阻塞时长
		ch <- struct{ id uint }{42}
	}(&wg)

	wg.Wait()
}
