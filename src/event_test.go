package eventdispatcher

import (
	"testing"
)

func TestPropagationIsFalse(t *testing.T) {
	event := NewEvent()
	if event.IsPropagationStopped() {
		t.Errorf("Expected %v got %v", false, event.IsPropagationStopped())
	}
}
