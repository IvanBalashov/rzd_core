package usecase

import (
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"rzd/app/entity"
	"rzd/mocks"
	"testing"
)

var logs = make(chan string)
var app = App{
	Trains:  &mocks.TrainsGateway{},
	Users:   &mocks.UsersGateway{},
	Routes:  &mocks.RzdGateway{},
	Cache:   &mocks.CacheGateway{},
	LogChan: logs,
}

func init() {

}

func TestApp_GetStationCodes(t *testing.T) {
	expectedCode1 := 0
	expectedCode2 := 0
	rzdMock := mocks.RzdGateway{}
	app.Routes = &rzdMock

	rzdMock.On("GetDirectionsCode", "").Return(0, nil)

	code1, code2, err := app.GetStationCodes("", "")
	if err != nil {
		t.Errorf("error in GetStationCodes - %s\n", err)
		t.FailNow()
	}

	assert.Equal(t, expectedCode1, code1)
	assert.Equal(t, expectedCode2, code2)
}

func TestApp_GetInfoAboutTrains(t *testing.T) {
	args := entity.RouteArgs{}
	expectedTrains := []entity.Train{}
	rzdMock := mocks.RzdGateway{}
	cache := mocks.CacheGateway{}

	app.Cache = &cache
	app.Routes = &rzdMock

	cache.On("Set", "", []byte{}).Return(nil)
	rzdMock.On("GetRid", entity.RidArgs{}).Return(entity.Rid{}, []*http.Cookie{}, nil)
	rzdMock.On("GetRoutes", entity.RouteArgs{Rid: "0"}, []*http.Cookie{}).Return(entity.Route{
		Result: "OK",
		Tp: []entity.Tp{
			0: {
				List: []entity.List{},
			},
		},
	}, nil)

	go func() {
		for msg := range logs {
			log.Println(msg)
		}
	}()

	trains, err := app.GetInfoAboutTrains(args)
	if err != nil {
		t.Errorf("error in GetInfoAboutTrains - %s\n", err)
		t.FailNow()
	}

	assert.Equal(t, expectedTrains, trains)
}

func TestApp_GenerateTrainsList(t *testing.T) {
	route := entity.Route{
		Result: "OK",
		Tp: []entity.Tp{
			0: {
				List: []entity.List{},
			},
		},
	}
	args := entity.RouteArgs{}
	expectedTrains := []entity.Train{}

	trains, err := app.GenerateTrainsList(route, args)
	if err != nil {
		t.Errorf("error in GenerateTrainsList - %s\n", err)
		t.FailNow()
	}

	assert.Equal(t, trains, expectedTrains)
}

func TestApp_GetUsersList(t *testing.T) {
	expectedUsers := []entity.User{}
	usersMock := mocks.UsersGateway{}
	app.Users = &usersMock

	usersMock.On("ReadMany").Return([]entity.User{}, nil)

	users, err := app.GetUsersList()
	if err != nil {
		t.Errorf("error in GetUsersList - %s\n", err)
		t.FailNow()
	}

	assert.Equal(t, users, expectedUsers)
}

func TestApp_AddUser(t *testing.T) {
	user := entity.User{}
	usersMock := mocks.UsersGateway{}
	app.Users = &usersMock

	usersMock.On("Create", user).Return(nil)

	err := app.AddUser(user)
	if err != nil {
		t.Errorf("error in AddUser - %s\n", err)
		t.FailNow()
	}
}

func TestApp_DeleteUser(t *testing.T) {
	user := entity.User{}
	usersMock := mocks.UsersGateway{}
	app.Users = &usersMock

	usersMock.On("Delete", user).Return(nil)

	err := app.DeleteUser(user)
	if err != nil {
		t.Errorf("error in AddUser - %s\n", err)
		t.FailNow()
	}
}

func TestApp_SaveTrainInUser(t *testing.T) {
	user := entity.User{}
	usersMock := mocks.UsersGateway{}
	app.Users = &usersMock

	usersMock.On("ReadOne", entity.User{}).Return(entity.User{}, nil)
	usersMock.On("Update", entity.User{}).Return( nil)

	err := app.SaveTrainInUser(user, "123321")
	if err != nil {
		t.Errorf("error in AddUser - %s\n", err)
		t.FailNow()
	}

}

func TestApp_CheckAndRefreshTrainInfo(t *testing.T) {
	train := entity.Train{}
	expectedResult := false
	rzdMock := mocks.RzdGateway{}

	rzdMock.On("GetRid", entity.RidArgs{}).Return(entity.Rid{RID: 0}, []*http.Cookie{}, nil)
	rzdMock.On("GetInfoAboutOneTrain", entity.Train{
		QueryArgs: entity.RouteArgs{Rid: "0"},
	}, []*http.Cookie{}).Return(entity.Route{}, nil)
	app.Routes = &rzdMock

	result := app.CheckAndRefreshTrainInfo(train)

	assert.Equal(t, result, expectedResult)
}

func TestApp_CheckUsers(t *testing.T)  {
	start := 0
	end := 10
	expectedUsers := []entity.User{}
	usersMock := mocks.UsersGateway{}
	app.Users = &usersMock

	usersMock.On("ReadSection", start, end).Return([]entity.User{}, nil)

	users, err := app.CheckUsers(start, end)
	if err != nil {
		t.Errorf("error in CheckUsers - %s\n", err)
		t.FailNow()
	}

	assert.Equal(t, users, expectedUsers)
}