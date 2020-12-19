package main

import (
	"log"
)

type ConnLimiter struct {
	concurrentConn int
	bucket         chan int // 即使是计数 也尽量不要用锁，我们用channel
}

// 比如同时只能提供10个服务，那把channel设置10， go语言协程同步是用channel，而不是 全局变量+锁

func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the rate limitation.")
		return false
	}
	cl.bucket <- 1
	log.Printf("Successfully got connection")
	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	c := <-cl.bucket
	log.Printf("New connection coming: %d", c)
}
