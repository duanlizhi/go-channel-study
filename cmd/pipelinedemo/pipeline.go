package main

import (
	"fmt"
	"go-channel-study/pipeline"
)

/**
 * <p>Description: (一句话描述该方法的作用) </p>
 * @author lizhi_duan
 * @date 2018/10/29 22:24
 */
func main() {
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
