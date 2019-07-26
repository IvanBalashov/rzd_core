package trains_gateway

import (
	"rzd/app/entity"
)

type TrainsGateway interface {
	Create(train entity.Train) (string, error)
	ReadOne(trainID string) (entity.Train, error)
	ReadMany() ([]entity.Train, error)
	Update(train entity.Train) error
	Delete(train entity.Train) error
}
