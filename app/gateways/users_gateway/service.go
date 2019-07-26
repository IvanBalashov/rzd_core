package users_gateway

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
type ServiceUsers struct {
	Connection io.Reader
}

func NewService(reader *io.Reader) ServiceUsers {
	return ServiceUsers{Connection: *reader}
}

func (s *ServiceUsers) Create(user entity.User) error {
	panic("Service:Gateways->Trains_Gateway->Create: Not implemented method")

	return nil
}

func (s *ServiceUsers) ReadOne() (entity.User, error) {
	panic("Service:Gateways->Trains_Gateway->ReadOne: Not implemented method")

	return entity.User{}, nil
}

func (s *ServiceUsers) ReadMany(ids []int) ([]entity.User, error) {
	panic("Service:Gateways->Trains_Gateway->ReadMany: Not implemented method")

	return nil, nil
}

func (s *ServiceUsers) Update(user entity.User) error {
	panic("Service:Gateways->Trains_Gateway->Update: Not implemented method")

	return nil
}

func (s *ServiceUsers) Delete(user entity.User) error {
	panic("Service:Gateways->Trains_Gateway->Delete: Not implemented method")

	return nil
}
