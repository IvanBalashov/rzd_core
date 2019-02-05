package users_gateway

import "rzd/app/entity"

type UsersGateway interface {
	Create(user entity.User) error
	ReadOne() (entity.User, error)
	ReadMany() ([]entity.User, error)
	Update(user entity.User) error
	Delete(user entity.User) error
	ReadSection(start, end int) ([]entity.User, error)
}
