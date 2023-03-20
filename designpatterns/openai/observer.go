package main

// 在这个示例中，我们定义了一个 Subject 接口和一个 Observer 接口，它们分别代表事件主题和观察者。然后，我们实现了一个具体的 ConcreteSubject 类型和一个具体的 ConcreteObserver 类型，它们实现了 Subject 接口和 Observer 接口的方法。
//
//在 main 函数中，我们创建了一个具体的主题 subject，并添加了两个具体的观察者 observer1 和 observer2，然后设置主题的状态并通知所有观察者。接着，我们移除了一个观察者 observer1，并再次设置主题的状态并通知所有观察者。
//
//观察者模式可以让我们轻松地实现事件系统，当一个对象的状态发生变化时，它可以通知所有观察者进行
import "fmt"

// Subject is the interface that defines the subject to be observed
type Subject interface {
	Attach(observer Observer)
	Detach(observer Observer)
	Notify()
	GetState() int
	SetState(state int)
}

// Observer is the interface that defines the observer
type Observer interface {
	Update(subject Subject)
}

// ConcreteSubject is a concrete implementation of the Subject interface
type ConcreteSubject struct {
	state     int
	observers []Observer
}

// Attach adds an observer to the subject's observers list
func (s *ConcreteSubject) Attach(observer Observer) {
	s.observers = append(s.observers, observer)
}

// Detach removes an observer from the subject's observers list
func (s *ConcreteSubject) Detach(observer Observer) {
	for i, o := range s.observers {
		if o == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notify notifies all the subject's observers when a state change occurs
func (s *ConcreteSubject) Notify() {
	for _, o := range s.observers {
		o.Update(s)
	}
}

// GetState returns the subject's current state
func (s *ConcreteSubject) GetState() int {
	return s.state
}

// SetState sets the subject's current state and notifies all observers of the change
func (s *ConcreteSubject) SetState(state int) {
	s.state = state
	s.Notify()
}

// ConcreteObserver is a concrete implementation of the Observer interface
type ConcreteObserver struct {
	name string
}

// Update prints out the observer's name and the subject's current state
func (o *ConcreteObserver) Update(subject Subject) {
	fmt.Printf("Observer %s received the update. New state is %d\n", o.name, subject.GetState())
}

func main() {
	// Create a new subject
	subject := &ConcreteSubject{}

	// Create and attach two observers
	observer1 := &ConcreteObserver{"Observer 1"}
	observer2 := &ConcreteObserver{"Observer 2"}
	subject.Attach(observer1)
	subject.Attach(observer2)

	// Set the subject's state and notify all observers
	subject.SetState(1)

	// Detach an observer and set the subject's state again
	subject.Detach(observer1)
	subject.SetState(2)
}
