// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// ReaperConfig is an autogenerated mock type for the ReaperConfig type
type ReaperConfig struct {
	mock.Mock
}

// EthTxReaperInterval provides a mock function with given fields:
func (_m *ReaperConfig) EthTxReaperInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// EthTxReaperThreshold provides a mock function with given fields:
func (_m *ReaperConfig) EthTxReaperThreshold() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// EvmFinalityDepth provides a mock function with given fields:
func (_m *ReaperConfig) EvmFinalityDepth() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}
