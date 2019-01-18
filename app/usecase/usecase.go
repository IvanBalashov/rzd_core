package usecase

import "rzd/app/entity"

type Usecase interface {
	GetSeats(ids []int) ([]entity.Train, error)
}
