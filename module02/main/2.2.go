package main

import (
	"log"
	"strconv"
	"sync"
	"time"
)

var fmeng = false

func producer(threadId int, wg *sync.WaitGroup, ch chan string) {
	count := 0
	for !fmeng {
		time.Sleep(time.Second)
		count++
		data := strconv.Itoa(threadId) + "---" + strconv.Itoa(count)
		log.Println("producer", data)
		ch <- data
	}
	wg.Done()
}

func consumer(wg *sync.WaitGroup, ch chan string) {
	for data := range ch {
		time.Sleep(time.Second)
		log.Println("consumer", data)
	}
	wg.Done()
}

func main() {
	// 多生产者，多消费者
	chanStream := make(chan string, 10)
	// 生产者和消费者计数器
	producers := new(sync.WaitGroup)
	consumers := new(sync.WaitGroup)

	for i := 0; i < 3; i++ {
		producers.Add(1)
		go producer(i, producers, chanStream)
	}

	for j := 0; j < 2; j++ {
		consumers.Add(1)
		go consumer(consumers, chanStream)
	}

	// 制造超时
	go func() {
		time.Sleep(time.Second * 3)
		fmeng = true
	}()

	producers.Wait()

	// 生产完关闭ch
	close(chanStream)

	consumers.Wait()
}
