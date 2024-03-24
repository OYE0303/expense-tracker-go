// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/OYE0303/expense-tracker-go/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// TransactionUC is an autogenerated mock type for the TransactionUC type
type TransactionUC struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, trans
func (_m *TransactionUC) Create(ctx context.Context, trans domain.CreateTransactionInput) error {
	ret := _m.Called(ctx, trans)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.CreateTransactionInput) error); ok {
		r0 = rf(ctx, trans)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id, user
func (_m *TransactionUC) Delete(ctx context.Context, id int64, user domain.User) error {
	ret := _m.Called(ctx, id, user)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, domain.User) error); ok {
		r0 = rf(ctx, id, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAccInfo provides a mock function with given fields: ctx, query, user
func (_m *TransactionUC) GetAccInfo(ctx context.Context, query domain.GetAccInfoQuery, user domain.User) (domain.AccInfo, error) {
	ret := _m.Called(ctx, query, user)

	if len(ret) == 0 {
		panic("no return value specified for GetAccInfo")
	}

	var r0 domain.AccInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.GetAccInfoQuery, domain.User) (domain.AccInfo, error)); ok {
		return rf(ctx, query, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.GetAccInfoQuery, domain.User) domain.AccInfo); ok {
		r0 = rf(ctx, query, user)
	} else {
		r0 = ret.Get(0).(domain.AccInfo)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.GetAccInfoQuery, domain.User) error); ok {
		r1 = rf(ctx, query, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx, query, user
func (_m *TransactionUC) GetAll(ctx context.Context, query domain.GetQuery, user domain.User) ([]domain.Transaction, error) {
	ret := _m.Called(ctx, query, user)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.GetQuery, domain.User) ([]domain.Transaction, error)); ok {
		return rf(ctx, query, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.GetQuery, domain.User) []domain.Transaction); ok {
		r0 = rf(ctx, query, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.GetQuery, domain.User) error); ok {
		r1 = rf(ctx, query, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBarChartData provides a mock function with given fields: ctx, dataRange, transactionType, user
func (_m *TransactionUC) GetBarChartData(ctx context.Context, dataRange domain.ChartDateRange, transactionType domain.TransactionType, user domain.User) (domain.ChartData, error) {
	ret := _m.Called(ctx, dataRange, transactionType, user)

	if len(ret) == 0 {
		panic("no return value specified for GetBarChartData")
	}

	var r0 domain.ChartData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChartDateRange, domain.TransactionType, domain.User) (domain.ChartData, error)); ok {
		return rf(ctx, dataRange, transactionType, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChartDateRange, domain.TransactionType, domain.User) domain.ChartData); ok {
		r0 = rf(ctx, dataRange, transactionType, user)
	} else {
		r0 = ret.Get(0).(domain.ChartData)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ChartDateRange, domain.TransactionType, domain.User) error); ok {
		r1 = rf(ctx, dataRange, transactionType, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPieChartData provides a mock function with given fields: ctx, dataRange, transactionType, user
func (_m *TransactionUC) GetPieChartData(ctx context.Context, dataRange domain.ChartDateRange, transactionType domain.TransactionType, user domain.User) (domain.ChartData, error) {
	ret := _m.Called(ctx, dataRange, transactionType, user)

	if len(ret) == 0 {
		panic("no return value specified for GetPieChartData")
	}

	var r0 domain.ChartData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChartDateRange, domain.TransactionType, domain.User) (domain.ChartData, error)); ok {
		return rf(ctx, dataRange, transactionType, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChartDateRange, domain.TransactionType, domain.User) domain.ChartData); ok {
		r0 = rf(ctx, dataRange, transactionType, user)
	} else {
		r0 = ret.Get(0).(domain.ChartData)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ChartDateRange, domain.TransactionType, domain.User) error); ok {
		r1 = rf(ctx, dataRange, transactionType, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTransactionUC creates a new instance of TransactionUC. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionUC(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionUC {
	mock := &TransactionUC{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
