package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	for {
		time.Sleep(time.Second)
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://127.0.0.1:1234/hello", nil)

		if err != nil {
			fmt.Println("NewRequest", err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("client.Do:", err)
		}
		fmt.Println("resp:", resp)
	}
}
