// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// URLSaver is an autogenerated mock type for the URLSaver type
type URLSaver struct {
	mock.Mock
}

// SaveURL provides a mock function with given fields: URL, alias
func (_m *URLSaver) SaveURL(URL string, alias string) (int64, error) {
	ret := _m.Called(URL, alias)

	if len(ret) == 0 {
		panic("no return value specified for SaveURL")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (int64, error)); ok {
		return rf(URL, alias)
	}
	if rf, ok := ret.Get(0).(func(string, string) int64); ok {
		r0 = rf(URL, alias)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(URL, alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewURLSaver creates a new instance of URLSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLSaver {
	mock := &URLSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
