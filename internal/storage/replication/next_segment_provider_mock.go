// Code generated by mockery v2.53.3. DO NOT EDIT.

package replication

import mock "github.com/stretchr/testify/mock"

// MockNextSegmentProvider is an autogenerated mock type for the NextSegmentProvider type
type MockNextSegmentProvider struct {
	mock.Mock
}

type MockNextSegmentProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNextSegmentProvider) EXPECT() *MockNextSegmentProvider_Expecter {
	return &MockNextSegmentProvider_Expecter{mock: &_m.Mock}
}

// NextSegment provides a mock function with given fields: filename
func (_m *MockNextSegmentProvider) NextSegment(filename string) (string, error) {
	ret := _m.Called(filename)

	if len(ret) == 0 {
		panic("no return value specified for NextSegment")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(filename)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(filename)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNextSegmentProvider_NextSegment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NextSegment'
type MockNextSegmentProvider_NextSegment_Call struct {
	*mock.Call
}

// NextSegment is a helper method to define mock.On call
//   - filename string
func (_e *MockNextSegmentProvider_Expecter) NextSegment(filename interface{}) *MockNextSegmentProvider_NextSegment_Call {
	return &MockNextSegmentProvider_NextSegment_Call{Call: _e.mock.On("NextSegment", filename)}
}

func (_c *MockNextSegmentProvider_NextSegment_Call) Run(run func(filename string)) *MockNextSegmentProvider_NextSegment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockNextSegmentProvider_NextSegment_Call) Return(_a0 string, _a1 error) *MockNextSegmentProvider_NextSegment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNextSegmentProvider_NextSegment_Call) RunAndReturn(run func(string) (string, error)) *MockNextSegmentProvider_NextSegment_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockNextSegmentProvider creates a new instance of MockNextSegmentProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNextSegmentProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNextSegmentProvider {
	mock := &MockNextSegmentProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
