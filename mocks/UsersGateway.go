// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "rzd/app/entity"
import mock "github.com/stretchr/testify/mock"

// UsersGateway is an autogenerated mock type for the UsersGateway type
type UsersGateway struct {
	mock.Mock
}

// Create provides a mock function with given fields: user
func (_m *UsersGateway) Create(user entity.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(entity.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: user
func (_m *UsersGateway) Delete(user entity.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(entity.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadMany provides a mock function with given fields:
func (_m *UsersGateway) ReadMany() ([]entity.User, error) {
	ret := _m.Called()

	var r0 []entity.User
	if rf, ok := ret.Get(0).(func() []entity.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadOne provides a mock function with given fields: filter
func (_m *UsersGateway) ReadOne(filter entity.User) (entity.User, error) {
	ret := _m.Called(filter)

	var r0 entity.User
	if rf, ok := ret.Get(0).(func(entity.User) entity.User); ok {
		r0 = rf(filter)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(entity.User) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadSection provides a mock function with given fields: start, end
func (_m *UsersGateway) ReadSection(start int, end int) ([]entity.User, error) {
	ret := _m.Called(start, end)

	var r0 []entity.User
	if rf, ok := ret.Get(0).(func(int, int) []entity.User); ok {
		r0 = rf(start, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(start, end)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: user
func (_m *UsersGateway) Update(user entity.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(entity.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
