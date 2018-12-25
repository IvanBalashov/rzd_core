package usecase

import "rzd/app/entity"

type Usecase interface {
	GetSeats(args entity.RouteArgs) error
}
