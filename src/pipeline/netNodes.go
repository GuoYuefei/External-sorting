package pipeline

import (
	"net"
)

//执行这个的一般是数据所在地的机器，需要其他机器来访问它
//接受连接并向网络流中发送数据
func NetWorkSink(addr string,in <-chan int) {
	listener, err := net.Listen("tcp", addr)
	if err != nil{
		panic(err)
	}

	go func() {
		defer listener.Close()
		//接收连接，conn继承了ReaderWriter Closer等
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		//将数据写入Writer
		WriterSinkBuf(conn, in)
	}()
}

//从网络流中接收数据
func NetWorkSource(addr string) <-chan int {
	out := make(chan int, 1024)
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		ch := ReaderSource(conn,-1)
		for v := range ch {				//不断的从ch中取数据知道ch关闭
			out <- v
		}
		close(out)
	}()
	return out
}
