// Code generated by mockery v2.53.3. DO NOT EDIT.

package storage

import mock "github.com/stretchr/testify/mock"

// MockTransaction is an autogenerated mock type for the Transaction type
type MockTransaction struct {
	mock.Mock
}

type MockTransaction_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTransaction) EXPECT() *MockTransaction_Expecter {
	return &MockTransaction_Expecter{mock: &_m.Mock}
}

// Commit provides a mock function with no fields
func (_m *MockTransaction) Commit() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Commit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTransaction_Commit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Commit'
type MockTransaction_Commit_Call struct {
	*mock.Call
}

// Commit is a helper method to define mock.On call
func (_e *MockTransaction_Expecter) Commit() *MockTransaction_Commit_Call {
	return &MockTransaction_Commit_Call{Call: _e.mock.On("Commit")}
}

func (_c *MockTransaction_Commit_Call) Run(run func()) *MockTransaction_Commit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTransaction_Commit_Call) Return(_a0 error) *MockTransaction_Commit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransaction_Commit_Call) RunAndReturn(run func() error) *MockTransaction_Commit_Call {
	_c.Call.Return(run)
	return _c
}

// Del provides a mock function with given fields: key
func (_m *MockTransaction) Del(key string) error {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Del")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTransaction_Del_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Del'
type MockTransaction_Del_Call struct {
	*mock.Call
}

// Del is a helper method to define mock.On call
//   - key string
func (_e *MockTransaction_Expecter) Del(key interface{}) *MockTransaction_Del_Call {
	return &MockTransaction_Del_Call{Call: _e.mock.On("Del", key)}
}

func (_c *MockTransaction_Del_Call) Run(run func(key string)) *MockTransaction_Del_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockTransaction_Del_Call) Return(_a0 error) *MockTransaction_Del_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransaction_Del_Call) RunAndReturn(run func(string) error) *MockTransaction_Del_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: key
func (_m *MockTransaction) Get(key string) (string, bool) {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 string
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (string, bool)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockTransaction_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockTransaction_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - key string
func (_e *MockTransaction_Expecter) Get(key interface{}) *MockTransaction_Get_Call {
	return &MockTransaction_Get_Call{Call: _e.mock.On("Get", key)}
}

func (_c *MockTransaction_Get_Call) Run(run func(key string)) *MockTransaction_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockTransaction_Get_Call) Return(_a0 string, _a1 bool) *MockTransaction_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTransaction_Get_Call) RunAndReturn(run func(string) (string, bool)) *MockTransaction_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Rollback provides a mock function with no fields
func (_m *MockTransaction) Rollback() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Rollback")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTransaction_Rollback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rollback'
type MockTransaction_Rollback_Call struct {
	*mock.Call
}

// Rollback is a helper method to define mock.On call
func (_e *MockTransaction_Expecter) Rollback() *MockTransaction_Rollback_Call {
	return &MockTransaction_Rollback_Call{Call: _e.mock.On("Rollback")}
}

func (_c *MockTransaction_Rollback_Call) Run(run func()) *MockTransaction_Rollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTransaction_Rollback_Call) Return(_a0 error) *MockTransaction_Rollback_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransaction_Rollback_Call) RunAndReturn(run func() error) *MockTransaction_Rollback_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: key, value
func (_m *MockTransaction) Set(key string, value string) error {
	ret := _m.Called(key, value)

	if len(ret) == 0 {
		panic("no return value specified for Set")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTransaction_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type MockTransaction_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - key string
//   - value string
func (_e *MockTransaction_Expecter) Set(key interface{}, value interface{}) *MockTransaction_Set_Call {
	return &MockTransaction_Set_Call{Call: _e.mock.On("Set", key, value)}
}

func (_c *MockTransaction_Set_Call) Run(run func(key string, value string)) *MockTransaction_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockTransaction_Set_Call) Return(_a0 error) *MockTransaction_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransaction_Set_Call) RunAndReturn(run func(string, string) error) *MockTransaction_Set_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTransaction creates a new instance of MockTransaction. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTransaction(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTransaction {
	mock := &MockTransaction{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
