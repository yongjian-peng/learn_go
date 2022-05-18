// 版权所有 2010 Go 作者。 版权所有。
// 此源代码的使用由 BSD 样式管理
// 可以在 LICENSE 文件中找到的许可证。

package main

import (
	"log"
	"net/http"
	"time"
)

const (
	numPollers     = 2               // 要启动的轮询 goroutine 的数量
	pollInterval   = 1 * time.Second // 轮询每个 URL 的频率
	statusInterval = 1 * time.Second // 多久将状态记录到标准输出
	errTimeout     = 1 * time.Second // 错误后退超时
)

var urls = []string{
	"http://test.co/notify.php",
	"http://api-yong.hk.blueoceantech.co/",
}

// State represents the last-known state of a URL.
type State struct {
	url    string
	status string
}

// StateMonitor 维护了一个映射，该映射存储了正在运行的 URL 的状态
// 轮询，并每隔 updateInterval 纳秒打印当前状态。
// 它返回一个资源状态应该发送到的 chan 状态。
func StateMonitor(updateInterval time.Duration) chan<- State {
	updates := make(chan State)
	urlStatus := make(map[string]string)
	ticker := time.NewTicker(updateInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				logState(urlStatus)
			case s := <-updates:
				urlStatus[s.url] = s.status
			}
		}
	}()
	return updates
}

// logState 打印状态图。
func logState(s map[string]string) {
	log.Println("Current state:")
	for k, v := range s {
		log.Printf(" %s %s", k, v)
	}
}

// Resource 表示该程序要轮询的 HTTP URL
type Resource struct {
	url      string
	errCount int
}

// Poll 对 url 执行 HTTP HEAD 请求
// 并返回 HTTP 状态字符串或错误字符串
func (r *Resource) Poll() string {
	resp, err := http.Head(r.url)
	if err != nil {
		log.Println("Error", r.url, err)
		r.errCount++
		return err.Error()
	}
	r.errCount = 0
	return resp.Status
}

// Sleep 休眠适当的时间间隔（取决于错误状态）
// 在发送资源之前完成。
func (r *Resource) Sleep(done chan<- *Resource) {
	time.Sleep(pollInterval + errTimeout*time.Duration(r.errCount))
	done <- r
}

func Poller(in <-chan *Resource, out chan<- *Resource, status chan<- State) {
	for r := range in {
		s := r.Poll()
		status <- State{r.url, s}
		out <- r
	}
}

func main() {
	// 创建我们的输入和输出通道。
	pending, complete := make(chan *Resource), make(chan *Resource)

	// 启动状态监视器。
	status := StateMonitor(statusInterval)

	// 启动一些 Poller goroutine
	for i := 0; i < numPollers; i++ {
		go Poller(pending, complete, status)
	}

	// 将一些资源发送到待处理队列。
	go func() {
		for _, url := range urls {
			pending <- &Resource{url: url}
		}
	}()

	for r := range complete {
		go r.Sleep(pending)
	}
}
