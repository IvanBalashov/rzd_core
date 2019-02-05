package usecase

import "rzd/app/entity"

// TODO: split on two interfaces, for trains and users.
type Usecase interface {
	GetInfoAboutTrains(args entity.RouteArgs) ([]entity.Train, error)
	GetStationCodes(target, source string) (int, int, error)
	SaveInfoAboutTrain(trainID string) error
	Run(refreshTimeSec string)
	// User Block
	AddUser(user entity.User) error
	UpdateUserTrainInfo(user entity.User) error
	DeleteUser(user entity.User) error
	GetUsersList() ([]entity.User, error)
	CheckUsers(start, end int) (bool, error)
}
