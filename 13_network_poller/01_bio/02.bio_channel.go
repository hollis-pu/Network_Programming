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

// BIOChannel channel的阻塞
func BIOChannel() {
	// 0.初始化数据
	wg := sync.WaitGroup{}
	ch := make(chan struct{}) // IO channel

	// 1.模拟读，体会读的阻塞状态
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println(time.Now().Format("03:04:05.000"), "start read.")
		content := <-ch // IO Read
		log.Println(time.Now().Format("03:04:05.000"), "content:", content)
	}(&wg)

	// 2.模拟写
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		time.Sleep(time.Second) // 阻塞时长
		ch <- struct{}{}
	}(&wg)

	wg.Wait()
}
