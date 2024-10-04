package service

import (
	"errors"
	"fmt"
	"serviceNest/interfaces"
	"serviceNest/model"
	"serviceNest/util"
)

type UserService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) interfaces.UserService {
	return &UserService{userRepo: userRepo}
}

// View User

func (s *UserService) ViewProfileByID(userID string) (*model.User, error) {
	//if err := s.userRepo.EnsureConnection(); err != nil {
	//	return nil, err
	//}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserService) UpdateUser(userID string, newEmail, newPassword, newAddress, newPhone *string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("could not find user: %v", err)
	}

	// Update email
	if newEmail != nil {
		if err := util.ValidateEmail(*newEmail); err != nil {
			return err
		}
		existingUser, err := s.userRepo.GetUserByEmail(*newEmail)
		if err == nil && existingUser.ID != userID {
			return errors.New("email already in use by another user")
		}
		user.Email = *newEmail
	}

	// Update password
	if newPassword != nil {
		if err := util.ValidatePassword(*newPassword); err != nil {
			return err
		}
		user.Password = *newPassword
	}

	// Update contact
	if newPhone != nil {
		if err := util.ValidatePhoneNumber(*newPhone); err != nil {
			return err
		}
		user.Contact = *newPhone
	}
	// Update address
	if newAddress != nil {
		user.Address = *newAddress
	}

	// Save the updated user back to the repository_test
	if err := s.userRepo.UpdateUser(user); err != nil {
		return fmt.Errorf("could not update user: %v", err)
	}

	return nil
}

func (s *UserService) CheckUserExists(email string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}
	return user, nil
}

func (s *UserService) CreateUser(user *model.User) error {
	user.ID = GetUniqueID()
	err := s.userRepo.SaveUser(user)
	if err != nil {
		return fmt.Errorf("could not save user: %v", err)
	}
	return nil
}
