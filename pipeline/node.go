/**
 * <p>Description: (外部排序和内部排序) </>
 * @author lizhi_duan
 * @date 2018/10/29 22:17
 * @version 1.0
 */
package pipeline

import "sort"

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
		//go中已经做了处理，过来时两方已经排完序，机制未了解
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
