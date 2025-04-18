// Code generated by mockery v2.53.3. DO NOT EDIT.

package storage

import (
	parser "github.com/DaniilZ77/InMemDB/internal/compute/parser"
	mock "github.com/stretchr/testify/mock"
)

// MockCompute is an autogenerated mock type for the Compute type
type MockCompute struct {
	mock.Mock
}

type MockCompute_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCompute) EXPECT() *MockCompute_Expecter {
	return &MockCompute_Expecter{mock: &_m.Mock}
}

// Parse provides a mock function with given fields: source
func (_m *MockCompute) Parse(source string) (*parser.Command, error) {
	ret := _m.Called(source)

	if len(ret) == 0 {
		panic("no return value specified for Parse")
	}

	var r0 *parser.Command
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*parser.Command, error)); ok {
		return rf(source)
	}
	if rf, ok := ret.Get(0).(func(string) *parser.Command); ok {
		r0 = rf(source)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*parser.Command)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(source)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCompute_Parse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Parse'
type MockCompute_Parse_Call struct {
	*mock.Call
}

// Parse is a helper method to define mock.On call
//   - source string
func (_e *MockCompute_Expecter) Parse(source interface{}) *MockCompute_Parse_Call {
	return &MockCompute_Parse_Call{Call: _e.mock.On("Parse", source)}
}

func (_c *MockCompute_Parse_Call) Run(run func(source string)) *MockCompute_Parse_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCompute_Parse_Call) Return(_a0 *parser.Command, _a1 error) *MockCompute_Parse_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCompute_Parse_Call) RunAndReturn(run func(string) (*parser.Command, error)) *MockCompute_Parse_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCompute creates a new instance of MockCompute. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCompute(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCompute {
	mock := &MockCompute{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
