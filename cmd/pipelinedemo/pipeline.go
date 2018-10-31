package main

import (
	"bufio"
	"fmt"
	"go-channel-study/pipeline"
	"os"
)

/**
 * <p>Description: (一句话描述该方法的作用) </p>
 * @author lizhi_duan
 * @date 2018/10/29 22:24
 */
func main() {
	const (
		filename = "big.in"
		count    = 100000000
	)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	source := pipeline.RandomSource(count)
	//此处writeSink的第一个参数是一个writer接口，传递file对象也可以的原因在于File已经实现了Writer接口
	//所以可以作为Writer类型传递和使用
	//pipeline.WriteSink(file, source)
	//此处的bufio实现了buffer缓冲区，增加了写入速度
	pipeline.WriteSink(bufio.NewWriter(file), source)
	//读取二进制文件
	open, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer open.Close()
	//readSource := pipeline.ReadSource(open)
	//此处的bufio实现了buffer缓冲区，增加了写入速度
	readSource := pipeline.ReadSource(bufio.NewReader(open))
	fmt.Println()
	fmt.Println("读取文件：")
	for i := range readSource {
		fmt.Print(i, " ")
	}
}

func mergeDemo() {
	p := pipeline.Merge(
		pipeline.InMemorySort(pipeline.ArraySource(2, 8, 5, 7, 1)),
		pipeline.InMemorySort(pipeline.ArraySource(3, 2, 5, 12, 10)))
	//死循环遍历channel，判断channel关闭后结束循环
	//for {
	//	//如果不判断channel是否关闭，那么拿到值永远是0，即类型的初值
	//	//i := <-p
	//	//fmt.Println(i)
	//	if i, ok := <-p; ok {
	//		fmt.Println(i)
	//	} else {
	//		break
	//	}
	//}
	//使用range遍历channel，没有之后结束遍历
	for v := range p {
		fmt.Println(v)
	}
}
