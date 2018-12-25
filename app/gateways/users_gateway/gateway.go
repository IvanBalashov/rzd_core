package users_gateway

import "rzd/app/entity"

type UsersGateway interface {
	Create(user entity.User) error
	Read(offset, limit int) ([]entity.User, error)
	Update(user entity.User) error
	Delete(user entity.User) error
}
