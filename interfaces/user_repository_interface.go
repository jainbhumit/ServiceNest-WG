package interfaces

import (
	"serviceNest/model"
)

type UserRepository interface {
	SaveUser(user *model.User) error
	GetUserByID(userID string) (*model.User, error)
	UpdateUser(updatedUser *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	DeActivateUser(userID string) error
	GetSecurityAnswerByEmail(userEmail string) (*string, error)
	UpdatePassword(userEmail, updatedPassword string) error
}
