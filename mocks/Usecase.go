// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "rzd/app/entity"
import mock "github.com/stretchr/testify/mock"

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// AddUser provides a mock function with given fields: user
func (_m *Usecase) AddUser(user *entity.User) (bool, error) {
	ret := _m.Called(user)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*entity.User) bool); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*entity.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckUsers provides a mock function with given fields: start, end
func (_m *Usecase) CheckUsers(start int64, end int64) ([]*entity.User, error) {
	ret := _m.Called(start, end)

	var r0 []*entity.User
	if rf, ok := ret.Get(0).(func(int64, int64) []*entity.User); ok {
		r0 = rf(start, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(start, end)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: user
func (_m *Usecase) DeleteUser(user *entity.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetInfoAboutTrains provides a mock function with given fields: args
func (_m *Usecase) GetInfoAboutTrains(args *entity.RouteArgs) ([]*entity.Train, error) {
	ret := _m.Called(args)

	var r0 []*entity.Train
	if rf, ok := ret.Get(0).(func(*entity.RouteArgs) []*entity.Train); ok {
		r0 = rf(args)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.Train)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*entity.RouteArgs) error); ok {
		r1 = rf(args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStationCodes provides a mock function with given fields: target, source
func (_m *Usecase) GetStationCodes(target string, source string) (int, int, error) {
	ret := _m.Called(target, source)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, string) int); ok {
		r0 = rf(target, source)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(string, string) int); ok {
		r1 = rf(target, source)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(target, source)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetUsersList provides a mock function with given fields:
func (_m *Usecase) GetUsersList() ([]*entity.User, error) {
	ret := _m.Called()

	var r0 []*entity.User
	if rf, ok := ret.Get(0).(func() []*entity.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.User)
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

// Run provides a mock function with given fields: refreshTimeSec
func (_m *Usecase) Run(refreshTimeSec string) {
	_m.Called(refreshTimeSec)
}

// SaveInfoAboutTrain provides a mock function with given fields: trainID
func (_m *Usecase) SaveInfoAboutTrain(trainID string) (string, error) {
	ret := _m.Called(trainID)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(trainID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(trainID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveTrainInUser provides a mock function with given fields: user, trainID
func (_m *Usecase) SaveTrainInUser(user string, trainID string) error {
	ret := _m.Called(user, trainID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(user, trainID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUserTrainInfo provides a mock function with given fields: user
func (_m *Usecase) UpdateUserTrainInfo(user *entity.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UsersCount provides a mock function with given fields:
func (_m *Usecase) UsersCount() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
