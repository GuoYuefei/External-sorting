package main

import "fmt"

func main() {
	ch := make(chan int, 1)
	s := tt(ch)
	s.name = "lala1"
	for {
		if v, ok := <-ch; ok {
			fmt.Println(s.name, s.id, v)
		} else {
			break
		}
	}
}

type ss struct {
	name string
	id   int
}

func tt(ch chan int) ss {
	var s ss
	for i := 0; i < 10; i++ {
		go func(ii int) {
			s = ss{"lala", ii}
			ch <- ii //闭包
		}(i)
	}
	return s
}
