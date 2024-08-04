package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

// Do 传入函数f有返回值error，如果初始化失败，需要返回失败的error
// Do方法会把这个error返回给调用者
func (o *Once) Do(fn func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}
	return o.doSlow(fn)
}

// doSlow 初始化fn函数，成功则标注
func (o *Once) doSlow(fn func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 {
		err = fn()
		if err != nil { // 只有初始化成功才打标注
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}

func main() {
	urls := []string{
		"127.0.0.1:3453",
		"127.0.0.1:9002",
		"127.0.0.1:9003",
		"baidu.com:80",
	}

	var conn net.Conn
	var o Once
	count := 0
	var err error
	for _, url := range urls {
		err := o.Do(func() error {
			count++
			fmt.Printf("初始化%d次\n", count)
			conn, err = net.DialTimeout("tcp", url, time.Second)
			fmt.Println("err=>", err)
			return err
		})
		if err == nil {
			break
		}
		if count == 3 {
			fmt.Println("初始化失败，不再重试")
			break
		}
	}
	if conn != nil {
		_, _ = conn.Write([]byte("GET / HTTP/1.1\nHost: google.com\n Accept: */*\n\n"))
		_, _ = io.Copy(os.Stdout, conn)
	}
}
