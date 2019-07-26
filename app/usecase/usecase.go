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
	AddUser(user *entity.User) (bool, error)
	SaveTrainInUser(user string, trainID string) error
	UpdateUserTrainInfo(user *entity.User) error
	DeleteUser(user *entity.User) error
	GetUsersList() ([]entity.User, error)
	UsersCount() (int, error)
	CheckUsers(start, end int64) ([]*entity.User, error)
}
