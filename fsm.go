// Finite State Machines
//
// Suitable for use anytime, anywhere. Bring it.
// Here is the basic API:
//
//     sm := NewStateMachine(&delegate,
//
//       Transition{ From: "locked",    Event: "coin",     To: "unlocked",  Action: "token_inc" },
//       Transition{ From: "locked",    Event: OnEntry,                     Action: "enter" },
//       Transition{ From: "locked",    Event: Default,    To: "locked",    Action: "default" },
//
//       Transition{ From: "unlocked",  Event: "turn",     To: "locked",    },
//       Transition{ From: "unlocked",  Event: OnExit,                      Action: "exit" },
//
//       )
//
//     sm.Process("coin")
//     sm.Process("turn", optionalArg, ...)
//     sm.Process("break")
//
// For a more complete usage, see the test file.
package fsm

import "errors"
import "fmt"

const (
	OnEntry = "ON_ENTRY"
	OnExit  = "ON_EXIT"
	Default = "DEFAULT"
)

type Transition struct {
	From   string
	Event  string
	To     string
	Action string
}

// `action' corresponds to what's in a Transition
type Delegate interface {
	StateMachineCallback(action string, args []interface{})
}

type StateMachine struct {
	delegate     Delegate
	transitions  []Transition
	currentState *Transition
}

// Use this in conjunction with Transition literals, keeping
// in mind that To may be omitted for actions, and Action may
// always be omitted. See the overview above for an example.
func NewStateMachine(delegate Delegate, transitions ...Transition) StateMachine {
	return StateMachine{delegate: delegate, transitions: transitions, currentState: &transitions[0]}
}

func (m *StateMachine) Process(event string, args ...interface{}) error {
	trans := m.findTransMatching(m.currentState.From, event)
	if trans == nil {
		trans = m.findTransMatching(m.currentState.From, Default)
	}

	if trans == nil {
		return errors.New(fmt.Sprintf("state machine error: cannot find transition for event [%s] when in state [%s]\n", event, m.currentState.From))
	}

	changing_states := trans.From != trans.To

	if changing_states {
		m.runAction(m.currentState.From, OnExit, args)
	}

	if trans.Action != "" {
		m.delegate.StateMachineCallback(trans.Action, args)
	}

	if changing_states {
		m.runAction(trans.To, OnEntry, args)
	}

	m.currentState = m.findState(trans.To)

	return nil
}

func (m *StateMachine) findTransMatching(fromState string, event string) *Transition {
	for _, v := range m.transitions {
		if v.From == fromState && v.Event == event {
			return &v
		}
	}
	return nil
}

func (m *StateMachine) runAction(state string, event string, args []interface{}) {
	if trans := m.findTransMatching(state, event); trans != nil && trans.Action != "" {
		m.delegate.StateMachineCallback(trans.Action, args)
	}
}

func (m *StateMachine) findState(state string) *Transition {
	for _, v := range m.transitions {
		if v.From == state {
			return &v
		}
	}
	return nil
}
