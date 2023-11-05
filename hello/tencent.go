package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Book 书
type Book struct {
	name   string
	writer string
	Intro  []byte
}

// saleBok 卖书
func SaleBook() {
	book := Book{
		name:   "xxx",
		writer: "yyy",
		Intro:  []byte("zzz"),
	}

	book_str, _ := json.Marshal(book)
	fmt.Println(string(book_str))
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
