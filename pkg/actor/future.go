package actor

import (
	"errors"
)

// Future promises the result of a long running job.
type Future interface {
	GetResult() (interface{}, error)
}

type futureState uint8

const (
	invalid futureState = iota
	scheduled
	computed
	cancelled
	thrown
)

// FutureTask is the concrete implementation of a Future.
type FutureTask struct {
	state  futureState
	result interface{}
	done   chan bool
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
	switch f.state {
	case invalid:
		// TODO: Return error for getting result twice
	case scheduled:
		<-f.done
		return f.result, nil
	case computed:
		return f.result, nil
	case cancelled:
		// TODO: Return error
	case thrown:
		// TODO: Return error
	}

	return nil, errors.New("unsupported future state; this should not happen")
}

// SetResult completes the computation. Any thread or goroutine waiting for GetResult() will be unblocked.
func (f *FutureTask) SetResult(result interface{}) error {
	switch f.state {
	case invalid:
		// TODO: Return error for setting result after read
	case scheduled:
		f.result = result
		f.state = computed
		f.done <- true
		return nil
	case computed:
		// TODO: Return error for setting result twice
	case cancelled:
		// TODO: Return error for setting result after cancel
	case thrown:
		// TODO: Return error for setting result after error
	}

	return errors.New("unsupported future state; this should not happen")
}
