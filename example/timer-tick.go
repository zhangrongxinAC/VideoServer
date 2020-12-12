package main

import (
	"fmt"
	"time"
)

func main() {
	//创建定时器，每隔1秒后，定时器就会给channel发送一个事件(当前时间)
	ticker := time.NewTicker(time.Second * 1)

	i := 0
	go func() {
		for { //循环
			<-ticker.C // 定时器触发
			i++
			fmt.Println("i = ", i)

			if i == 5 {
				ticker.Stop() //停止定时器
				fmt.Println("stop timer")
			}
		}
	}() //别忘了()

	//死循环，特地不让main goroutine结束
	for {
	}
}
