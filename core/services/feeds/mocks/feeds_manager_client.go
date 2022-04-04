// Code generated by mockery v2.10.1. DO NOT EDIT.

package mocks

import (
	context "context"

	proto "github.com/smartcontractkit/chainlink/core/services/feeds/proto"
	mock "github.com/stretchr/testify/mock"
)

// FeedsManagerClient is an autogenerated mock type for the FeedsManagerClient type
type FeedsManagerClient struct {
	mock.Mock
}

// ApprovedJob provides a mock function with given fields: ctx, in
func (_m *FeedsManagerClient) ApprovedJob(ctx context.Context, in *proto.ApprovedJobRequest) (*proto.ApprovedJobResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *proto.ApprovedJobResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ApprovedJobRequest) *proto.ApprovedJobResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ApprovedJobResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.ApprovedJobRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CancelledJob provides a mock function with given fields: ctx, in
func (_m *FeedsManagerClient) CancelledJob(ctx context.Context, in *proto.CancelledJobRequest) (*proto.CancelledJobResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *proto.CancelledJobResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.CancelledJobRequest) *proto.CancelledJobResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.CancelledJobResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.CancelledJobRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RejectedJob provides a mock function with given fields: ctx, in
func (_m *FeedsManagerClient) RejectedJob(ctx context.Context, in *proto.RejectedJobRequest) (*proto.RejectedJobResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *proto.RejectedJobResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.RejectedJobRequest) *proto.RejectedJobResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.RejectedJobResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.RejectedJobRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateNode provides a mock function with given fields: ctx, in
func (_m *FeedsManagerClient) UpdateNode(ctx context.Context, in *proto.UpdateNodeRequest) (*proto.UpdateNodeResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *proto.UpdateNodeResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.UpdateNodeRequest) *proto.UpdateNodeResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.UpdateNodeResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.UpdateNodeRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
