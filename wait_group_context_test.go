package async_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/reugn/async"

	"github.com/reugn/async/internal/assert"
)

func TestWaitGroupContext(t *testing.T) {
	var result atomic.Int32
	wgc := async.NewWaitGroupContext(t.Context())
	wgc.Add(2)

	go func() {
		defer wgc.Done()
		time.Sleep(10 * time.Millisecond)
		result.Add(1)
	}()
	go func() {
		defer wgc.Done()
		time.Sleep(20 * time.Millisecond)
		result.Add(2)
	}()
	go func() {
		wgc.Wait()
		result.Add(3)
	}()

	wgc.Wait()
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, int(result.Load()), 6)
}

func TestWaitGroupContextCanceled(t *testing.T) {
	var result atomic.Int32
	ctx, cancelFunc := context.WithCancel(t.Context())
	go func() {
		time.Sleep(100 * time.Millisecond)
		result.Add(10)
		cancelFunc()
	}()
	wgc := async.NewWaitGroupContext(ctx)
	wgc.Add(2)

	go func() {
		defer wgc.Done()
		time.Sleep(10 * time.Millisecond)
		result.Add(1)
	}()
	go func() {
		defer wgc.Done()
		time.Sleep(300 * time.Millisecond)
		result.Add(2)
	}()
	go func() {
		wgc.Wait()
		result.Add(100)
	}()

	wgc.Wait()
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, int(result.Load()), 111)
}

func TestWaitGroupContextPanicNegativeCounter(t *testing.T) {
	negativeCounter := func() {
		wgc := async.NewWaitGroupContext(t.Context())
		wgc.Add(-2)
	}
	assert.PanicMsgContains(t, negativeCounter, "negative")
}

func TestWaitGroupContextPanicReused(t *testing.T) {
	reusedBeforeWaitReturned := func() {
		var result atomic.Int32
		wgc := async.NewWaitGroupContext(t.Context())

		n := 10
		for range n {
			wgc.Add(1)
			go func() {
				defer wgc.Add(1)
				defer wgc.Done()
				result.Add(1)
			}()
			wgc.Wait()
		}
	}
	assert.PanicMsgContains(t, reusedBeforeWaitReturned, "reused")
}

func TestWaitGroupContextReused(t *testing.T) {
	var result atomic.Int32
	wgc := async.NewWaitGroupContext(t.Context())

	n := 1000
	for i := range n {
		assert.Equal(t, int(result.Load()), i*3)
		wgc.Add(2)
		go func() {
			defer wgc.Done()
			result.Add(1)
		}()
		go func() {
			defer wgc.Done()
			result.Add(1)
		}()
		go func() {
			wgc.Wait()
			result.Add(1)
		}()
		wgc.Wait()
		time.Sleep(time.Millisecond)
	}

	assert.Equal(t, int(result.Load()), n*3)
}
