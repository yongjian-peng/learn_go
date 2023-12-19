package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Book 书
type Book struct {
	name   string // 书名
	writer string // 作者
	Intro  []byte // 简介
}

// SaleBook 卖书
func SaleBook() {
	book := Book{
		name:   "xxx",
		writer: "yyy",
		Intro:  []byte("zzz"),
	}
	book_str, _ := json.Marshal(book)
	fmt.Println("book_str", string(book_str))
}

func main() {
	t1 := time.NewTimer(time.Minute / 2)
	t2 := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-t2.C:
			go func() {
				SaleBook()
			}()
		case <-t1.C:
			t2.Stop()
			return
		}
	}
}
