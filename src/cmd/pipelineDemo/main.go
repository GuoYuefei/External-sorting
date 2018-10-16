package main

import (
	"bufio"
	"fmt"
	"os"
	"pipeline"
)

//生成一个文件，作为排序的源文件
func main() {
	filename := "small.in"
	size := 2000000
	bufSize := 100
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	p := pipeline.RandomSource(size)

	pipeline.WriterSinkBufSize(file, p, bufSize)

	file, err = os.Open(filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	p = pipeline.ReaderSource(bufio.NewReader(file),-1)
	count := 0
	for v := range p {
		count++
		fmt.Println(count,v)
	}

}

func mergeDemo() {
	p := pipeline.InMemorySort(
		pipeline.ArraySource(7,5,6,1,3,9,10,15,12))
	//for {
	//	if num, ok := <-p; ok {
	//		fmt.Println(num)
	//	} else {
	//		fmt.Println("over")
	//		break
	//	}
	//}
	q := pipeline.ArraySource(3,4,21,7,6,8,2,0)
	q = pipeline.InMemorySort(q)
	//for v := range p {				//发送方写了close了
	//	fmt.Println(v)
	//}
	fmt.Println("-----------------")

	w := pipeline.Merge(p,q)
	for v := range w {
		fmt.Println(v)
	}
}