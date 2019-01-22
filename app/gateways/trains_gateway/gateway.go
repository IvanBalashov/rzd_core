package trains_gateway

import "rzd/app/entity"

type TrainsGateway interface {
	Create(train entity.Train) error
	ReadOne() (entity.Train, error)
	ReadMany(ids []int) ([]entity.Train, error)
	Update(train entity.Train) error
	Delete(train entity.Train) error
}
