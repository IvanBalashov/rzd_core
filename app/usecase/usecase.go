package usecase

import "rzd/app/entity"

type Usecase interface {
	GetInfoAboutTrains(args entity.RouteArgs) ([]entity.Train, error)
	GetStationCodes(target, source string) (int, int, error)
	SaveInfoAboutTrain(trainID string) error
	Run(refreshTimeSec string)
}
