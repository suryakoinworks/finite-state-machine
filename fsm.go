package fsm

import (
	"errors"
	"sync"

	"golang.org/x/exp/slices"
)

type (
	Transitionable interface {
		SetState(state string)
		GetState() string
		BeforeTransition(machine Transitionable, action Action)
		AfterTransition(machine Transitionable)
	}

	Transition struct {
		To   string
		From []string
	}

	Action struct {
		From string
		To   string
	}

	Machine struct {
		State string
	}

	finiteStateMachine struct {
		lock        sync.Mutex
		machine     Transitionable
		states      []string
		transitions []Transition
	}
)

func (m *Machine) SetState(state string) {
	m.State = state
}

func (m *Machine) GetState() string {
	return m.State
}

func (m *Machine) BeforeTransition(machine Transitionable, action Action) {
}

func (m *Machine) AfterTransition(machine Transitionable) {
}

/**
 * - machine: represent of managed object
 * - states: represent of available places/states/transitions
 * - transitions: represent of transition rules
 *
 * Initiation will return error when the initial state is not valid
 */
func NewFSM(machine Transitionable, states []string, transtitions []Transition) (*finiteStateMachine, error) {
	fsm := finiteStateMachine{
		lock:        sync.Mutex{},
		machine:     machine,
		states:      states,
		transitions: transtitions,
	}

	if !fsm.validateInitiation() {
		return nil, errors.New("invalid initial state")
	}

	return &fsm, nil
}

func (f *finiteStateMachine) AvailableStates() []string {
	return f.states
}

func (f *finiteStateMachine) Actions() []Action {
	actions := make([]Action, 0, len(f.transitions))
	for _, i := range f.transitions {
		for _, f := range i.From {
			actions = append(actions, Action{From: f, To: i.To})
		}
	}

	return actions
}

func (f *finiteStateMachine) GetCurrentState() string {
	f.lock.Lock()
	defer f.lock.Unlock()

	return f.machine.GetState()
}

func (f *finiteStateMachine) Do(state string) error {
	if !slices.Contains(f.states, state) {
		return errors.New("invalid state")
	}

	for _, t := range f.transitions {
		if t.To == state {
			if !slices.Contains(t.From, f.GetCurrentState()) {
				return errors.New("invalid transition")
			}

			f.machine.BeforeTransition(f.machine, Action{From: f.GetCurrentState(), To: state})

			f.lock.Lock()
			defer f.lock.Unlock()

			f.machine.SetState(state)
			f.machine.AfterTransition(f.machine)

			return nil
		}
	}

	return errors.New("invalid transition")
}

func (f *finiteStateMachine) validateInitiation() bool {
	if slices.Contains(f.states, f.GetCurrentState()) {
		return true
	}

	return false
}
