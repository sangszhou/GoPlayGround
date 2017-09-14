package ch8

import "time"

func main()  {

}

func spinner(delay time.Duration)  {

}

func fib(x int) int {
	if x < 2 {
		return x
	}

	return fib(x-1) + fib(x-2)
}
