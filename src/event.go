package eventdispatcher

import (
	_ "fmt"
)

type Event struct {
	propagationStopped bool
}

func NewEvent() *Event {
	return &Event{propagationStopped: false}
}

func (this *Event) IsPropagationStopped() bool {
	return this.propagationStopped
}

func (this *Event) StopPropagation() {
	this.propagationStopped = true
}

type EventCallback struct {
	call       func(event *Event)
	priority   int
	subscriber EventSubscriber
	name       string
	id	int
}

func NewEventCall(n string, c func(*Event), p int) *EventCallback {
	return &EventCallback{name: n, call: c, priority: p}
}

// define a type so that we can implement a sort
// the default sort is in descending order (see Less where comparison is done using >=)
// change >= to < for an ascending order or create a different type with an anonymous field
// EventCallbacks
type EventCallbacks []*EventCallback

func (calls EventCallbacks) Len() int {
	return len(calls)
}

func (calls EventCallbacks) Less(i, j int) bool {
	return calls[i].priority > calls[j].priority
}

func (calls EventCallbacks) Swap(i, j int) {
	calls[i], calls[j] = calls[j], calls[i]
}
