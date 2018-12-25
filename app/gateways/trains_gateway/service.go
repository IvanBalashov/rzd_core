package trains_gateway

import (
	"io"
	"rzd/app/entity"
)

//EXAMPLE
type ServiceTrains struct {
	Connection io.Reader
}

func NewService(reader *io.Reader) ServiceTrains {
	return ServiceTrains{Connection: *reader}
}

func (s *ServiceTrains) Create(user entity.User) error {
	panic("IMPLIMENT ME!!!")
	return nil
}

func (s *ServiceTrains) Read(offset, limit int) ([]entity.User, error) {
	panic("IMPLIMENT ME!!!")
	return nil, nil
}

func (s *ServiceTrains) Update(user entity.User) error {
	panic("IMPLIMENT ME!!!")
	return nil
}

func (s *ServiceTrains) Delete(user entity.User) error {
	panic("IMPLIMENT ME!!!")
	return nil
}
