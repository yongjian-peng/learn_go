package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
)

// 一个功能更强大的Once
type Once struct {
	m    sync.Mutex
	done uint32
}

// 传入的函数f有返回值error, 如果初始化失败，需要返回失败的error
// Do 方法会把这个error返回给调用者
func (o *Once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}
	return o.slowDo(f)
}

// 如果还没有初始化
func (o *Once) slowDo(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 { // 双检查，还没有初始化
		err = f()
		if err == nil { // 初始化成功才将标记置为已初始化
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}

func main() {
	var once Once
	var googleConn net.Conn

	err := once.Do(func() error {
		var err error
		googleConn, err = net.Dial("tcp", "baidu.com:80")
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Println("err=>", err)
		return
	}

	_, _ = googleConn.Write([]byte("GET / HTTP/1.1\r\nHost: baidu.com\r\n Accept: */*\r\n\r\n"))
	written, err := io.Copy(os.Stdout, googleConn)
	if err != nil {
		fmt.Println("err=>", err)
		return
	}

	fmt.Println("written=>", written)
}
