package repository

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
	"serviceNest/model"
)

type UserRepository struct {
	filePath string
}

func NewUserRepository(filePath string) *UserRepository {
	return &UserRepository{filePath: filePath}
}

// SaveUser adds a new user to the repository
func (repo *UserRepository) SaveUser(user model.User) error {
	users, err := repo.loadUsers()
	if err != nil {
		return err
	}

	users = append(users, user)

	return repo.saveUsers(users)
}

// GetUserByEmail retrieves a user by their email
func (repo *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	users, err := repo.loadUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

// UpdateUser updates an existing user's details
func (repo *UserRepository) UpdateUser(updatedUser *model.User) error {
	users, err := repo.loadUsers()
	if err != nil {
		return err
	}

	for i, user := range users {
		if user.ID == updatedUser.ID {
			// Update the user's details
			if updatedUser.Email != "" && user.Email != updatedUser.Email {
				// Ensure the new email doesn't already exist in the system
				if _, err := repo.GetUserByEmail(updatedUser.Email); err == nil {
					return fmt.Errorf("email already in use")
				}
				user.Email = updatedUser.Email
			}
			if updatedUser.Password != "" {
				hashPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
				if err != nil {
					return err
				}
				user.Password = string(hashPassword)
			}
			if updatedUser.Address != "" {
				user.Address = updatedUser.Address
			}
			if updatedUser.Contact != "" {
				user.Contact = updatedUser.Contact
			}
			// Update the slice with the modified user
			users[i] = user
			break
		}
	}

	// Save the updated user list back to the file
	return repo.saveUsers(users)
}

// loadUsers loads users from the file
func (repo *UserRepository) loadUsers() ([]model.User, error) {
	var users []model.User

	file, err := ioutil.ReadFile(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return users, nil
		}
		return nil, err
	}

	err = json.Unmarshal(file, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// saveUsers saves the list of users to the file
func (repo *UserRepository) saveUsers(users []model.User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(repo.filePath, data, 0644)
}

// GetUserByID retrieves a user by their unique ID
func (repo *UserRepository) GetUserByID(userID string) (*model.User, error) {
	users, err := repo.loadUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.ID == userID {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}
