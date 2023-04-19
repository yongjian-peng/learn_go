package main

import "fmt"

// 在这个例子中，我们首先定义了一个抽象工厂接口ShapeFactory，它有两个抽象方法CreateCurvedShape和CreateStraightShape，用于创建曲线形状和直线形状的产品。

// 然后我们定义了两个具体的工厂RoundedRectangleFactory和EllipseFactory，它们分别实现了ShapeFactory接口，并提供了自己的实现来创建产品。

// 我们还定义了四个产品接口，其中CurvedShape和StraightShape是抽象的，RoundedRectangle、Ellipse、StraightLine和DiagonalLine是具体的。

// 每个具体的产品都实现了其相应的接口方法Draw，用于绘制产品。

// 最后，我们在客户端代码中使用这些工厂和产品来创建和绘制不同的形状

// Abstract Factory interface
type ShapeFactory interface {
	CreateCurvedShape() CurvedShape
	CreateStraightShape() StraightShape
}

// Concrete factory 1
type RoundedRectangleFactory struct{}

func (r *RoundedRectangleFactory) CreateCurvedShape() CurvedShape {
	return &RoundedRectangle{}
}

func (r *RoundedRectangleFactory) CreateStraightShape() StraightShape {
	return &StraightLine{}
}

// Concrete factory 2
type EllipseFactory struct{}

func (e *EllipseFactory) CreateCurvedShape() CurvedShape {
	return &Ellipse{}
}

func (e *EllipseFactory) CreateStraightShape() StraightShape {
	return &DiagonalLine{}
}

// Abstract product 1
type CurvedShape interface {
	Draw() string
}

// Abstract product 2
type StraightShape interface {
	Draw() string
}

// Concrete product 1
type RoundedRectangle struct{}

func (r *RoundedRectangle) Draw() string {
	return "Drawing a rounded rectangle"
}

// Concrete product 2
type Ellipse struct{}

func (e *Ellipse) Draw() string {
	return "Drawing an ellipse"
}

// Concrete product 3
type StraightLine struct{}

func (s *StraightLine) Draw() string {
	return "Drawing a straight line"
}

// Concrete product 4
type DiagonalLine struct{}

func (d *DiagonalLine) Draw() string {
	return "Drawing a diagonal line"
}

// Client code
func main() {
	// Create the first factory
	factory1 := &RoundedRectangleFactory{}

	// Create the products using the first factory
	curvedShape1 := factory1.CreateCurvedShape()
	straightShape1 := factory1.CreateStraightShape()

	// Draw the products
	fmt.Println(curvedShape1.Draw())
	fmt.Println(straightShape1.Draw())

	// Create the second factory
	factory2 := &EllipseFactory{}

	// Create the products using the second factory
	curvedShape2 := factory2.CreateCurvedShape()
	straightShape2 := factory2.CreateStraightShape()

	// Draw the products
	fmt.Println(curvedShape2.Draw())
	fmt.Println(straightShape2.Draw())
}
