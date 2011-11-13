package eventdispatcher

import (
	"testing"
)

type testSubscriber struct {
	EventSubscriber
	callbacks EventCallbacks
}

func NewTestSubscriber() *testSubscriber {
	calls := make(EventCallbacks, 0)
	calls = append(calls, NewEventCall("OnFoo", OnFooCallback, 0))
    calls = append(calls, NewEventCall("OnBoo", OnBooCallback, 1))
	return &testSubscriber{callbacks: calls}
}

func (this *testSubscriber) GetSubscribedEvents() EventCallbacks {
	return this.callbacks
}

func OnFooCallback(e *Event) {
}

func OnBooCallback(e *Event) {
}

type testSubscriber2 struct {
	EventSubscriber
	callbacks EventCallbacks
}

func NewTestSubscriber2() *testSubscriber2 {
	calls := make(EventCallbacks, 0)
	calls = append(calls, NewEventCall("OnFoo", OnFooCallback2, 1))
	calls = append(calls, NewEventCall("OnBoo", OnBooCallback2, 0))
	return &testSubscriber2{callbacks: calls}
}

func (this *testSubscriber2) GetSubscribedEvents() EventCallbacks {
	return this.callbacks
}

func OnFooCallback2(e *Event) {
}

func OnBooCallback2(e *Event) {
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

	if fooCallbacks[0].call != OnFooCallback2{
		t.Errorf("Expected callback %v got %v in OnFoo\n", OnFooCallback2, fooCallbacks[0])
	}

	booCallbacks := dispatcher.callbacks["OnBoo"]
	realSize = len(booCallbacks)

	if realSize != 1 {
		t.Errorf("Expected callbacks size %v got %v\n", 1, realSize)
	}

	if booCallbacks[0].call != OnBooCallback2{
		t.Errorf("Expected callback %v got %v in OnBoo\n", OnBooCallback2, booCallbacks[0])
	}
}
