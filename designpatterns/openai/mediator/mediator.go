package main

// 在这个示例中，我们定义了一个 Mediator 接口和一个 Colleague 接口，其中 Mediator 接口定义了一个 SendMessage 方法，用于向所有同事发送消息，而 Colleague 接口定义了 Send 和 Receive 方法，用于发送和接收消息。
//
//然后，我们实现了一个具体的中介者 ConcreteMediator，它维护了一个同事列表，并实现了 SendMessage 方法来向所有同事发送消息。
//
//接着，我们实现了两个具体的同事 ConcreteColleagueA 和 ConcreteColleagueB，它们都有一个指向中介者的指针
import "fmt"

// Mediator interface
type Mediator interface {
	SendMessage(message string, sender Colleague)
}

// Colleague interface
type Colleague interface {
	Send(message string)
	Receive(message string)
}

// ConcreteMediator is a concrete implementation of the Mediator interface
type ConcreteMediator struct {
	colleagues []Colleague
}

// AddColleague adds a colleague to the mediator2
func (m *ConcreteMediator) AddColleague(colleague Colleague) {
	m.colleagues = append(m.colleagues, colleague)
}

// SendMessage sends a message to all colleagues except the sender
func (m *ConcreteMediator) SendMessage(message string, sender Colleague) {
	for _, colleague := range m.colleagues {
		if colleague != sender {
			colleague.Receive(message)
		}
	}
}

// ConcreteColleagueA is a concrete implementation of the Colleague interface
type ConcreteColleagueA struct {
	mediator Mediator
}

// SetMediator sets the colleague's mediator2
func (c *ConcreteColleagueA) SetMediator(mediator Mediator) {
	c.mediator = mediator
}

// Send sends a message to the mediator2
func (c *ConcreteColleagueA) Send(message string) {
	c.mediator.SendMessage(message, c)
}

// Receive receives a message from the mediator2
func (c *ConcreteColleagueA) Receive(message string) {
	fmt.Printf("ConcreteColleagueA received message: %s\n", message)
}

// ConcreteColleagueB is a concrete implementation of the Colleague interface
type ConcreteColleagueB struct {
	mediator Mediator
}

// SetMediator sets the colleague's mediator2
func (c *ConcreteColleagueB) SetMediator(mediator Mediator) {
	c.mediator = mediator
}

// Send sends a message to the mediator2
func (c *ConcreteColleagueB) Send(message string) {
	c.mediator.SendMessage(message, c)
}

// Receive receives a message from the mediator2
func (c *ConcreteColleagueB) Receive(message string) {
	fmt.Printf("ConcreteColleagueB received message: %s\n", message)
}

func main() {
	// Create a new mediator2
	mediator := &ConcreteMediator{}

	// Create some colleagues and add them to the mediator2
	colleagueA := &ConcreteColleagueA{}
	colleagueB := &ConcreteColleagueB{}
	mediator.AddColleague(colleagueA)
	mediator.AddColleague(colleagueB)

	// Set the colleagues' mediator2
	colleagueA.SetMediator(mediator)
	colleagueB.SetMediator(mediator)

	// Send messages between colleagues
	colleagueA.Send("Hello from colleague A")
	colleagueB.Send("Hello from colleague B")
}
