/**
 * <p>Description: (外部排序和内部排序) </>
 * @author lizhi_duan
 * @date 2018/10/29 22:17
 * @version 1.0
 */
package pipeline

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
)

/**
 * <p>Description: (排序原数组数组) </p>
 * @author lizhi_duan
 * @date 2018/10/29 22:23
 */
//作为返回值的channel上加上操作符号，限定操作权限，只出不进或者只进不出
//当前返回值channel中，代表返回一个只出不进的channel，是针对使用者来说的
//可以弹出的从箭头来看，箭头指向使用者，代表channel是将之传递给使用者的
func ArraySource(in ...int) <-chan int {
	out := make(chan int)
	//并发，开一个goroutine
	//操作channel,向channel中set值
	//使用channel就是配合goroutine使用的，所以给channel set值的时候开启goroutine
	go func() {
		for _, v := range in {
			out <- v
		}
		//这里有明显的结束标记，所以结束后关闭channel，否则使用者遍历时会造成死锁
		//造成 fatal error: all goroutines are asleep - deadlock!
		close(out)
	}()
	return out
}

/**
 * <p>Description: (排序) </p>
 * @author lizhi_duan
 * @date 2018/10/30 6:46
 */
//该函数参数集中的in，作为参数，类型是一个channel，只出不进的channel
//如果向该channel中set值，编译期会报错（invalid operation: in <- v (send to receive-only type <-chan int)）
//receive-only针对使用者来说的，使用者只能从该channel中接收值，不能set值
func InMemorySort(in <-chan int) <-chan int {
	dest := make(chan int)
	go func() {
		var source []int
		//将原要排序的数组从channel中取出，放到内存中
		for i := range in {
			source = append(source, i)
		}
		//对数组进行排序
		sort.Ints(source)
		//将排完序的数组放回到channel中
		for _, v := range source {
			dest <- v
		}
		//关闭channel
		close(dest)
	}()
	return dest
}

/**
 * <p>Description: (归并节点) </p>
 * @author lizhi_duan
 * @date 2018/10/30 7:17
 */
func Merge(int1, int2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		//检测两个节点中是否含有数据
		//如果其中一个已经没有数据了那么就不需要比较了，只从另一个中获取就可以了
		//排序接口中当排完序才会往channel中set值
		v1, ok1 := <-int1
		v2, ok2 := <-int2
		for ok1 || ok2 {
			//检测channel2，如果channel2没有数据了，或者channel2中的数据大于channel1中的数据，那么获取channel1中的值
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				//更新channel1的状态和值
				v1, ok1 = <-int1
			} else {
				out <- v2
				//更新channel的状态和值
				v2, ok2 = <-int2
			}
		}
		close(out)
	}()
	return out
}

/**
 * <p>Description: (从reader中读取信息) </p>
 * @author lizhi_duan
 * @date 2018/10/30 21:57
 */
func ReadSource(reader io.Reader) <-chan int {
	out := make(chan int)
	go func() {
		buffer := make([]byte, 8)
		for {
			n, err := reader.Read(buffer)
			//n>0代表还有剩余数据，err不为nil代表buffer中未读满8个字节
			//先判断n则是因为如果剩余字节不到8个，也将剩余的字节全部读到
			if n > 0 {
				//字节序包括：大端序和小端序
				//而所谓大字节序（big endian），便是指其“最高有效位（most significant byte）”落在低地址上的存储方式
				//而对于小字节序（little endian）来说就正好相反了，它把“最低有效位（least significant byte）”放在低地址上
				u := binary.BigEndian.Uint64(buffer)
				out <- int(u)
			}
			if err != nil {
				break
			}
		}
		close(out)
	}()
	return out
}

/**
 * <p>Description: (写数据) </p>
 * @author lizhi_duan
 * @date 2018/10/30 22:24
 */
func WriteSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

/**
 * <p>Description: (随机生成元数据) </p>
 * @author lizhi_duan
 * @date 2018/10/30 22:31
 */
func RandomSource(count int) <-chan int {
	out := make(chan int)
	//给channel设置值的时候不放到goroutine内，同样造成deadlock，原理目前不明
	//for i := 0; i < count; i++ {
	//	r := rand.Int()
	//	fmt.Print(r, " ")
	//	out <- r
	//}
	//close(out)
	go func() {
		for i := 0; i < count; i++ {
			r := rand.Int()
			fmt.Print(r, " ")
			out <- r
		}
		close(out)
	}()

	return out
}
