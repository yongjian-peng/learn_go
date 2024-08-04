package main

import "fmt"

/**
qcount : 当前 channel 中存在多少个元素；
dataqsize: 当前 channel 能存放的元素容量；
buf: channel 中用于存放元素的环形缓冲区
elemsize: channel 元素类型的大小；
closed: 标识 channel 是否关闭；
elemtype: channel 元素类型；
sendx: 发送元素进入环形缓冲区的 index；
recvx: 接受元素所处的环形缓冲区的 index;
recvq: 因接受而陷入阻塞的协程队列;
sendq: 因发送而陷入阻塞的协程队列;

sudog: 用于包装协程的节点
g: goroutine, 协程；
next: 队列中的下一个节点；
prev: 队列中的前一个节点；
elem: 读取/写入 channel 的数据的容器；
c: 标识当前 sudog 交互的 chan;

	Channel应用场景
		消息传递、消息过滤
		信号广播
		并发控制，同步与异步
		事件订阅与广播等
	数据结构 hchan
	写流程
		写时存在阻塞读协程（同步发送）
		写时无阻塞读协程且环形缓冲区仍有空间（异步发送）
		写时无阻塞读协程且环形缓冲区无空间（阻塞发送）
	读流程
		读时有阻塞的写协程（同步接收）
		读时无阻塞写协程且缓冲区有元素（异步接收）
		读时无阻塞写协程且缓冲区无元素（阻塞接收）
*/

// 接受返回值的时候，第二个值是bool类型
// 代表是否成功的从 chan 中读取到一个值，如果第二个参数是false,
// chan 已经被 close 而且 chan 中没有缓存的数据。
// 这个时候，第一个值是零值，所以，如果从 chan 读取到一个零值，可能是 sender 真正的发送零值。
// 也可能是 closed 的并且没有缓存元素产生的零值
func main() {
	// 创建一个带有缓冲区的通道
	ch := make(chan int, 1)

	// 发送数据到通道
	ch <- 0

	// 从通道接受的值
	value, ok := <-ch
	// 检查是否成功的接受到值
	if ok {
		fmt.Println("接受到的值：", value)
	} else {
		fmt.Println("未接受到值")
	}
}
