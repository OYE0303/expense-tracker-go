// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	domain "github.com/OYE0303/expense-tracker-go/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// MainCategUC is an autogenerated mock type for the MainCategUC type
type MainCategUC struct {
	mock.Mock
}

// Create provides a mock function with given fields: categ, userID
func (_m *MainCategUC) Create(categ domain.MainCateg, userID int64) error {
	ret := _m.Called(categ, userID)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.MainCateg, int64) error); ok {
		r0 = rf(categ, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *MainCategUC) Delete(id int64) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: userID, transType
func (_m *MainCategUC) GetAll(userID int64, transType domain.TransactionType) ([]domain.MainCateg, error) {
	ret := _m.Called(userID, transType)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.MainCateg
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, domain.TransactionType) ([]domain.MainCateg, error)); ok {
		return rf(userID, transType)
	}
	if rf, ok := ret.Get(0).(func(int64, domain.TransactionType) []domain.MainCateg); ok {
		r0 = rf(userID, transType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.MainCateg)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, domain.TransactionType) error); ok {
		r1 = rf(userID, transType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: categ, userID
func (_m *MainCategUC) Update(categ *domain.MainCateg, userID int64) error {
	ret := _m.Called(categ, userID)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.MainCateg, int64) error); ok {
		r0 = rf(categ, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMainCategUC creates a new instance of MainCategUC. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMainCategUC(t interface {
	mock.TestingT
	Cleanup(func())
}) *MainCategUC {
	mock := &MainCategUC{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
