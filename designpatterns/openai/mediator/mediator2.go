package main

import "fmt"

// 在这个示例中，我们定义了一个 Mediator2 接口和一个 Colleague2 接口 其中 Mediator2 接口定义了一个 SendMessage2 方法，用于向所有的同事发送消息，而 Colleague2 接口定义了 Send2 和 Receive 方法，用于发送和接收消息
// 然后，我们具体实现了一个具体的中介者 ConcreteMediator2 它维护了一个同事的列表，并实现了 SendMessage2 方法来向所有同事发送消息
// 接着，我们实现了两个具体的同事， ConcreteColleagueA 和 ConcreteColleagueB 它们都有一个指向中介者的指针

// Mediator2 interface
type Mediator2 interface {
	SendMessage2(message string, sender Colleague2)
}

// Colleague2 interface
type Colleague2 interface {
	Send2(message string)
	Receive2(message string)
}

// ConcreteMediator2 is a concrete implementation of the Mediator interface
type ConcreteMediator2 struct {
	colleagues []Colleague2
}

// AddColleague add a colleague to the mediator2
func (m *ConcreteMediator2) AddColleague2(colleague Colleague2) {
	m.colleagues = append(m.colleagues, colleague)
}

// SendMessage2 send a message to all colleagues except the sender
func (m *ConcreteMediator2) SendMessage2(message string, sender Colleague2) {
	for _, colleague := range m.colleagues {
		if colleague != sender {
			colleague.Receive2(message)
		}
	}
}

// ConcreteColleagueC is a concrete implementation of the Colleague interface
type ConcreteColleagueC struct {
	mediator2 Mediator2
}

// SetMediator sets the colleague's mediator2
func (c *ConcreteColleagueC) SetMediator(mediator Mediator2) {
	c.mediator2 = mediator
}

// Send2 sends a message to the mediator2
func (c *ConcreteColleagueC) Send2(message string) {
	c.mediator2.SendMessage2(message, c)
}

// Receive receives a message form the mediator2
func (c *ConcreteColleagueC) Receive2(message string) {
	fmt.Printf("ConcreteColleagueC received message: %s\n", message)
}

type ConcreteColleagueD struct {
	mediator2 Mediator2
}

// ConcreteColleagueD is a concrete implementation of the Colleague interface
func (c *ConcreteColleagueD) SetMediator(mediator Mediator2) {
	c.mediator2 = mediator
}

func (c *ConcreteColleagueD) Send2(message string) {
	c.mediator2.SendMessage2(message, c)
}

func (c *ConcreteColleagueD) Receive2(message string) {
	fmt.Printf("ConcreteColleagueD received message: %s\n", message)
}

func main() {
	mediator := &ConcreteMediator2{}

	colleagueC := &ConcreteColleagueC{}
	colleagueD := &ConcreteColleagueD{}

	mediator.AddColleague2(colleagueC)
	mediator.AddColleague2(colleagueD)

	colleagueC.SetMediator(mediator)
	colleagueD.SetMediator(mediator)

	colleagueC.Send2("Hello from colleague A")
	colleagueD.Send2("Hello from colleague B")
}
