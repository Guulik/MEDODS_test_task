// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// TokenModifier is an autogenerated mock type for the TokenModifier type
type TokenModifier struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, userID
func (_m *TokenModifier) Delete(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: ctx, userID, tokenHash, ipAddress, expiresAt
func (_m *TokenModifier) Insert(ctx context.Context, userID string, tokenHash string, ipAddress string, expiresAt time.Time) error {
	ret := _m.Called(ctx, userID, tokenHash, ipAddress, expiresAt)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, time.Time) error); ok {
		r0 = rf(ctx, userID, tokenHash, ipAddress, expiresAt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTokenModifier creates a new instance of TokenModifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenModifier(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenModifier {
	mock := &TokenModifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
