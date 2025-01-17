// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ReaperConfig is an autogenerated mock type for the ReaperConfig type
type ReaperConfig struct {
	mock.Mock
}

// FinalityDepth provides a mock function with given fields:
func (_m *ReaperConfig) FinalityDepth() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

type mockConstructorTestingTNewReaperConfig interface {
	mock.TestingT
	Cleanup(func())
}

// NewReaperConfig creates a new instance of ReaperConfig. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReaperConfig(t mockConstructorTestingTNewReaperConfig) *ReaperConfig {
	mock := &ReaperConfig{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
