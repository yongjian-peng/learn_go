package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/coreos/etcd/clientv3"
	recipe "github.com/coreos/etcd/contrib/recipes"
)

var (
	addr      = flag.String("addr", "http://127.0.0.1:12379", "http://127.0.0.1:12379")
	queueName = flag.String("name", "my-test-queue", "my-test-queue")
)

func main() {
	flag.Parse()

	// 解析etcd地址
	endpoints := strings.Split(*addr, ",")

	// 创建etcd的client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 创建/获取队列
	queue := recipe.NewQueue(cli, *queueName)

	// 从命令行读取命令
	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "push": // 加入队列
			if len(items) != 2 {
				fmt.Println("must set value to push")
				continue
			}
			queue.Enqueue(items[1]) // 入队
		case "pop": // 从队列弹出
			v, err := queue.Dequeue() // 出队
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("v=>", v) // 输出队列到元素
		case "quit", "exit": // 退出
			return
		default:
			fmt.Println("unknown action")
		}
	}
}
