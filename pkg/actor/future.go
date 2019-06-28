package actor

import (
	"errors"
	"sync"
)

// Future promises the result of a long running job.
type Future interface {
	GetResult() (interface{}, error)
}

type futureState uint8

const (
	scheduled futureState = iota
	computed
	thrown
)

// FutureTask is the concrete implementation of a Future.
type FutureTask struct {
	state  futureState
	result interface{}
	done   chan bool
	mux    sync.Mutex
}

// NewFutureTask creates a new FutureTask.
func NewFutureTask() *FutureTask {
	return &FutureTask{
		state:  scheduled,
		result: nil,
		done:   make(chan bool, 1), // Using a buffered channel allows the future resolver to async notify any listeners.
	}
}

// GetResult will provide the result of the computation. If the future is not computed yet, it will block.
func (f *FutureTask) GetResult() (interface{}, error) {
	// The current state can be changed by another thread using the Resolve function, so we lock the read if the state is
	// being changed.
	currentState := func() futureState {
		f.mux.Lock()
		defer f.mux.Unlock()

		return f.state
	}()

	switch currentState {
	case thrown:
		return nil, f.result.(error)
	case scheduled:
		<-f.done // Blocking wait until the future is resolved
		if f.state == thrown {
			return nil, f.result.(error)
		}
		fallthrough
	case computed:
		return f.result, nil
	}

	return nil, errors.New("unsupported future state; this should not happen")
}

// Resolve completes the computation. Any thread or goroutine waiting for GetResult() will be unblocked.
func (f *FutureTask) Resolve(result interface{}) error {
	f.mux.Lock()
	defer f.mux.Unlock()

	switch f.state {
	case thrown:
		return errors.New("future has already been rejected")
	case computed:
		return errors.New("future has already been resolved")
	case scheduled:
		f.result = result
		f.state = computed
		f.done <- true
		return nil
	}

	return errors.New("unsupported future state; this should not happen")
}

func (f *FutureTask) Reject(err error) error {
	f.mux.Lock()
	defer f.mux.Unlock()

	switch f.state {
	case thrown:
		return errors.New("future has already been rejected")
	case computed:
		return errors.New("future has already been resolved")
	case scheduled:
		f.result = err
		f.state = thrown
		f.done <- true
		return nil
	}

	return errors.New("unsupported future state; this should not happen")
}
