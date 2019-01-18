package users_gateway

import (
	"io"
	"rzd/app/entity"
)

/*
This struct released only for example.
I'm wana show how CA works with several data flows.
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
	panic("IMPLIMENT ME!!!")
	return nil
}

func (s *ServiceUsers) Read(offset, limit int) ([]entity.User, error) {
	panic("IMPLIMENT ME!!!")
	return nil, nil
}

func (s *ServiceUsers) Update(user entity.User) error {
	panic("IMPLIMENT ME!!!")
	return nil
}

func (s *ServiceUsers) Delete(user entity.User) error {
	panic("IMPLIMENT ME!!!")
	return nil
}
