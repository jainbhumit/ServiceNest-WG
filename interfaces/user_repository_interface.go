package interfaces

import (
	"serviceNest/model"
)

type UserRepository interface {
	SaveUser(user *model.User) error
	GetUserByID(userID string) (*model.User, error)
	UpdateUser(updatedUser *model.User) error
	GetUserByEmail(email string) (*model.User, error)
}
