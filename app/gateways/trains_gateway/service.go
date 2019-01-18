package trains_gateway

import (
	"io"
	"rzd/app/entity"
)

/*
This struct released only for example.
I'm wanna show how CA works with several data flows.
First - sql, second - another place, like remote service.
DON'T USE THIS IN CODE!!!!
*/
type ServiceTrains struct {
	Connection io.Reader
}

func NewService(reader *io.Reader) ServiceTrains {
	return ServiceTrains{Connection: *reader}
}

func (s *ServiceTrains) Create(train entity.Train) error {
	panic("IMPLIMENT ME!!!")
	return nil
}

func (s *ServiceTrains) ReadOne(id int) (entity.Train, error) {
	panic("IMPLIMENT ME!!!")
	return entity.Train{}, nil
}

func (s *ServiceTrains) ReadMany(ids []int) ([]entity.Train, error) {
	panic("IMPLIMENT ME!!!")
	return nil, nil
}

func (s *ServiceTrains) Update(train entity.Train) error {
	panic("IMPLIMENT ME!!!")
	return nil
}

func (s *ServiceTrains) Delete(train entity.Train) error {
	panic("IMPLIMENT ME!!!")
	return nil
}
