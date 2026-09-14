package extensions

import (
	"errors"
	"testing"
)

func TestLifecycleTransitionMatrix(t *testing.T) {
	states := []State{StateDiscovered, StateStarting, StateReady, StateDegraded, StateStopping, StateStopped, StateFailed}
	allowed := map[[2]State]bool{
		{StateDiscovered, StateStarting}: true,
		{StateStarting, StateReady}:      true,
		{StateStarting, StateFailed}:     true,
		{StateStarting, StateStopping}:   true,
		{StateReady, StateDegraded}:      true,
		{StateReady, StateStopping}:      true,
		{StateReady, StateFailed}:        true,
		{StateDegraded, StateReady}:      true,
		{StateDegraded, StateStopping}:   true,
		{StateDegraded, StateFailed}:     true,
		{StateStopping, StateStopped}:    true,
		{StateStopping, StateFailed}:     true,
		{StateStopped, StateStarting}:    true,
		{StateFailed, StateStarting}:     true,
		{StateFailed, StateStopped}:      true,
	}
	for _, from := range states {
		for _, to := range states {
			err := ValidateTransition(from, to)
			if allowed[[2]State{from, to}] && err != nil {
				t.Fatalf("allowed transition %s -> %s rejected: %v", from, to, err)
			}
			if !allowed[[2]State{from, to}] && !errors.Is(err, ErrInvalidTransition) {
				t.Fatalf("invalid transition %s -> %s returned %v", from, to, err)
			}
		}
	}
	for _, pair := range [][2]State{{"unknown", StateReady}, {StateReady, "unknown"}} {
		if !errors.Is(ValidateTransition(pair[0], pair[1]), ErrInvalidTransition) {
			t.Fatalf("unknown state transition accepted: %#v", pair)
		}
	}
}

func TestValidState(t *testing.T) {
	for _, state := range []State{StateDiscovered, StateStarting, StateReady, StateDegraded, StateStopping, StateStopped, StateFailed} {
		if !validState(state) {
			t.Fatalf("valid state rejected: %s", state)
		}
	}
	if validState("") {
		t.Fatal("empty state accepted")
	}
}
