package pipeline

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

//前期试验用，数据来源使用console输入
func ArraySource(a ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

//内部排序
func InMemorySort(in <-chan int) <-chan int {
	out := make(chan int, 1024)

	go func() {
		a := make([]int,0,16)
		for v := range in {
			a = append(a, v)
		}
		//日志打印
		fmt.Println("Read done:", time.Now().Sub(startTime))

		sort.Ints(a)
		//日志打印
		fmt.Println("InMemorySort done:", time.Now().Sub(startTime))

		for _, v := range a {
			out <- v
		}
		close(out)
	}()

	return out

}

//两两归并
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
		fmt.Println("Merge done:", time.Now().Sub(startTime))
	}()

	return out
}

//file其实就是一种io，取名不取filesource这样有具体含义的
//其实就是在读数据
func ReaderSource(reader io.Reader, chunkSize int64) <-chan int {
	out := make(chan int, 1024)				//给一个1024的缓存空间，有利于速度  不用收一个用一个

	go func() {
		buffer := make([]byte,8)
		var bytesRead int64 = 0
		for {
			//读取完毕会有一个ROF的错误
			n, err := reader.Read(buffer)
			bytesRead += int64(n)
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chunkSize != -1 && bytesRead >= chunkSize) {
				break
			}
		}
		close(out)
	}()

	return out
}

//将in中数据写入 io.Writer
//writer未必就是文件写入流，可以是网络流哦   为集群式外排做准备
func WriterSink(writer io.Writer,in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

func WriterSinkFromBuf(writer bufio.Writer,in <-chan int) {

	WriterSink(&writer, in)
	writer.Flush()
}

func WriterSinkBufSize(writer io.Writer,in <-chan int, bufSize int) {
	w := bufio.NewWriterSize(writer,bufSize)

	WriterSinkFromBuf(*w, in)

}

func WriterSinkBuf(writer io.Writer,in <-chan int) {
	w := bufio.NewWriter(writer)

	WriterSinkFromBuf(*w,in)

}

//用于随机生成数据
func RandomSource(count int) <-chan int {
	out := make(chan int,512)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

//将n块chan int归并，返回归并后的chan int
func MergeN(inputs ...<-chan int) <-chan int {
	//如果只有一块就不用归并
	if len(inputs) == 1 {
		return inputs[0]
	}

	m := len(inputs) / 2
	p := Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))

	//PrintChanI(p,100)
	//递归至一块，然后两两归并到一块并返回，之后不断递归，直至所有的chan全部并归到一块
	return p
}

////功能与上同，没有使用递归
//func MergeM(inputs ...<-chan int) <-chan int {
//	for i := 0; i < len(inputs) / 2; i++ {
//
//	}
//}

//从通道p中打印num个数据 用于调试
func PrintChanI(p <-chan int,num int) {
	count := 1
	for v := range p {
		fmt.Println(count,v)
		count++
		if count > num {
			break
		}
	}
}







