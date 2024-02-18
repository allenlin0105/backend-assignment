// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "dcard-backend/domain"

	mock "github.com/stretchr/testify/mock"
)

// AdUsecase is an autogenerated mock type for the AdUsecase type
type AdUsecase struct {
	mock.Mock
}

// Create provides a mock function with given fields: c, ad
func (_m *AdUsecase) Create(c context.Context, ad *domain.Ad) error {
	ret := _m.Called(c, ad)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Ad) error); ok {
		r0 = rf(c, ad)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByCondition provides a mock function with given fields: c, condition
func (_m *AdUsecase) GetByCondition(c context.Context, condition map[string][]string) ([]domain.Ad, error) {
	ret := _m.Called(c, condition)

	if len(ret) == 0 {
		panic("no return value specified for GetByCondition")
	}

	var r0 []domain.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string][]string) ([]domain.Ad, error)); ok {
		return rf(c, condition)
	}
	if rf, ok := ret.Get(0).(func(context.Context, map[string][]string) []domain.Ad); ok {
		r0 = rf(c, condition)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Ad)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, map[string][]string) error); ok {
		r1 = rf(c, condition)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAdUsecase creates a new instance of AdUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAdUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *AdUsecase {
	mock := &AdUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
