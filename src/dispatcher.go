package eventdispatcher

import (
_	"fmt"
	"sort"
)

type Dispatcher struct {
	callbacks map[string]EventCallbacks
	count int
}

func NewEventDispatcher() *Dispatcher {
	this := &Dispatcher{}
	this.callbacks = make(map[string]EventCallbacks)
	this.count = 0
	return this
}

func (this *Dispatcher) AddSubscriber(subscriber EventSubscriber) {
	for _, call := range subscriber.GetSubscribedEvents() {
		this.AddListener(call, subscriber)
	}

	this.sortSubscribers()
}

func (this *Dispatcher) AddListener(call *EventCallback, subscriber EventSubscriber) {
	calls, ok := this.callbacks[call.name]
	if !ok {
		calls = make(EventCallbacks, 0)
	}
	call.subscriber = subscriber
	call.id = this.count
	this.callbacks[call.name] = append(calls, call)
	this.count += 1
}

func (this *Dispatcher) Dispatch(eventName string, event *Event) {
	callbacks, ok := this.callbacks[eventName]
	if !ok {
		return
	}

	for _, callback := range callbacks {
		callback.call(event)
		if event.IsPropagationStopped() {
			break
		}
	}
}

func (this *Dispatcher) RemoveSubscriber(subscriber EventSubscriber) {
    events := subscriber.GetSubscribedEvents()
	for _, callback := range events{
		this.RemoveListener(callback)
	}
	this.sortSubscribers()
}

func (this *Dispatcher) sortSubscribers() {
	for _, callbacks := range this.callbacks {
		sort.Sort(callbacks)
	}
}

func (this *Dispatcher) RemoveListener(callback *EventCallback) {
	callbacks, ok := this.callbacks[callback.name]
	if !ok {
		return
	}

	// TODO: replace this with a more efficient search
	for i, c := range callbacks{
		if c == callback{
	       this.callbacks[callback.name] = append(callbacks[:i], callbacks[i+1:]...)
	       break
		}
	}
}