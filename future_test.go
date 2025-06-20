package async_test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/reugn/async"

	"github.com/reugn/async/internal/assert"
	"github.com/reugn/async/internal/util"
)

func TestFuture(t *testing.T) {
	p := async.NewPromise[bool]()
	go func() {
		time.Sleep(100 * time.Millisecond)
		p.Success(true)
	}()

	res, err := p.Future().Join()

	assert.Equal(t, true, res)
	assert.IsNil(t, err)
}

func TestFuture_Utils(t *testing.T) {
	p1 := async.NewPromise[*int]()
	p2 := async.NewPromise[*int]()
	p3 := async.NewPromise[*int]()

	res1 := util.Ptr(1)
	res2 := util.Ptr(2)
	err3 := errors.New("error")

	go func() {
		time.Sleep(100 * time.Millisecond)
		p1.Success(res1)
		time.Sleep(200 * time.Millisecond)
		p2.Success(res2)
		time.Sleep(300 * time.Millisecond)
		p3.Failure(err3)
	}()

	arr := []async.Future[*int]{p1.Future(), p2.Future(), p3.Future()}
	res := []any{res1, res2, err3}
	futRes, _ := async.FutureSeq(arr).Join()

	assert.Equal(t, res, futRes)
}

func TestFuture_FirstCompleted(t *testing.T) {
	p := async.NewPromise[*bool]()
	go func() {
		time.Sleep(100 * time.Millisecond)
		p.Success(util.Ptr(true))
	}()

	timeout := async.FutureTimer[*bool](10 * time.Millisecond)
	futRes, futErr := async.FutureFirstCompletedOf(p.Future(), timeout).Join()

	assert.IsNil(t, futRes)
	assert.NotEqual(t, futErr, nil)
}

func TestFuture_Transform(t *testing.T) {
	p1 := async.NewPromise[*int]()
	go func() {
		time.Sleep(100 * time.Millisecond)
		p1.Success(util.Ptr(1))
	}()

	future := p1.Future().Map(func(v *int) (*int, error) {
		inc := *v + 1
		return &inc, nil
	}).FlatMap(func(v *int) (async.Future[*int], error) {
		inc := *v + 1
		p2 := async.NewPromise[*int]()
		p2.Success(&inc)
		return p2.Future(), nil
	}).Recover(func() (*int, error) {
		return util.Ptr(5), nil
	})

	res, _ := future.Get(t.Context())
	assert.Equal(t, 3, *res)

	res, _ = future.Join()
	assert.Equal(t, 3, *res)
}

func TestFuture_Recover(t *testing.T) {
	p1 := async.NewPromise[int]()
	p2 := async.NewPromise[int]()
	go func() {
		time.Sleep(10 * time.Millisecond)
		p1.Success(1)
		time.Sleep(10 * time.Millisecond)
		p2.Failure(errors.New("recover Future failure"))
	}()

	future := p1.Future().Map(func(_ int) (int, error) {
		return 0, errors.New("map error")
	}).FlatMap(func(_ int) (async.Future[int], error) {
		p2 = async.NewPromise[int]()
		p2.Failure(errors.New("flatMap Future failure"))
		return p2.Future(), nil
	}).FlatMap(func(_ int) (async.Future[int], error) {
		return nil, errors.New("flatMap error")
	}).Recover(func() (int, error) {
		return 0, errors.New("recover error")
	}).RecoverWith(p2.Future()).Recover(func() (int, error) {
		return 2, nil
	})

	res, err := future.Join()
	assert.Equal(t, 2, res)
	assert.IsNil(t, err)
}

func TestFuture_Failure(t *testing.T) {
	p1 := async.NewPromise[int]()
	p2 := async.NewPromise[int]()
	p3 := async.NewPromise[int]()
	err := errors.New("error")
	go func() {
		time.Sleep(5 * time.Millisecond)
		p1.Failure(err)
		time.Sleep(5 * time.Millisecond)
		p2.Failure(err)
		time.Sleep(5 * time.Millisecond)
		p3.Success(2)
	}()

	res, _ := p1.Future().
		Map(func(_ int) (int, error) { return 0, err }).
		FlatMap(func(_ int) (async.Future[int], error) { return p1.Future(), err }).
		RecoverWith(p2.Future()).
		RecoverWith(p3.Future()).
		FlatMap(func(_ int) (async.Future[int], error) { return p2.Future(), err }).
		RecoverWith(p3.Future()).
		RecoverWith(p3.Future()).
		Join()

	assert.Equal(t, 2, res)
}

func TestFuture_Timeout(t *testing.T) {
	p := async.NewPromise[bool]()
	go func() {
		time.Sleep(100 * time.Millisecond)
		p.Success(true)
	}()

	future := p.Future()

	ctx, cancel := context.WithTimeout(t.Context(),
		10*time.Millisecond)
	defer cancel()

	_, err := future.Get(ctx)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	_, err = future.Join()
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestFuture_GoroutineLeak(t *testing.T) {
	var wg sync.WaitGroup
	fmt.Println(runtime.NumGoroutine())

	numFuture := 100
	for range numFuture {
		promise := async.NewPromise[*string]()
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			promise.Success(util.Ptr("OK"))
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			fut := promise.Future()
			_, _ = fut.Get(t.Context())
			time.Sleep(100 * time.Millisecond)
			_, _ = fut.Join()
		}()
	}

	wg.Wait()
	time.Sleep(10 * time.Millisecond)
	numGoroutine := runtime.NumGoroutine()
	fmt.Println(numGoroutine)

	if numGoroutine > numFuture {
		t.Fatalf("numGoroutine is %d", numGoroutine)
	}
}
