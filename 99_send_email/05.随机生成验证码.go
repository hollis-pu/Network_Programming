package main

import (
	"fmt"
	"math/rand"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-05 10:24
 */

func main() {
	code := generateRandomCode1()
	fmt.Println(code)
}

func generateRandomCode1() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 指定随机种子
	min := 100000                                        // 最小的六位数
	max := 999999                                        // 最大的六位数
	code := min + r.Intn(max-min+1)
	return fmt.Sprintf("%06d", code) // 将随机数格式化为六位数的字符串
}
