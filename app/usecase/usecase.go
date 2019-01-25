package usecase

import "rzd/app/entity"

type Usecase interface {
	GetSeats(args entity.RouteArgs) ([]entity.Train, error)
	GetCodes(target, source string) (int, int, error)
	SaveTrain()
	Run()
}
