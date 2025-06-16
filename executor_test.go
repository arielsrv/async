package async_test

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/reugn/async"

	"github.com/reugn/async/internal/assert"
)

func TestExecutor(t *testing.T) {
	ctx := t.Context()
	executor := async.NewExecutor[int](ctx, async.NewExecutorConfig(2, 2))

	job := func(_ context.Context) (int, error) {
		time.Sleep(time.Millisecond)
		return 1, nil
	}
	jobLong := func(_ context.Context) (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 1, nil
	}

	future1 := submitJob[int](t, executor, job)
	future2 := submitJob[int](t, executor, job)

	// wait for the first two jobs to complete
	time.Sleep(3 * time.Millisecond)

	// submit four more jobs
	future3 := submitJob[int](t, executor, jobLong)
	future4 := submitJob[int](t, executor, jobLong)
	future5 := submitJob[int](t, executor, jobLong)
	future6 := submitJob[int](t, executor, jobLong)

	// the queue has reached its maximum capacity
	future7, err := executor.Submit(job)
	assert.ErrorIs(t, err, async.ErrExecutorQueueFull)
	assert.IsNil(t, future7)

	assert.Equal(t, executor.Status(), async.ExecutorStatusRunning)

	routines := runtime.NumGoroutine()

	// shut down the executor
	_ = executor.Shutdown()
	time.Sleep(time.Millisecond)

	// verify that submit fails after the executor was shut down
	_, err = executor.Submit(job)
	assert.ErrorIs(t, err, async.ErrExecutorShutDown)

	// validate the executor status
	assert.Equal(t, executor.Status(), async.ExecutorStatusTerminating)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, executor.Status(), async.ExecutorStatusShutDown)

	assert.Equal(t, routines, runtime.NumGoroutine()+4)

	assertFutureResult(t, 1, future1, future2, future3, future4)
	assertFutureError(t, async.ErrExecutorShutDown, future5, future6)
}

func TestExecutor_context(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	executor := async.NewExecutor[int](ctx, async.NewExecutorConfig(2, 2))

	job := func(_ context.Context) (int, error) {
		return 0, errors.New("error")
	}

	future, err := executor.Submit(job)
	assert.IsNil(t, err)

	result, err := future.Join()
	assert.Equal(t, result, 0)
	assert.ErrorContains(t, err, "error")

	cancel()
	time.Sleep(5 * time.Millisecond)

	_, err = executor.Submit(job)
	assert.ErrorIs(t, err, async.ErrExecutorShutDown)

	assert.Equal(t, executor.Status(), async.ExecutorStatusShutDown)
}

func TestExecutor_jobPanic(t *testing.T) {
	ctx := t.Context()
	executor := async.NewExecutor[int](ctx, async.NewExecutorConfig(2, 2))

	job := func(_ context.Context) (int, error) {
		var i int
		return 1 / i, nil
	}

	future, err := executor.Submit(job)
	assert.IsNil(t, err)

	result, err := future.Join()
	assert.Equal(t, result, 0)
	assert.ErrorContains(t, err, "integer divide by zero")

	_ = executor.Shutdown()
}

func submitJob[T any](t *testing.T, executor async.ExecutorService[T],
	f func(context.Context) (T, error),
) async.Future[T] {
	future, err := executor.Submit(f)
	assert.IsNil(t, err)

	runtime.Gosched()
	return future
}

func assertFutureResult[T any](t *testing.T, expected T, futures ...async.Future[T]) {
	for _, future := range futures {
		result, err := future.Join()
		assert.IsNil(t, err)
		assert.Equal(t, expected, result)
	}
}

func assertFutureError[T any](t *testing.T, expected error, futures ...async.Future[T]) {
	for _, future := range futures {
		result, err := future.Join()
		var zero T
		assert.Equal(t, zero, result)
		assert.ErrorIs(t, err, expected)
	}
}
