package fsm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type machine struct {
	state string
	Machine
}

func (f *machine) GetState() string {
	return f.state
}

func (f *machine) SetState(state string) {
	f.state = state
}

func TestInvalidFsmInitiate(t *testing.T) {
	fsm, err := NewFSM(&machine{}, []string{"open", "close", "half_open"}, []Transition{})

	assert.Nil(t, fsm)
	assert.NotNil(t, err)
}

func TestValidFsmInitiate(t *testing.T) {
	machine := machine{
		state: "open",
	}

	fsm, err := NewFSM(&machine, []string{"open", "close", "half_open"}, []Transition{})

	assert.NotNil(t, fsm)
	assert.Nil(t, err)
}

func TestGetCurrentState(t *testing.T) {
	machine := machine{
		state: "open",
	}

	fsm, err := NewFSM(&machine, []string{"open", "close", "half_open"}, []Transition{})

	assert.NotNil(t, fsm)
	assert.Nil(t, err)
	assert.Equal(t, machine.GetState(), fsm.GetCurrentState())
}

func TestInvalidState(t *testing.T) {
	machine := machine{
		state: "open",
	}

	fsm, err := NewFSM(&machine, []string{"open", "close", "half_open"}, []Transition{})

	assert.NotNil(t, fsm)
	assert.Nil(t, err)
	assert.Error(t, fsm.Do("invalid"))
}

func TestEmptyTransition(t *testing.T) {
	machine := machine{
		state: "open",
	}

	states := []string{"open", "close", "half_open"}

	fsm, err := NewFSM(&machine, states, []Transition{})

	assert.NotNil(t, fsm)
	assert.Nil(t, err)
	assert.Error(t, fsm.Do("open"))
	assert.Equal(t, states, fsm.AvailableStates())
}

func TestInvalidTransition(t *testing.T) {
	machine := machine{
		state: "open",
	}

	fsm, err := NewFSM(&machine, []string{"open", "close", "half_open"}, []Transition{
		{From: []string{"open", "half_open"}, To: "close"},
		{From: []string{"close"}, To: "half_open"},
		{From: []string{"close"}, To: "open"},
	})

	assert.NotNil(t, fsm)
	assert.Nil(t, err)
	assert.Error(t, fsm.Do("half_open"))
}

func TestValidTransition(t *testing.T) {
	machine := machine{
		state: "open",
	}

	fsm, err := NewFSM(&machine, []string{"open", "close", "half_open"}, []Transition{
		{From: []string{"open", "half_open"}, To: "close"},
		{From: []string{"close"}, To: "half_open"},
		{From: []string{"close"}, To: "open"},
	})

	fmt.Println(fsm.Actions())

	assert.NotNil(t, fsm)
	assert.Nil(t, err)
	assert.Nil(t, fsm.Do("close"))
	assert.Equal(t, machine.GetState(), fsm.GetCurrentState())
	assert.Equal(t, machine.GetState(), "close")
	assert.Equal(t, 4, len(fsm.Actions()))
}
