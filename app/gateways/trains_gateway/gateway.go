package trains_gateway

import (
	"rzd/app/entity"
)

type TrainsGateway interface {
	Create(user *entity.Train) (string, error)
	Delete(user *entity.Train) error
	Update(user *entity.Train) error
	ReadOne(filter *entity.Train) (*entity.Train, error)
	ReadMany() ([]*entity.Train, error)
	ReadSection(start, end int64) ([]*entity.Train, error)
}
