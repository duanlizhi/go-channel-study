/**
 * <p>Description: (一句话描述一下该文件的作用) </>
 * @author lizhi_duan
 * @date 2018/10/28 10:51
 * @version 1.0
 */
package main

import (
	"fmt"
)

func main() {
	//开一个goroutine
	//这里非线程，而是一个协程
	//在操作系统中，有两个重要的概念：一个是进程、一个是线程
	//当我们运行一个程序的时候，比如你的IDE或者QQ等，操作系统会为这个程序创建一个进程，这个进程包含了运行这个程序所需的各种资源，可以说它是一个容器，
	// 是属于这个程序的工作空间，比如它里面有内存空间、文件句柄、设备和线程等等。、

	//概念	           说明
	//进程	一个程序对应一个独立程序空间
	//线程	一个执行空间，一个进程可以有多个线程
	//逻辑处理器	执行创建的goroutine，绑定一个线程
	//调度器	Go运行时中的，分配goroutine给不同的逻辑处理器
	//全局运行队列	所有刚创建的goroutine都会放到这里
	//本地运行队列	逻辑处理器的goroutine队列
	ch :=make(chan string)
	for i:=0; i<5000; i++ {
		go printHello(i,ch)
	}
	for {
		//取出管道里的字符串信息，赋值给msg变量
		msg:= <- ch
		fmt.Println(msg)
	}
}

func printHello(i int, ch chan string)  {
	for {
		//将生成的字符串放入管道里边，按照箭头符号顺序
		ch <- fmt.Sprintf("hello world! from goroutine %d\n",i)
	}

}
