package main

import "fmt"

// Mediator interface defines the contract for a mediator
type Mediator interface {
	Notify(sender Component, event string)
}

// ConcreteMediator implements the Mediator interface and coordinates communication between Components
type ConcreteMediator struct {
	components []Component
}

// Notify broadcasts an event to all components except the sender
func (m *ConcreteMediator) Notify(sender Component, event string) {
	for _, component := range m.components {
		if component != sender {
			component.Receive(event)
		}
	}
}

// Component interface defines the contract for a component
type Component interface {
	Send(event string)
	Receive(event string)
}

// ConcreteComponent implements the Component interface
type ConcreteComponent struct {
	mediator Mediator
}

// Send sends an event to the mediator
func (c *ConcreteComponent) Send(event string) {
	c.mediator.Notify(c, event)
}

// Receive receives an event from the mediator
func (c *ConcreteComponent) Receive(event string) {
	fmt.Printf("Received event: %s\n", event)
}

func main() {
	// Create a new ConcreteMediator instance
	mediator := &ConcreteMediator{}

	// Create two ConcreteComponent instances and set their mediator to the ConcreteMediator instance
	component1 := &ConcreteComponent{mediator: mediator}
	component2 := &ConcreteComponent{mediator: mediator}

	// Add the components to the ConcreteMediator's list of components
	mediator.components = append(mediator.components, component1, component2)

	// Send an event from one component and ensure the other component receives it
	component1.Send("Hello, world!")
}
