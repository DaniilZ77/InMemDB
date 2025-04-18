// Code generated by mockery v2.53.3. DO NOT EDIT.

package storage

import mock "github.com/stretchr/testify/mock"

// MockEngine is an autogenerated mock type for the Engine type
type MockEngine struct {
	mock.Mock
}

type MockEngine_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEngine) EXPECT() *MockEngine_Expecter {
	return &MockEngine_Expecter{mock: &_m.Mock}
}

// Del provides a mock function with given fields: key
func (_m *MockEngine) Del(key string) {
	_m.Called(key)
}

// MockEngine_Del_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Del'
type MockEngine_Del_Call struct {
	*mock.Call
}

// Del is a helper method to define mock.On call
//   - key string
func (_e *MockEngine_Expecter) Del(key interface{}) *MockEngine_Del_Call {
	return &MockEngine_Del_Call{Call: _e.mock.On("Del", key)}
}

func (_c *MockEngine_Del_Call) Run(run func(key string)) *MockEngine_Del_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockEngine_Del_Call) Return() *MockEngine_Del_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockEngine_Del_Call) RunAndReturn(run func(string)) *MockEngine_Del_Call {
	_c.Run(run)
	return _c
}

// Get provides a mock function with given fields: key
func (_m *MockEngine) Get(key string) (string, bool) {
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

// MockEngine_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockEngine_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - key string
func (_e *MockEngine_Expecter) Get(key interface{}) *MockEngine_Get_Call {
	return &MockEngine_Get_Call{Call: _e.mock.On("Get", key)}
}

func (_c *MockEngine_Get_Call) Run(run func(key string)) *MockEngine_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockEngine_Get_Call) Return(_a0 string, _a1 bool) *MockEngine_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEngine_Get_Call) RunAndReturn(run func(string) (string, bool)) *MockEngine_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: key, value
func (_m *MockEngine) Set(key string, value string) {
	_m.Called(key, value)
}

// MockEngine_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type MockEngine_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - key string
//   - value string
func (_e *MockEngine_Expecter) Set(key interface{}, value interface{}) *MockEngine_Set_Call {
	return &MockEngine_Set_Call{Call: _e.mock.On("Set", key, value)}
}

func (_c *MockEngine_Set_Call) Run(run func(key string, value string)) *MockEngine_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockEngine_Set_Call) Return() *MockEngine_Set_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockEngine_Set_Call) RunAndReturn(run func(string, string)) *MockEngine_Set_Call {
	_c.Run(run)
	return _c
}

// NewMockEngine creates a new instance of MockEngine. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEngine(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEngine {
	mock := &MockEngine{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
