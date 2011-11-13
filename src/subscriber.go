package eventdispatcher

type EventSubscriber interface {
	GetSubscribedEvents() EventCallbacks
}

// EventCalls
