package main

// IStrategy 定义一个策略类
type IStrategy interface {
	do(int, int) int
}

// 策略实现：加
type add struct {
}

func (*add) do(a, b int) int {
	return a + b
}

// 策略实现：减
type reduce struct {
}

func (*reduce) do(a, b int) int {
	return a - b
}

// Operator 具体策略的执行者
type Operator struct {
	strategy IStrategy
}

// setStrategy 设置策略
func (Operator *Operator) setStrategy(strategy IStrategy) {
	Operator.strategy = strategy
}

// calculate 调用策略中的方法
func (Operator *Operator) calculate(a, b int) int {
	return Operator.strategy.do(a, b)
}
