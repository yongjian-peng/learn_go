package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"../../selfworkgroup"
	"../../yongworkgroup"
)

func FmtWeb() {

	fmt.Println("fmtWeb")
}

func RunSelf() {
	var wg selfworkgroup.SelfGroup

	wg.Add(selfworkgroup.SelfSignal())

	srv := http.Server{Addr: "127.0.0.1:8082"}

	wg.Add(selfworkgroup.SelfServer(
		func() error {
			fmt.Printf("SelfServer listen at %v\n", srv.Addr)
			err := srv.ListenAndServe()
			fmt.Printf("SelfServer stopped listening with error: %v\n", err)
			if err != http.ErrServerClosed {
				return err
			}
			return nil
		},
		func() error {
			fmt.Println("SelfServer is about to shutdown")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
			defer cancel()

			err := srv.Shutdown(ctx)
			fmt.Printf("SelfServer shutdown with error: %v\n", err)
			return err
		},
	))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("SelfContext canceled")
		cancel()
	}()

	wg.Add(selfworkgroup.SelfContext(ctx))

	err := wg.Run()

	fmt.Printf("SelfWorkgroup run stopped with error: %v\n", err)
}

func RunYong() {
	var wg_yong yongworkgroup.YongGroup

	wg_yong.Add(yongworkgroup.YongSignal())

	srv := http.Server{Addr: "127.0.0.1:8088"}

	wg_yong.Add(yongworkgroup.YongServer(
		func() error {
			fmt.Printf("SelfServer listen at %v\n", srv.Addr)
			err := srv.ListenAndServe()
			fmt.Printf("SelfServer stopped listening with error: %v\n", err)
			if err != http.ErrServerClosed {
				return err
			}
			return nil
		},
		func() error {
			fmt.Println("YongServer is about to shotdown")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
			defer cancel()

			err := srv.Shutdown(ctx)
			fmt.Printf("YongServer shutdown with error: %v\n", err)

			return err
		},
	))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("YongContext canceled")
		cancel()
	}()

	wg_yong.Add(yongworkgroup.YongContext(ctx))

	err := wg_yong.RunYong()

	fmt.Printf("yongWorkGroup run stopped with error: %v\n", err)
}
