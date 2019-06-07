package actor

import (
	"errors"
)

// Future promises the result of a long running job.
type Future interface {
	GetResult() (interface{}, error)
}

// FutureState is the current state of the future.
type FutureState uint8

const (
	INVALID FutureState = iota
	SCHEDULED
	COMPUTED
	CANCELLED
	THROWN
)

// FutureTask is the concrete implementation of a Future.
type FutureTask struct {
	state  FutureState
	result interface{}
	done   chan bool
}

// NewFutureTask creates a new FutureTask.
func NewFutureTask() *FutureTask {
	return &FutureTask{
		state:  SCHEDULED,
		result: nil,
		done:   make(chan bool, 1), // Using a buffered channel allows the future resolver to async notify any listeners.
	}
}

// GetResult will provide the result of the computation. If the future is not computed yet, it will block.
func (f *FutureTask) GetResult() (interface{}, error) {
	switch f.state {
	case INVALID:
		// TODO: Return error for getting result twice
	case SCHEDULED:
		<-f.done
		return f.result, nil
	case COMPUTED:
		return f.result, nil
	case CANCELLED:
		// TODO: Return error
	case THROWN:
		// TODO: Return error
	}

	return nil, errors.New("unsupported future state; this should not happen")
}

// SetResult completes the computation. Any thread or goroutine waiting for GetResult() will be unblocked.
func (f *FutureTask) SetResult(result interface{}) error {
	switch f.state {
	case INVALID:
		// TODO: Return error for setting result after read
	case SCHEDULED:
		f.result = result
		f.state = COMPUTED
		f.done <- true
		return nil
	case COMPUTED:
		// TODO: Return error for setting result twice
	case CANCELLED:
		// TODO: Return error for setting result after cancel
	case THROWN:
		// TODO: Return error for setting result after error
	}

	return errors.New("unsupported future state; this should not happen")
}
