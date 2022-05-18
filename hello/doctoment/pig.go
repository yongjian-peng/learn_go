package main

import (
	"fmt"
	"math/rand"
)

const (
	win            = 25  // 猪游戏中获胜分数
	gamesPerSeries = 100 // 每个系列要模拟的游戏数量
)

// 一个分数包括每个玩家在前一回合中累计的分数
type score struct {
	player, opponent, thisTurn int
}

// 如果回合结束，如果中的player和opponent字段应该交换，因为现在轮到其他玩家了
type action func(current score) (result score, tureIsOver bool)

// 函数roll和stay每个返回一对值，它们匹配action类型签名，这些 action函数定义了pig的规则
func roll(s score) (score, bool) {
	outcome := rand.Intn(6) + 1 // [1, 6] 中的随机整数
	// 如果 roll 值为1 则放弃thisTrue得分 玩家的 角色互换，否则将滚动值添加到 thisTrue
	if outcome == 1 {
		return score{s.opponent, s.player, 0}, true
	}
	return score{s.player, s.opponent, outcome + s.thisTurn}, false
}

// 返回停留的（result trueIsover）结果
func stay(s score) (score, bool) {
	return score{s.opponent, s.player + s.thisTurn, 0}, true
}

// 策略为任何给定的分数选择一个动作
type strategy func(score) action

// stayAtk 返回一个策略，该策略滚动直到 thisTrue 至少为 k 然后停留
func stayStK(k int) strategy {
	return func(s score) action {
		if s.thisTurn >= k {
			return stay
		}
		return roll
	}
}

// play 模拟猪游戏并返回获胜者
// 通过调用 action 来模拟更新 score直到一个玩家达到100分。每个 action 都是通过
// 调用 strategy 与当前玩家关联函数来选择的
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

// roundRobin 在每对策略之间模拟一系列游戏 并计算获胜
// 每种策略都在玩其他策略 gamesPerSeries 时间
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
	gamesPerStrategy := gamesPerSeries * (len(strategies) - 1)
	return wins, gamesPerStrategy
}

// ratioString 接收一个整数值列表返回一个字符串，该字符串列出
// 每个值及其所有值总和的百分比
// 例如：ratios(1, 2, 3) = "1/6 (16.7%), 2/6 (33.3%), 3/6 (50.0%)"
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

// 模拟循环赛，然后打印每个策略的输赢记录
func main() {
	strategies := make([]strategy, win)
	for k := range strategies {
		strategies[k] = stayStK(k + 1)
	}
	wins, games := roundRobin(strategies)

	for k := range strategies {
		fmt.Printf("Wins, losses staying at k = % 4d: %s\n", k+1,
			ratioString(wins[k], games-wins[k]))
	}
}
