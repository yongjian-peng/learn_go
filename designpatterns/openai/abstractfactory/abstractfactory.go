package main

import (
	"fmt"
	"runtime"
)

// 在这个例子中，我们首先定义了一个抽象工厂接口GUIFactory，它有两个抽象方法CreateButton和CreateLabel，用于创建按钮和标签。

// 然后我们定义了两个具体的工厂WindowsFactory和MacOSFactory，它们分别实现了GUIFactory接口，并提供了自己的实现来创建Windows和macOS的界面元素。

// 我们还定义了四个产品接口，其中Button和Label是抽象的，WindowsButton、WindowsLabel、MacOSButton和MacOSLabel是具体的。

// 每个具体的产品都实现了其相应的接口方法Paint，用于在操作系统上绘制相应的界面元素。

// 最后，在客户端代码中，我们检查当前操作系统

// Abstract Factory interface
type GUIFactory interface {
	CreateButton() Button
	CreateLabel() Label
}

// Concrete factory for Windows
type WindowsFactory struct{}

func (w *WindowsFactory) CreateButton() Button {
	return &WindowsButton{}
}

func (w *WindowsFactory) CreateLabel() Label {
	return &WindowsLabel{}
}

// Concrete factory for macOS
type MacOSFactory struct{}

func (m *MacOSFactory) CreateButton() Button {
	return &MacOSButton{}
}

func (m *MacOSFactory) CreateLabel() Label {
	return &MacOSLabel{}
}

// Abstract product 1
type Button interface {
	Paint() string
}

// Abstract product 2
type Label interface {
	Paint() string
}

// Concrete product 1 for Windows
type WindowsButton struct{}

func (w *WindowsButton) Paint() string {
	return "Painting a Windows button"
}

// Concrete product 2 for Windows
type WindowsLabel struct{}

func (w *WindowsLabel) Paint() string {
	return "Painting a Windows label"
}

// Concrete product 1 for macOS
type MacOSButton struct{}

func (m *MacOSButton) Paint() string {
	return "Painting a macOS button"
}

// Concrete product 2 for macOS
type MacOSLabel struct{}

func (m *MacOSLabel) Paint() string {
	return "Painting a macOS label"
}

// Client code
func main() {
	// Determine the operating system
	var factory GUIFactory
	if runtime.GOOS == "windows" {
		factory = &WindowsFactory{}
	} else if runtime.GOOS == "darwin" {
		factory = &MacOSFactory{}
	} else {
		fmt.Println("Unsupported operating system")
		return
	}

	// Create the products using the selected factory
	button := factory.CreateButton()
	label := factory.CreateLabel()

	// Paint the products
	fmt.Println(button.Paint())
	fmt.Println(label.Paint())
}
