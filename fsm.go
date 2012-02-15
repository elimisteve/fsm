package fsm

import "fmt"

const (
  OnEntry = "ON_ENTRY"
  OnExit  = "ON_EXIT"
  Default = "DEFAULT"
)

type Transition struct {
  from   string
  event  string
  to     string
  action string
}

type Delegate interface {
  StateMachineCallback(action string, args []interface{})
}

type StateMachine struct {
  delegate     Delegate
  transitions  []Transition
  currentState *Transition
}

func NewStateMachine(delegate Delegate, transitions ...Transition) StateMachine {
  return StateMachine{delegate: delegate, transitions: transitions, currentState: &transitions[0]}
}

func (m *StateMachine) Process(event string, args ...interface{}) {
  trans := m.findTransMatching(m.currentState.from, event)
  if trans == nil {
    trans = m.findTransMatching(m.currentState.from, Default)
  }

  if trans == nil {
    panic(fmt.Sprintf("state machine error: cannot find transition for event [%s] when in state [%s]\n", event, m.currentState.from))
  }

  changing_states := trans.from != trans.to

  if changing_states {
    m.runAction(m.currentState.from, OnExit, args)
  }

  if trans.action != "" {
    m.delegate.StateMachineCallback(trans.action, args)
  }

  if changing_states {
    m.runAction(trans.to, OnEntry, args)
  }

  m.currentState = m.findState(trans.to)
}

func (m *StateMachine) findTransMatching(fromState string, event string) *Transition {
  for _, v := range m.transitions {
    if v.from == fromState && v.event == event {
      return &v
    }
  }
  return nil
}

func (m *StateMachine) runAction(state string, event string, args []interface{}) {
  if trans := m.findTransMatching(state, event); trans != nil && trans.action != "" {
    m.delegate.StateMachineCallback(trans.action, args)
  }
}

func (m *StateMachine) findState(state string) *Transition {
  for _, v := range m.transitions {
    if v.from == state {
      return &v
    }
  }
  return nil
}
