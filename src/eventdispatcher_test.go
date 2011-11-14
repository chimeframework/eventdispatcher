package eventdispatcher

import (
	"fmt"
	"testing"
)

func TestPropagationIsFalse(t *testing.T) {
	event := NewEvent()
	if event.IsPropagationStopped() {
		t.Errorf("Expected %v got %v", false, event.IsPropagationStopped())
	}
}

type testSubscriber struct {
	EventSubscriber
	callbacks EventCallbacks
	name      string
}

func NewTestSubscriber() *testSubscriber {
	calls := make(EventCallbacks, 0)
	calls = append(calls, NewEventCallback("OnFoo", OnFooCallback, 0))
	calls = append(calls, NewEventCallback("OnBoo", OnBooCallback, 1))
	return &testSubscriber{callbacks: calls, name: "TestSubscriber"}
}

func (this *testSubscriber) GetSubscribedEvents() EventCallbacks {
	return this.callbacks
}

func (this *testSubscriber) Print() {
	fmt.Printf("PASS %v\n", this.name)
}

func (this *testSubscriber2) Print() {
	fmt.Printf("PASS %v\n", this.id)
}

func OnFooCallback(receiver EventSubscriber, e Eventer) {
	this := receiver.(*testSubscriber)
	this.Print()
}

func OnBooCallback(this EventSubscriber, e Eventer) {
}

type testSubscriber2 struct {
	EventSubscriber
	callbacks EventCallbacks
	id        int
}

func NewTestSubscriber2() *testSubscriber2 {
	calls := make(EventCallbacks, 0)
	calls = append(calls, NewEventCallback("OnFoo", OnFooCallback2, 1))
	calls = append(calls, NewEventCallback("OnBoo", OnBooCallback2, 0))
	return &testSubscriber2{callbacks: calls, id: 123456}
}

func (this *testSubscriber2) GetSubscribedEvents() EventCallbacks {
	return this.callbacks
}

func OnFooCallback2(receiver EventSubscriber, e Eventer) {
	this := receiver.(*testSubscriber2)
	this.Print()
}

func OnBooCallback2(this EventSubscriber, e Eventer) {
}

func TestAddSubscriber(t *testing.T) {
	dispatcher := NewEventDispatcher()
	dispatcher.AddSubscriber(NewTestSubscriber())
	dispatcher.AddSubscriber(NewTestSubscriber2())

	realSize := len(dispatcher.callbacks["OnFoo"])

	if realSize != 2 {
		t.Errorf("Expected callbacks size %v got %v\n", 2, realSize)
	}

	realSize = len(dispatcher.callbacks["OnBoo"])

	if realSize != 2 {
		t.Errorf("Expected callbacks size %v got %v\n", 2, realSize)
	}
}

func TestPrioritySortSubscribers(t *testing.T) {
	dispatcher := NewEventDispatcher()
	dispatcher.AddSubscriber(NewTestSubscriber())
	dispatcher.AddSubscriber(NewTestSubscriber2())

	fooCallbacks := dispatcher.callbacks["OnFoo"]

	if fooCallbacks[0].call != OnFooCallback2 {
		t.Errorf("Expected callback %v got %v\n", OnFooCallback, fooCallbacks[0])
	}

	if fooCallbacks[1].call != OnFooCallback {
		t.Errorf("Expected callback %v got %v\n", OnFooCallback, fooCallbacks[1])
	}

	booCallbacks := dispatcher.callbacks["OnBoo"]

	if booCallbacks[0].call != OnBooCallback {
		t.Errorf("Expected callback %v got %v\n", OnBooCallback, fooCallbacks[0])
	}

	if booCallbacks[1].call != OnBooCallback2 {
		t.Errorf("Expected callback %v got %v\n", OnBooCallback, fooCallbacks[0])
	}
}

