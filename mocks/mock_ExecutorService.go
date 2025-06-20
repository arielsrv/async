// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package async

import (
	"context"

	"github.com/reugn/async"
	mock "github.com/stretchr/testify/mock"
)

// NewMockExecutorService creates a new instance of MockExecutorService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExecutorService[T any](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExecutorService[T] {
	mock := &MockExecutorService[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockExecutorService is an autogenerated mock type for the ExecutorService type
type MockExecutorService[T any] struct {
	mock.Mock
}

type MockExecutorService_Expecter[T any] struct {
	mock *mock.Mock
}

func (_m *MockExecutorService[T]) EXPECT() *MockExecutorService_Expecter[T] {
	return &MockExecutorService_Expecter[T]{mock: &_m.Mock}
}

// Shutdown provides a mock function for the type MockExecutorService
func (_mock *MockExecutorService[T]) Shutdown() error {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Shutdown")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func() error); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockExecutorService_Shutdown_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Shutdown'
type MockExecutorService_Shutdown_Call[T any] struct {
	*mock.Call
}

// Shutdown is a helper method to define mock.On call
func (_e *MockExecutorService_Expecter[T]) Shutdown() *MockExecutorService_Shutdown_Call[T] {
	return &MockExecutorService_Shutdown_Call[T]{Call: _e.mock.On("Shutdown")}
}

func (_c *MockExecutorService_Shutdown_Call[T]) Run(run func()) *MockExecutorService_Shutdown_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockExecutorService_Shutdown_Call[T]) Return(err error) *MockExecutorService_Shutdown_Call[T] {
	_c.Call.Return(err)
	return _c
}

func (_c *MockExecutorService_Shutdown_Call[T]) RunAndReturn(run func() error) *MockExecutorService_Shutdown_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Status provides a mock function for the type MockExecutorService
func (_mock *MockExecutorService[T]) Status() async.ExecutorStatus {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Status")
	}

	var r0 async.ExecutorStatus
	if returnFunc, ok := ret.Get(0).(func() async.ExecutorStatus); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(async.ExecutorStatus)
	}
	return r0
}

// MockExecutorService_Status_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Status'
type MockExecutorService_Status_Call[T any] struct {
	*mock.Call
}

// Status is a helper method to define mock.On call
func (_e *MockExecutorService_Expecter[T]) Status() *MockExecutorService_Status_Call[T] {
	return &MockExecutorService_Status_Call[T]{Call: _e.mock.On("Status")}
}

func (_c *MockExecutorService_Status_Call[T]) Run(run func()) *MockExecutorService_Status_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockExecutorService_Status_Call[T]) Return(executorStatus async.ExecutorStatus) *MockExecutorService_Status_Call[T] {
	_c.Call.Return(executorStatus)
	return _c
}

func (_c *MockExecutorService_Status_Call[T]) RunAndReturn(run func() async.ExecutorStatus) *MockExecutorService_Status_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Submit provides a mock function for the type MockExecutorService
func (_mock *MockExecutorService[T]) Submit(fn func(context.Context) (T, error)) (async.Future[T], error) {
	ret := _mock.Called(fn)

	if len(ret) == 0 {
		panic("no return value specified for Submit")
	}

	var r0 async.Future[T]
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(func(context.Context) (T, error)) (async.Future[T], error)); ok {
		return returnFunc(fn)
	}
	if returnFunc, ok := ret.Get(0).(func(func(context.Context) (T, error)) async.Future[T]); ok {
		r0 = returnFunc(fn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(async.Future[T])
		}
	}
	if returnFunc, ok := ret.Get(1).(func(func(context.Context) (T, error)) error); ok {
		r1 = returnFunc(fn)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockExecutorService_Submit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Submit'
type MockExecutorService_Submit_Call[T any] struct {
	*mock.Call
}

// Submit is a helper method to define mock.On call
//   - fn func(context.Context) (T, error)
func (_e *MockExecutorService_Expecter[T]) Submit(fn interface{}) *MockExecutorService_Submit_Call[T] {
	return &MockExecutorService_Submit_Call[T]{Call: _e.mock.On("Submit", fn)}
}

func (_c *MockExecutorService_Submit_Call[T]) Run(run func(fn func(context.Context) (T, error))) *MockExecutorService_Submit_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 func(context.Context) (T, error)
		if args[0] != nil {
			arg0 = args[0].(func(context.Context) (T, error))
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockExecutorService_Submit_Call[T]) Return(future async.Future[T], err error) *MockExecutorService_Submit_Call[T] {
	_c.Call.Return(future, err)
	return _c
}

func (_c *MockExecutorService_Submit_Call[T]) RunAndReturn(run func(fn func(context.Context) (T, error)) (async.Future[T], error)) *MockExecutorService_Submit_Call[T] {
	_c.Call.Return(run)
	return _c
}
