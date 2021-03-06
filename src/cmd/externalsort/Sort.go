package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"pipeline"
)

var line string = "---------------------------------------"

//外部排序实现，主要操作有pipeline包封装
func main() {
	var sourcefilename string = "small.in"
	var desfilename string = "small.out"
	var chunkCount int = 800
	p := createPipeline(sourcefilename, chunkCount)
	//for i := 0; i < 1e10; i++ {
	//
	//}
	//pipeline.PrintChanI(p,100)					//这里输出就已经不正常了，证明不是写文件的问题
	fmt.Println(line)
	writeToFile(p,desfilename)
	printFile(desfilename)
}

//打印文件，filename为文件的全名，亦可是相对路径
func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	p := pipeline.ReaderSource(file, -1)
	pipeline.PrintChanI(p, 100)
}

//将p中数据写入到filename代表的文件中
func writeToFile(p <-chan int, filename string) {
	file, err := os.Create(filename)
	pipeline.WriterSinkBuf(file,p)
	if err != nil {
		panic(err)
	}
	defer file.Close()
}

//这里最主要的函数
//filename代表文件名，用于获取文件信息
//chunkCount代表将文件最初分成多少份进行外部排序
//最后返回一个chan，这个chan中拥有最后的结果
func createPipeline(filename string, chunkCount int) <-chan int {
	var inputs []<-chan int
	fileSize := getFileSize(filename)
	//fmt.Println(fileSize)
	chunkSize := fileSize / int64(chunkCount)

	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		//fmt.Println(chunkSize)
		if err != nil {
			panic(err)
		}
		file.Seek(chunkSize * int64(i),0)
		//获得chan int
		p := pipeline.ReaderSource(bufio.NewReader(file),chunkSize)

		p = pipeline.InMemorySort(p)
		inputs = append(inputs,p)
		//fmt.Println(len(inputs))
		//printChanI(inputs[i],100)			//没问题，归并有问题
		//fmt.Println(i)
	}

	return pipeline.MergeN(inputs...)
}

func getFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}