func TestRemoveSubscribers(t *testing.T) {
	dispatcher := NewEventDispatcher()
	subscriber1 := NewTestSubscriber()
	subscriber2 := NewTestSubscriber2()
	dispatcher.AddSubscriber(subscriber1)
	dispatcher.AddSubscriber(subscriber2)

	dispatcher.RemoveSubscriber(subscriber1)

	fooCallbacks := dispatcher.callbacks["OnFoo"]
	realSize := len(fooCallbacks)
	if realSize != 1 {
		t.Errorf("Expected callbacks size %v got %v\n", 1, realSize)
	}

	if fooCallbacks[0].call != OnFooCallback2 {
		t.Errorf("Expected callback %v got %v in OnFoo\n", OnFooCallback2, fooCallbacks[0])
	}

	booCallbacks := dispatcher.callbacks["OnBoo"]
	realSize = len(booCallbacks)

	if realSize != 1 {
		t.Errorf("Expected callbacks size %v got %v\n", 1, realSize)
	}

	if booCallbacks[0].call != OnBooCallback2 {
		t.Errorf("Expected callback %v got %v in OnBoo\n", OnBooCallback2, booCallbacks[0])
	}
}
func TestDispatch(t *testing.T) {
	dispatcher := NewEventDispatcher()
	dispatcher.AddSubscriber(NewTestSubscriber())
	dispatcher.AddSubscriber(NewTestSubscriber2())
	dispatcher.Dispatch("OnFoo", NewEvent())
}

func TestHasSubscribers(t *testing.T) {
	dispatcher := NewEventDispatcher()
	dispatcher.AddSubscriber(NewTestSubscriber())
	dispatcher.AddSubscriber(NewTestSubscriber2())

	if !dispatcher.HasCallbacks("OnFoo") {
		t.Errorf("Expected subscribers size for OnFoo %v got %v\n", 2, 0)
	}
	if !dispatcher.HasCallbacks("OnBoo") {
		t.Errorf("Expected subscribers size for OnBoo %v got %v\n", 2, 0)
	}
	if dispatcher.HasCallbacks("OnFooBoo") {
		t.Errorf("Expected subscribers size for OnFooBoo %v\n", 0)
	}
}

type testSubscriber3 struct {
	EventSubscriber
	callbacks EventCallbacks
}

func NewTestSubscriber3() *testSubscriber3 {
	calls := make(EventCallbacks, 0)
	calls = append(calls, NewEventCallback("OnProp", OnProp1Callback, 1))
	return &testSubscriber3{callbacks: calls}
}

func (this *testSubscriber3) GetSubscribedEvents() EventCallbacks {
	return this.callbacks
}

func OnProp1Callback(receiver EventSubscriber, e Eventer) {
	e.StopPropagation()
}

type testSubscriber4 struct {
	EventSubscriber
	callbacks EventCallbacks
}

func NewTestSubscriber4() *testSubscriber4 {
	return &testSubscriber4{}
}

func (this *testSubscriber4) GetSubscribedEvents() EventCallbacks {
	return this.callbacks
}
func OnProp2Callback(receiver EventSubscriber, e Eventer) {
	fmt.Printf("FAIL! OnProp2Callback should never be called. OnProp1Callback should have stopped the propagation")
}

func TestEventPropagation(t *testing.T) {
	dispatcher := NewEventDispatcher()
	sub4 := NewTestSubscriber4()
	calls := make(EventCallbacks, 0)
	calls = append(calls, NewEventCallback("OnProp", func(receiver EventSubscriber, e Eventer) {
		t.Errorf("OnProp2Callback should never be called. OnProp1Callback should have stopped the propagation")
	}, 0))
	sub4.callbacks = calls

	dispatcher.AddSubscriber(sub4)
	dispatcher.AddSubscriber(NewTestSubscriber3())
	dispatcher.Dispatch("OnProp", NewEvent())
}

type testSubscriber5 struct {
	EventSubscriber
	callbacks EventCallbacks
}

type CustomEvent struct {
	Event
	name string
	id   int
}

func NewCustomEvent() *CustomEvent {
	return &CustomEvent{name: "custom event", id: 123456}
}

func NewTestSubscriber5() *testSubscriber5 {
	return &testSubscriber5{}
}

func (this *testSubscriber5) GetSubscribedEvents() EventCallbacks {
	return this.callbacks
}

func TestCustomEvent(t *testing.T) {
	dispatcher := NewEventDispatcher()
	sub := NewTestSubscriber5()
	calls := make(EventCallbacks, 0)

	callback := func(receiver EventSubscriber, e Eventer) {
		event := e.(*CustomEvent)
		if event.name != "custom event" {
			t.Errorf("Expected callback event name 'custom event' got %v\n", event.name)
		}
	}

	calls = append(calls, NewEventCallback("OnCustomEventDispatch", callback, 0))
	sub.callbacks = calls
	dispatcher.AddSubscriber(sub)
	dispatcher.Dispatch("OnCustomEventDispatch", NewCustomEvent())
}
