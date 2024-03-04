package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 派出所证明
	var psCertificate sync.Mutex
	// 物业证明
	var propertyCertificate sync.Mutex
	// 额外的资源锁
	var resourceLock sync.Mutex

	var wg sync.WaitGroup
	wg.Add(2) // 需要派出所和物业都处理

	// 派出所处理goroutine
	go func() {
		defer wg.Done() // 派出所处理完成

		// 请求额外的资源锁
		resourceLock.Lock()
		defer resourceLock.Unlock()

		psCertificate.Lock()
		defer psCertificate.Unlock()

		// 检查材料
		time.Sleep(5 * time.Second)

		// 请求物业的证明
		propertyCertificate.Lock()
		propertyCertificate.Unlock()
	}()

	// 物业处理goroutine
	go func() {
		defer wg.Done() // 物业处理完成

		// 请求额外的资源锁
		resourceLock.Lock()
		defer resourceLock.Unlock()

		propertyCertificate.Lock()
		defer propertyCertificate.Unlock()

		// 检查材料
		time.Sleep(5 * time.Second)

		// 请求派出所的证明
		psCertificate.Lock()
		psCertificate.Unlock()
	}()

	wg.Wait()
	fmt.Println("成功完成")
}
