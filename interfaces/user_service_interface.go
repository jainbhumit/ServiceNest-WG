package interfaces

import "serviceNest/model"

type UserService interface {
	CreateUser(user *model.User) error
	CheckUserExists(email string) (*model.User, error)
	UpdateUser(userID string, newEmail, newPassword, newAddress, newPhone *string) error
	ViewProfileByID(userID string) (*model.User, error)
}
