package eventdispatcher

import (
	"sort"
)

// EventDispatcher stores a table map of event names to corresponding callbacks.
type EventDispatcher struct {
	callbacks map[string]EventCallbacks
}

// NewEventDispatcher returns a new instance of a EventDispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{callbacks: make(map[string]EventCallbacks)}
}

// AddSubscriber adds an event subscriber and all its subscribed events to the callbacks map.
func (this *EventDispatcher) AddSubscriber(subscriber EventSubscriber) {
	for _, call := range subscriber.GetSubscribedEvents() {
		this.AddCallback(call, subscriber)
	}

	this.sortSubscribers()
}

// AddCallback adds a callback to EventDispatcher.
// It also modifies EventCallback by storing a pointer to the subscriber. This is important
// as when a subscriber is notified about an event, it needs to have a pointer to itself to be
// able to call other methods.
func (this *EventDispatcher) AddCallback(call *EventCallback, subscriber EventSubscriber) {
	calls, ok := this.callbacks[call.name]
	if !ok {
		calls = make(EventCallbacks, 0)
	}
	call.subscriber = subscriber
	this.callbacks[call.name] = append(calls, call)
}

// Dispatch calls all the corresponding callbacks which are subscribed to the event as indicated
// by eventName. The calls are made based on their priority number (0 denotes the lowest priority).
// The event propagation can be stopped by a subscriber by calling StopPropagation on a passed event.
func (this *EventDispatcher) Dispatch(eventName string, event Eventer) {
	callbacks, ok := this.callbacks[eventName]
	if !ok {
		return
	}

	for _, callback := range callbacks {
		callback.call(callback.subscriber, event)
		if event.IsPropagationStopped() {
			break
		}
	}
}

// HasCallbacks returns true if there are callbacks for a given eventName.
func (this *EventDispatcher) HasCallbacks(eventName string) bool {
	_, ok := this.callbacks[eventName]
	return ok
}

// RemoveSubscriber removes all the subscribed events of a subscriber from the callbacks table.
func (this *EventDispatcher) RemoveSubscriber(subscriber EventSubscriber) {
	events := subscriber.GetSubscribedEvents()
	for _, callback := range events {
		this.RemoveCallback(callback)
	}
	this.sortSubscribers()
}

// RemoveCallback unregisters (and removes) a particular callback from the table. 
func (this *EventDispatcher) RemoveCallback(callback *EventCallback) {
	callbacks, ok := this.callbacks[callback.name]
	if !ok {
		return
	}

	// TODO: replace this with a more efficient search
	for i, c := range callbacks {
		if c == callback {
			this.callbacks[callback.name] = append(callbacks[:i], callbacks[i+1:]...)
			break
		}
	}
}

// sortSubscribers sorts all the callbacks based on their priority (highest to lowest).
func (this *EventDispatcher) sortSubscribers() {
	for _, callbacks := range this.callbacks {
		sort.Sort(callbacks)
	}
}

// EventSubscriber is an interface to be implemented by a subscriber that intends to receive
// an event.
type EventSubscriber interface {
	GetSubscribedEvents() EventCallbacks
}

type Eventer interface {
	IsPropagationStopped() bool
	StopPropagation()
}

// Event is a type for a basic Event.
type Event struct {
	propagationStopped bool
}

// NewEvent reaturns a new instance of Event.
func NewEvent() *Event {
	return &Event{propagationStopped: false}
}

// IsPropagationStopped returns true if an event has been stopped from being propagated to low priority subscribers.
func (this *Event) IsPropagationStopped() bool {
	return this.propagationStopped
}

// StopPropagation forces an event to be propagated further to low priority subscribers.
func (this *Event) StopPropagation() {
	this.propagationStopped = true
}

// EventCallback defines a type for a callback. Every callback will have a name, a priority, a callback, as well as
// a pointer to the subscriber itself.
type EventCallback struct {
	name       string
	call       func(EventSubscriber, Eventer)
	priority   int
	subscriber EventSubscriber
}

// NewEventCallback returns an instance of EventCallback.
func NewEventCallback(n string, c func(EventSubscriber, Eventer), p int) *EventCallback {
	return &EventCallback{name: n, call: c, priority: p}
}

// EventCallbacks defines a type so that we can implement a sort.
//
// The default sort is in descending order (see Less where comparison is done using >=). Change >= to < for an 
// ascending order or create a different type with an anonymous field EventCallbacks.
type EventCallbacks []*EventCallback

// Len() returns the length of EventCallbacks
func (calls EventCallbacks) Len() int {
	return len(calls)
}

// Less compares two EventCallback based on their priority.
func (calls EventCallbacks) Less(i, j int) bool {
	return calls[i].priority > calls[j].priority
}

// Swaps two callbacks from EventCallbacks.
func (calls EventCallbacks) Swap(i, j int) {
	calls[i], calls[j] = calls[j], calls[i]
}
