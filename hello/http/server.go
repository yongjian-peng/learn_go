package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", HelloServer)

	http.ListenAndServe(":1234", nil)
}

func HelloServer(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println(req.Form)
	fmt.Println("path", req.URL.Path)
	fmt.Println("This is xc test")
	fmt.Fprintf(resp, "Hello xc")
}
