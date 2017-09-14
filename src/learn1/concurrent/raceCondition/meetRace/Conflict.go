package main

import (
"fmt"
"runtime"
"sync"
)
var (
	count int32
	wg    sync.WaitGroup
)

func main() {
	/**
	返回的结果有时候是 2 有时候是 4 但是大部分是 4
	当循环的次数变大时，还是稳定的结果为主
	 */
	wg.Add(2)
	go incCount()
	go incCount()
	wg.Wait()
	fmt.Println(count)
}

func incCount() {
	defer wg.Done()
	for i := 0; i < 200; i++ {
		value := count
		//runtime.Gosched()是让当前goroutine暂停的意思，退回执行队列，让其他等待的goroutine运行
		runtime.Gosched()
		value++
		count = value
	}
}

