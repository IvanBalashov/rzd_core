package usecase

import "rzd/app/entity"

// TODO: split on two interfaces, for trains and users.
type Usecase interface {
	// Trains Block
	GetInfoAboutTrains(args entity.RouteArgs) ([]entity.Train, error)
	GetStationCodes(target, source string) (int, int, error)
	SaveInfoAboutTrain(trainID string) (string, error)
	Run(refreshTimeSec string)
	// User Block
	AddUser(user entity.User) error
	SaveTrainInUser(user entity.User, trainID string) error
	UpdateUserTrainInfo(user entity.User) error
	DeleteUser(user entity.User) error
	GetUsersList() ([]entity.User, error)
	CheckUsers(start, end int) (bool, error)
}
