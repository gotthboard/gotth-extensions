package extensions

type State string

const (
	StateDiscovered State = "discovered"
	StateStarting   State = "starting"
	StateReady      State = "ready"
	StateDegraded   State = "degraded"
	StateStopping   State = "stopping"
	StateStopped    State = "stopped"
	StateFailed     State = "failed"
)

func ValidateTransition(from, to State) error {
	if !validState(from) || !validState(to) {
		return ErrInvalidTransition
	}
	allowed := false
	switch from {
	case StateDiscovered:
		allowed = to == StateStarting
	case StateStarting:
		allowed = to == StateReady || to == StateFailed || to == StateStopping
	case StateReady:
		allowed = to == StateDegraded || to == StateStopping || to == StateFailed
	case StateDegraded:
		allowed = to == StateReady || to == StateStopping || to == StateFailed
	case StateStopping:
		allowed = to == StateStopped || to == StateFailed
	case StateStopped, StateFailed:
		allowed = to == StateStarting || from == StateFailed && to == StateStopped
	}
	if !allowed {
		return ErrInvalidTransition
	}
	return nil
}

func validState(state State) bool {
	switch state {
	case StateDiscovered, StateStarting, StateReady, StateDegraded, StateStopping, StateStopped, StateFailed:
		return true
	default:
		return false
	}
}
