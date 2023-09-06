package fsm

import (
	"errors"
	"sync"

	"golang.org/x/exp/slices"
)

type (
	transitionable interface {
		SetState(state string)
		GetState() string
		BeforeTransition(machine transitionable, transition string)
		AfterTransition(machine transitionable)
	}

	Transition struct {
		To   string
		From []string
	}

	Machine struct {
	}

	finiteStateMachine struct {
		lock        sync.Mutex
		machine     transitionable
		states      []string
		transitions []Transition
	}
)

func (f *Machine) BeforeTransition(machine transitionable, state string) {
}

func (f *Machine) AfterTransition(machine transitionable) {
}

/**
 * - machine: represent of managed object
 * - states: represent of available places/states/transitions
 * - transitions: represent of transition rules
 *
 * Initiation will return error when the initial state is not valid
 */
func NewFSM(machine transitionable, states []string, transtitions []Transition) (*finiteStateMachine, error) {
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
			if !slices.Contains(t.From, f.machine.GetState()) {
				return errors.New("invalid transition")
			}

			f.lock.Lock()
			defer f.lock.Unlock()

			f.machine.BeforeTransition(f.machine, state)
			f.machine.SetState(state)
			f.machine.AfterTransition(f.machine)

			return nil
		}
	}

	return errors.New("invalid transition")
}

func (f *finiteStateMachine) validateInitiation() bool {
	if slices.Contains(f.states, f.machine.GetState()) {
		return true
	}

	return false
}
