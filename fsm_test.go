package fsm

import "testing"

type testTokenMachineContext struct {
  count int
  char rune
  entered int
}

func (ctx *testTokenMachineContext) StateMachineCallback(action string, args []interface{}) {
  switch action {
  case "token_inc":
    ctx.count++
    ctx.char = args[0].(rune)
  case "enter":
    ctx.entered++
  case "exit":
    ctx.entered = 7
  case "default":
    ctx.entered = 88
  }
}

func TestTokenMachine(t *testing.T) {
  var ctx testTokenMachineContext

  tm := NewStateMachine(&ctx,
      Transition{ from: "locked",    event: "coin",     to: "unlocked",  action: "token_inc" },
      Transition{ from: "locked",    event: OnEntry,                     action: "enter" },
      Transition{ from: "locked",    event: Default,    to: "locked",    action: "default" },
      Transition{ from: "unlocked",  event: "turn",     to: "locked",    },
      Transition{ from: "unlocked",  event: OnExit,                      action: "exit" },
      )

  if ! (tm.currentState.from == "locked") { t.Errorf("state machine failure") }
  if ! (ctx.count == 0) { t.Errorf("state machine failure") }
  if ! (ctx.char == 0) { t.Errorf("state machine failure") }

  tm.Process("coin", 'i')
  if ! (tm.currentState.from == "unlocked") { t.Errorf("state machine failure") }
  if ! (ctx.count == 1) { t.Errorf("state machine failure") }
  if ! (ctx.char == 'i') { t.Errorf("state machine failure") }

  tm.Process("turn", 'q')
  if ! (tm.currentState.from == "locked") { t.Errorf("state machine failure") }
  if ! (ctx.count == 1) { t.Errorf("state machine failure") }
  if ! (ctx.entered == 8) { t.Errorf("state machine failure, %d", ctx.entered) }

  tm.Process("random", 'p')
  if ! (tm.currentState.from == "locked") { t.Errorf("state machine failure") }
  if ! (ctx.entered == 88) { t.Errorf("state machine failure, %d", ctx.entered) }
}
