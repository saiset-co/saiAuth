package internal

import "github.com/Limpid-LLC/go-auth/internal/entities"

func (is InternalService) createUser(email string, phone string, password string, data interface{}) *entities.User {
	return &entities.User{
		Email:          email,
		Phone:          phone,
		HashedPassword: is.hashAndSaltPassword(password),
		Roles:          nil,
		Data:           data,
	}
}
