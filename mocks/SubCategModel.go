// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	domain "github.com/OYE0303/expense-tracker-go/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// SubCategModel is an autogenerated mock type for the SubCategModel type
type SubCategModel struct {
	mock.Mock
}

// Create provides a mock function with given fields: categ, userID
func (_m *SubCategModel) Create(categ *domain.SubCateg, userID int64) error {
	ret := _m.Called(categ, userID)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.SubCateg, int64) error); ok {
		r0 = rf(categ, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *SubCategModel) Delete(id int64) error {
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

// GetAll provides a mock function with given fields: userID
func (_m *SubCategModel) GetAll(userID int64) ([]*domain.SubCateg, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*domain.SubCateg
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) ([]*domain.SubCateg, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(int64) []*domain.SubCateg); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.SubCateg)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id, userID
func (_m *SubCategModel) GetByID(id int64, userID int64) (*domain.SubCateg, error) {
	ret := _m.Called(id, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *domain.SubCateg
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int64) (*domain.SubCateg, error)); ok {
		return rf(id, userID)
	}
	if rf, ok := ret.Get(0).(func(int64, int64) *domain.SubCateg); ok {
		r0 = rf(id, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.SubCateg)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(id, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByMainCategID provides a mock function with given fields: userID, mainCategID
func (_m *SubCategModel) GetByMainCategID(userID int64, mainCategID int64) ([]*domain.SubCateg, error) {
	ret := _m.Called(userID, mainCategID)

	if len(ret) == 0 {
		panic("no return value specified for GetByMainCategID")
	}

	var r0 []*domain.SubCateg
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int64) ([]*domain.SubCateg, error)); ok {
		return rf(userID, mainCategID)
	}
	if rf, ok := ret.Get(0).(func(int64, int64) []*domain.SubCateg); ok {
		r0 = rf(userID, mainCategID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.SubCateg)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(userID, mainCategID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: categ
func (_m *SubCategModel) Update(categ *domain.SubCateg) error {
	ret := _m.Called(categ)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.SubCateg) error); ok {
		r0 = rf(categ)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewSubCategModel creates a new instance of SubCategModel. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSubCategModel(t interface {
	mock.TestingT
	Cleanup(func())
}) *SubCategModel {
	mock := &SubCategModel{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
