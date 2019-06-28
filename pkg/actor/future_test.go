package actor

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFutureTask_GetResult_Resolve(t *testing.T) {
	expectedResult := 42

	// Create a future, sleep for a second, and resolve
	task := NewFutureTask()

	go func() {
		time.Sleep(time.Second)
		err := task.Resolve(expectedResult)
		assert.Nil(t, err)
	}()

	// Wait for the response
	result, err := task.GetResult()
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result.(int))
}

func TestFutureTask_GetResult_QuickResolve(t *testing.T) {
	expectedResult := 42

	// Create a future and resolve immediately
	task := NewFutureTask()
	err := task.Resolve(expectedResult)
	assert.Nil(t, err)

	// Get the response
	result, err := task.GetResult()
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result.(int))
}

func TestFutureTask_GetResult_Reject(t *testing.T) {
	expectedMsg := "rejected"

	// Create a future, sleep for a second, and reject
	task := NewFutureTask()

	go func() {
		time.Sleep(time.Second)
		err := task.Reject(errors.New(expectedMsg))
		assert.Nil(t, err)
	}()

	// Wait for the response
	result, err := task.GetResult()
	assert.Nil(t, result)
	assert.EqualError(t, err, expectedMsg)
}

func TestFutureTask_GetResult_QuickReject(t *testing.T) {
	expectedMsg := "rejected"

	// Create a future and reject immediately
	task := NewFutureTask()
	err := task.Reject(errors.New(expectedMsg))
	assert.Nil(t, err)

	// Get the response
	result, err := task.GetResult()
	assert.Nil(t, result)
	assert.EqualError(t, err, expectedMsg)
}

func TestFutureTask_Resolve_AlreadyRejected(t *testing.T) {
	expectedRejectErr := "rejected"
	expectedMsg := "future has already been rejected"

	// Create a future, reject, and then resolve.
	task := NewFutureTask()

	err := task.Reject(errors.New(expectedRejectErr))
	assert.Nil(t, err)

	err = task.Resolve(42)
	assert.EqualError(t, err, expectedMsg)

	// Assert that the result is rejected, not resolved
	result, err := task.GetResult()
	assert.Nil(t, result)
	assert.EqualError(t, err, expectedRejectErr)
}

func TestFutureTask_Resolve_AlreadyResolved(t *testing.T) {
	expectedResult := 42
	expectedMsg := "future has already been resolved"

	// Create a future, resolve, and then resolve.
	task := NewFutureTask()

	err := task.Resolve(expectedResult)
	assert.Nil(t, err)

	err = task.Resolve(expectedResult + 1)
	assert.EqualError(t, err, expectedMsg)

	// Assert that the result is resolved with the first value
	result, err := task.GetResult()
	assert.Nil(t, err)
	assert.Equal(t, result, result.(int))
}

func TestFutureTask_Reject_AlreadyRejected(t *testing.T) {
	expectedRejectErr := "rejected"
	expectedMsg := "future has already been rejected"

	// Create a future, reject, and then reject.
	task := NewFutureTask()

	err := task.Reject(errors.New(expectedRejectErr))
	assert.Nil(t, err)

	err = task.Reject(errors.New("rejected again"))
	assert.EqualError(t, err, expectedMsg)

	// Assert that the result is rejected with the first rejection value
	result, err := task.GetResult()
	assert.Nil(t, result)
	assert.EqualError(t, err, expectedRejectErr)
}

func TestFutureTask_Reject_AlreadyResolved(t *testing.T) {
	expectedResult := 42
	expectedMsg := "future has already been resolved"

	// Create a future, resolve, and then reject.
	task := NewFutureTask()

	err := task.Resolve(expectedResult)
	assert.Nil(t, err)

	err = task.Reject(errors.New("rejected"))
	assert.EqualError(t, err, expectedMsg)

	// Assert that the result is resolved, not rejected
	result, err := task.GetResult()
	assert.Nil(t, err)
	assert.Equal(t, result, result.(int))
}
