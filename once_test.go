package async_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/reugn/async"

	"github.com/reugn/async/internal/assert"
	"github.com/reugn/async/internal/util"
)

func TestOnce(t *testing.T) {
	var once async.Once[int]
	var count int

	for range 10 {
		count, _ = once.Do(func() (int, error) {
			count++
			return count, nil
		})
	}
	assert.Equal(t, 1, count)
}

func TestOnce_Ptr(t *testing.T) {
	var once async.Once[*int]
	count := new(int)

	for range 10 {
		count, _ = once.Do(func() (*int, error) {
			*count++
			return count, nil
		})
	}
	assert.Equal(t, util.Ptr(1), count)
}

func TestOnce_Concurrent(t *testing.T) {
	var once async.Once[*int32]
	var count atomic.Int32
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, _ := once.Do(func() (*int32, error) {
				newCount := count.Add(1)
				return &newCount, nil
			})
			count.Store(*result)
		}()
	}
	wg.Wait()
	assert.Equal(t, 1, int(count.Load()))
}

func TestOnce_Panic(t *testing.T) {
	var once async.Once[*int]
	count := new(int)
	var err error

	for range 10 {
		count, err = once.Do(func() (*int, error) {
			*count /= *count
			return count, nil
		})
	}
	assert.ErrorContains(t, err, "integer divide by zero")
}
