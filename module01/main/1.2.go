package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
基于 Channel 编写一个简单的单线程生产者消费者模型：

- 队列：
    队列长度 10，队列元素类型为 int
- 生产者：
    每 1 秒往队列中放入一个类型为 int 的元素，队列满时生产者可以阻塞
- 消费者：
    每一秒从队列中获取一个元素并打印，队列为空时消费者阻塞

*/
func main() {
	ch := make(chan int, 10)

	go producer(ch)

	consumer(ch)

}

func consumer(ch chan int) {
	for v := range ch {
		fmt.Println("receiving:", v)
	}
}

func producer(ch chan<- int) {
	rand.Seed(time.Now().UnixNano())
	timer := time.NewTicker(time.Second)
	for {
		t := <-timer.C
		next := rand.Intn(10)
		fmt.Println(t.Format("2006-01-02 15:04:05"), "sending:", next)
		ch <- next
	}
}
