// Code generated by mockery v2.53.3. DO NOT EDIT.

package wal

import mock "github.com/stretchr/testify/mock"

// MockDisk is an autogenerated mock type for the Disk type
type MockDisk struct {
	mock.Mock
}

type MockDisk_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDisk) EXPECT() *MockDisk_Expecter {
	return &MockDisk_Expecter{mock: &_m.Mock}
}

// Read provides a mock function with no fields
func (_m *MockDisk) Read() ([]byte, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDisk_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type MockDisk_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
func (_e *MockDisk_Expecter) Read() *MockDisk_Read_Call {
	return &MockDisk_Read_Call{Call: _e.mock.On("Read")}
}

func (_c *MockDisk_Read_Call) Run(run func()) *MockDisk_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDisk_Read_Call) Return(_a0 []byte, _a1 error) *MockDisk_Read_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDisk_Read_Call) RunAndReturn(run func() ([]byte, error)) *MockDisk_Read_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: _a0
func (_m *MockDisk) Write(_a0 []byte) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDisk_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type MockDisk_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - _a0 []byte
func (_e *MockDisk_Expecter) Write(_a0 interface{}) *MockDisk_Write_Call {
	return &MockDisk_Write_Call{Call: _e.mock.On("Write", _a0)}
}

func (_c *MockDisk_Write_Call) Run(run func(_a0 []byte)) *MockDisk_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockDisk_Write_Call) Return(_a0 error) *MockDisk_Write_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDisk_Write_Call) RunAndReturn(run func([]byte) error) *MockDisk_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDisk creates a new instance of MockDisk. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDisk(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDisk {
	mock := &MockDisk{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
