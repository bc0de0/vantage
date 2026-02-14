package core

import "vantage/core/state"

// StateMachine controls legal operation phase transitions.
type StateMachine struct {
	phase state.OperationPhase
}

func NewStateMachine(initial state.OperationPhase) (*StateMachine, error) {
	if err := initial.Validate(); err != nil {
		return nil, err
	}
	return &StateMachine{phase: initial}, nil
}

func (m *StateMachine) Phase() state.OperationPhase {
	return m.phase
}

func (m *StateMachine) Advance() bool {
	next, ok := m.phase.Next()
	if !ok {
		return false
	}
	m.phase = next
	return true
}
