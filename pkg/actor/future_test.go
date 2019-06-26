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
		if err != nil {
			assert.Fail(t, "Resolve returned an error", err)
		}
	}()

	// Wait for the response
	result, err := task.GetResult()
	if err != nil {
		assert.Fail(t, "GetResult returned an error", err)
	}

	assert.Equal(t, expectedResult, result.(int))
}

func TestFutureTask_GetResult_QuickResolve(t *testing.T) {
	expectedResult := 42

	// Create a future and resolve immediately
	task := NewFutureTask()
	err := task.Resolve(expectedResult)
	if err != nil {
		assert.Fail(t, "Resolve returned an error", err)
	}

	// Get the response
	result, err := task.GetResult()
	if err != nil {
		assert.Fail(t, "GetResult returned an error", err)
	}

	assert.Equal(t, expectedResult, result.(int))
}

func TestFutureTask_GetResult_Reject(t *testing.T) {
	expectedMsg := "rejected"

	// Create a future, sleep for a second, and reject
	task := NewFutureTask()

	go func() {
		time.Sleep(time.Second)
		err := task.Reject(errors.New(expectedMsg))
		if err != nil {
			assert.Fail(t, "Reject returned an error", err)
		}
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
	if err != nil {
		assert.Fail(t, "Reject returned an error", err)
	}

	// Get the response
	result, err := task.GetResult()
	assert.Nil(t, result)
	assert.EqualError(t, err, expectedMsg)
}
