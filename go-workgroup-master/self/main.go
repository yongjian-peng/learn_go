package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	// 	"./self2"
	"./self2/web"
)

func main() {
	// web.RunSelf()
	web.RunYong()
	// Create workgroup
	// var wg workgroup.Group
	// // Add function to cancel execution using os signal
	// wg.Add(workgroup.Signal())
	// // Create http server
	// srv := http.Server{Addr: "127.0.0.1:8081"}
	// // Add function to start and stop http server
	// wg.Add(workgroup.Server(
	// 	func() error {
	// 		fmt.Printf("Server listen at %v\n", srv.Addr)
	// 		err := srv.ListenAndServe()
	// 		fmt.Printf("Server stopped listening with error: %v\n", err)
	// 		if err != http.ErrServerClosed {
	// 			return err
	// 		}
	// 		return nil
	// 	},
	// 	func() error {
	// 		fmt.Println("Server is about to shutdown")
	// 		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	// 		defer cancel()

	// 		err := srv.Shutdown(ctx)
	// 		fmt.Printf("Server shutdown with error: %v\n", err)
	// 		return err
	// 	},
	// ))
	// // Create context to cancel execution after 5 seconds
	// ctx, cancel := context.WithCancel(context.Background())
	// go func() {
	// 	time.Sleep(time.Second * 5)
	// 	fmt.Println("Context canceled")
	// 	cancel()
	// }()
	// // Add function to cancel execution using context
	// wg.Add(workgroup.Context(ctx))
	// // Execute each function
	// err := wg.Run()
	// fmt.Printf("Workgroup run stopped with error: %v\n", err)

	// assert(t, nil, g.Run())
	// self2.Self()
	// web.FmtWeb()
	// tr := NewTracker()
	// go tr.Run()
	// _ = tr.Event(context.Background(), "test1")
	// _ = tr.Event(context.Background(), "test2")
	// _ = tr.Event(context.Background(), "test3")
	// time.Sleep(3 * time.Second)
	// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	// defer cancel()
	// tr.Shutdown(ctx)
}

func assert(t *testing.T, want, got error) {
	t.Helper()
	if want != got {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}

func NewTracker() *Tracker {
	return &Tracker{
		ch: make(chan string, 10),
	}
}

// 写的 owener 操作 channel 生命周期
//

type Tracker struct {
	ch   chan string
	stop chan struct{}
}

func (t *Tracker) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Tracker) Run() {
	for data := range t.ch {
		time.Sleep(1 * time.Second)
		fmt.Println(data)
	}
	t.stop <- struct{}{}
}

func (t *Tracker) Shutdown(ctx context.Context) {
	close(t.ch)
	select {
	case <-t.stop:
	case <-ctx.Done():
	}
}
