package main

import (
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// chopstick 筷子  美[ˈtʃɑːpstɪk]
type ChopstickLearn struct{ sync.Mutex }

// Philosopher 代表哲学家 	美[fəˈlɑːsəfər]
type PhilosopherLearn struct {
	// 哲学家名称
	name string
	// 左手一只筷子和右手一只筷子
	leftChopstick, rightChopstick *ChopstickLearn
	status                        string
}

// 无休止的进餐和冥想
// 吃完睡（冥想、打坐），睡完吃
// 可以调整吃睡的时间来增加或减少抢夺叉子的机会
func (p *PhilosopherLearn) dineLearn() {
	for {
		markLearn(p, "冥想")
		randomPauseLearn(10)

		markLearn(p, "饿了")
		p.leftChopstick.Lock() // 先尝试拿起左手边的筷子
		markLearn(p, "拿起左手筷子")
		p.rightChopstick.Lock()
		markLearn(p, "拿起右手筷子，开始用膳")
		randomPauseLearn(10)

		p.rightChopstick.Unlock() // 先尝试放下右手边的筷子
		p.leftChopstick.Unlock()  // 再尝试放下左手边的筷子
	}
}

// 随机暂停一段时间
func randomPauseLearn(max int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max)))
}

// 显示此哲学家的状态
func markLearn(p *PhilosopherLearn, action string) {
	fmt.Printf("%s开始%s\n", p.name, action)
	p.status = fmt.Sprintf("%s开始%s\n", p.name, action)
}

func main() {
	go http.ListenAndServe("localhost:8972", nil)

	// 哲学家数量
	count := 5

	// 创建5根筷子
	Chopsticks := make([]*ChopstickLearn, count)

	for i := 0; i < count; i++ {
		Chopsticks[i] = new(ChopstickLearn)
	}

	//
	names := []string{color.RedString("孔子"), color.MagentaString("庄子"), color.CyanString("墨子"), color.GreenString("孙子"), color.WhiteString("老子")}

	// 创建哲学家，分配给他们左右手边的叉子，领他们坐到圆餐桌上
	Philosophers := make([]*PhilosopherLearn, count)
	for i := 0; i < count; i++ {
		Philosophers[i] = &PhilosopherLearn{
			name:           names[i],
			leftChopstick:  Chopsticks[i],
			rightChopstick: Chopsticks[(i+1)%count],
		}
		go Philosophers[i].dineLearn()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println("退出中... 每个哲学家的状态：")
	for _, p := range Philosophers {
		fmt.Print(p.status)
	}
}
