package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

// 由于嵌套调用do里面的lock导致死锁
func ErrorOne() {
	var o sync.Once
	o.Do(func() {
		o.Do(func() {
			fmt.Println("初始化")
		})
	})
}

// 初始化未成功，但是此时once done的值已经是1，默认初始化成功
func ErrTwo() {
	var once sync.Once
	var googleConn net.Conn

	once.Do(func() {
		googleConn, _ = net.Dial("tcp", "google.com:80")
	})
	/*googleConn, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		fmt.Println("err=>", err)
		return
	}*/
	_, _ = googleConn.Write([]byte("GET / HTTP/1.1\r\nHost: google.com\r\n Accept: */*\r\n\r\n"))
	_, _ = io.Copy(os.Stdout, googleConn)
}

func main() {
	//ErrorOne()
	ErrTwo()
}
