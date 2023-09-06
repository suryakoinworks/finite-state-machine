# Finite State Machine

Finite State Machine is a Golang Library to handle [FSM (Finite State Machine)](https://en.wikipedia.org/wiki/Finite-state_machine) easly

## Install

```shell
go get github.com/suryakoinworks/finite-state-machine
```

## Basic Usage

```go
type machine struct {
	state string
	fsm.Machine
}

func (f *machine) GetState() string {
	return f.state
}

func (f *machine) SetState(state string) {
	f.state = state
}

func main() {
    machine := machine{
		state: "open",
	}

	fsm, err := NewFSM(&machine, []string{"open", "close", "half_open"}, []Transition{
		{From: []string{"open", "half_open"}, To: "close"},
		{From: []string{"close"}, To: "half_open"},
		{From: []string{"close"}, To: "open"},
	})

    fsm.GetCurrentState() // return: open
    fsm.Do("half_open") // return: error (invalid transition)
    fsm.Do("close")
    fsm.GetCurrentState() // return: close
    machine.GetState() // return: close
}
```
