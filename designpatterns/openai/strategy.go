// 在这个示例中，我们定义了一个 Strategy 接口，它有一个 Execute 方法，用于执行特定的策略。然后，我们实现了三个具体的策略：AddStrategy、SubtractStrategy 和 MultiplyStrategy，它们分别实现了 Strategy 接口的 Execute 方法。
//
// 接着，我们定义了一个 Context 类型，它包含当前使用的策略，并提供了 SetStrategy 和 ExecuteStrategy 方法，分别用于设置当前策略和执行当前策略。
//
// 在 main 函数中，我们创建了一个新的 Context 实例，并分别设置了三个具体的策略：加法、减法和乘法，并执行每个策略。
//
// 策略模式可以使代码更加灵活，我们可以根据需要更改策略而不必更改客户端代码，并且可以轻松地添加新的策略。
package main

import "fmt"

// Strategy is the interface that defines the strategy
type Strategy interface {
	Execute(a int, b int) int
}

// AddStrategy is a concrete implementation of the Strategy interface
type AddStrategy struct{}

// Execute adds two integers and returns the result
func (s *AddStrategy) Execute(a int, b int) int {
	return a + b
}

// SubtractStrategy is a concrete implementation of the Strategy interface
type SubtractStrategy struct{}

// Execute subtracts two integers and returns the result
func (s *SubtractStrategy) Execute(a int, b int) int {
	return a - b
}

// MultiplyStrategy is a concrete implementation of the Strategy interface
type MultiplyStrategy struct{}

// Execute multiplies two integers and returns the result
func (s *MultiplyStrategy) Execute(a int, b int) int {
	return a * b
}

// Context contains the current strategy
type Context struct {
	strategy Strategy
}

// SetStrategy sets the current strategy
func (c *Context) SetStrategy(strategy Strategy) {
	c.strategy = strategy
}

// ExecuteStrategy executes the current strategy and returns the result
func (c *Context) ExecuteStrategy(a int, b int) int {
	return c.strategy.Execute(a, b)
}

func main() {
	// Create a new context
	context := &Context{}

	// Set the context's strategy to add and execute it
	context.SetStrategy(&AddStrategy{})
	result := context.ExecuteStrategy(2, 3)
	fmt.Printf("Add result is: %d\n", result)

	// Set the context's strategy to subtract and execute it
	context.SetStrategy(&SubtractStrategy{})
	result = context.ExecuteStrategy(2, 3)
	fmt.Printf("Subtract result is: %d\n", result)

	// Set the context's strategy to multiply and execute it
	context.SetStrategy(&MultiplyStrategy{})
	result = context.ExecuteStrategy(2, 3)
	fmt.Printf("Multiply result is: %d\n", result)
}
