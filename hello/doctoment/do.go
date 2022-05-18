// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// 猪是一个两人游戏，玩一个 6 面骰子。每转一圈，你可以滚动或停留。
// 如果您掷出 1，您将失去本回合的所有分数并将传球传给您的对手。任何其他掷骰都会为您的回合得分增加其价值。
// 如果您留下，您的回合得分将被添加到您的总得分中，并将传球传给您的对手。
// 第一个达到 100 分的人获胜。

// 除了当前回合累积的积分外，该score类型还存储当前和对方玩家的分数。
package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
)

const (
	win            = 25  // 猪游戏中的获胜分数
	gamesPerSeries = 100 // 每个系列要模拟的游戏数量
)

var ResultSuccessNum int64

// 一个分数包括每个玩家在前一回合中累积的分数，
// 以及当前玩家本回合得分。
type score struct {
	player, opponent, thisTurn int
}

// 一个动作随机转换为一个结果分数
// 如果回合结束，结果中的player和opponent字段score应该交换，因为现在轮到其他玩家了。
type action func(current score) (result score, turnIsOver bool)

// 函数roll和stay每个返回一对值。它们也匹配action类型签名。这些 action函数定义了 Pig 的规则。
// roll 返回模拟掷骰的 (result, turnIsOver) 结果。
// 如果roll值为1，则放弃thisTurn得分，玩家的
// 角色交换。 否则，将滚动值添加到 thisTurn。
func roll(s score) (score, bool) {
	outcome := rand.Intn(6) + 1 // [1, 6] 中的随机整数
	atomic.AddInt64(&ResultSuccessNum, 1)
	if outcome == 1 {
		return score{s.opponent, s.player, 0}, true
	}
	return score{s.player, s.opponent, outcome + s.thisTurn}, false
}

// stay 返回停留的 (result, turnIsOver) 结果。
// thisTurn 分数被添加到玩家的分数中，并且玩家的角色互换。
func stay(s score) (score, bool) {
	return score{s.opponent, s.player + s.thisTurn, 0}, true
}

// 策略为任何给定的分数选择一个动作。
// Astrategy是一个以 ascore作为输入并返回action要执行的函数。
// （请记住，an action本身就是一个函数。）
type strategy func(score) action

// stayAtK 返回一个策略，该策略滚动直到 thisTurn 至少为 k，然后停留。
// 高阶函数
// 一个函数可以使用其他函数作为参数和返回值。

func stayAtK(k int) strategy {
	return func(s score) action {
		if s.thisTurn >= k {
			return stay
		}
		return roll
	}
}

// play 模拟猪游戏并返回获胜者（0 或 1）
//  模拟游戏
// 我们通过调用 an 来模拟 Pig 游戏action来更新， score直到一个玩家达到 100 分。每一个 action都是通过调用strategy与当前播放器关联的函数来选择的。
func play(strategy0, strategy1 strategy) int {
	strategies := []strategy{strategy0, strategy1}
	var s score
	var turnIsOver bool
	currentPlayer := rand.Intn(2) // 随机决定谁先玩
	for s.player+s.thisTurn < win {
		action := strategies[currentPlayer](s)
		s, turnIsOver = action(s)
		if turnIsOver {
			currentPlayer = (currentPlayer + 1) % 2
		}
	}
	return currentPlayer
}

// roundRobin 在每对策略之间模拟一系列游戏。
// 模拟比赛
// 该roundRobin函数模拟锦标赛并计算获胜。每种策略都在玩其他策略gamesPerSeries时间
func roundRobin(strategies []strategy) ([]int, int) {
	wins := make([]int, len(strategies))
	for i := 0; i < len(strategies); i++ {
		for j := i + 1; j < len(strategies); j++ {
			for k := 0; k < gamesPerSeries; k++ {
				winner := play(strategies[i], strategies[j])
				if winner == 0 {
					wins[i]++
				} else {
					wins[j]++
				}
			}
		}
	}
	gamesPerStrategy := gamesPerSeries * (len(strategies) - 1) // 没有自我游戏
	return wins, gamesPerStrategy
}

// 可变函数声明
// 可变参数函数，例如ratioString采用可变数量的参数。这些参数可作为函数内部的切片使用。
// ratioString 接受一个整数值列表并返回一个字符串，该字符串列出
// 每个值及其占所有值总和的百分比。
// 例如，ratios(1, 2, 3) = "1/6 (16.7%), 2/6 (33.3%), 3/6 (50.0%)"
func ratioString(vals ...int) string {
	total := 0
	for _, val := range vals {
		total += val
	}
	s := ""
	for _, val := range vals {
		if s != "" {
			s += ", "
		}
		pct := 100 * float64(val) / float64(total)
		s += fmt.Sprintf("%d/%d (%0.1f%%)", val, total, pct)
	}
	return s
}

// 仿真结果
// 该main函数定义了100个基本策略，模拟循环赛，然后打印每个策略的输赢记录。

// 在这些策略中，保持在 25 是最好的，但Pig 的最优策略要复杂得多
func main() {
	strategies := make([]strategy, win)
	for k := range strategies {
		strategies[k] = stayAtK(k + 1)
	}
	wins, games := roundRobin(strategies)

	for k := range strategies {
		fmt.Printf("Wins, losses staying at k =% 4d: %s\n",
			k+1, ratioString(wins[k], games-wins[k]))
	}

	fmt.Println(ResultSuccessNum)
}
