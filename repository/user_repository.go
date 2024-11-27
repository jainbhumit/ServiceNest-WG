package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"serviceNest/config"
	"serviceNest/errs"
	"serviceNest/interfaces"
	"serviceNest/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) SaveUser(user *model.User) error {
	column := []string{"id", "name", "email", "password", "role", "address", "contact", "security_answer", "is_active"}
	query := config.InsertQuery("users", column)

	_, err := repo.db.Exec(query, user.ID, user.Name, user.Email, user.Password, user.Role, user.Address, user.Contact, user.SecurityAnswer, true)
	return err
}

func (repo *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	column := []string{"id", "name", "email", "password", "role", "address", "contact", "is_active"}
	query := config.SelectQuery("users", "email", "", column)

	row := repo.db.QueryRow(query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Address, &user.Contact, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(errs.UserNotFound)
		}
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdateUser(updatedUser *model.User) error {
	// Ensure the new email doesn't already exist in the system
	existingUser, err := repo.GetUserByEmail(updatedUser.Email)
	if err == nil && existingUser.ID != updatedUser.ID {
		return fmt.Errorf(errs.EmailAlreadyUse)
	}
	column := []string{"name", "email", "password", "role", "address", "contact"}
	query := config.UpdateQuery("users", "id", "", column)

	_, err = repo.db.Exec(query, updatedUser.Name, updatedUser.Email, updatedUser.Password, updatedUser.Role, updatedUser.Address, updatedUser.Contact, updatedUser.ID)
	return err
}

func (repo *UserRepository) GetUserByID(userID string) (*model.User, error) {
	column := []string{"id", "name", "email", "password", "role", "address", "contact"}
	query := config.SelectQuery("users", "id", "", column)

	row := repo.db.QueryRow(query, userID)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Address, &user.Contact)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(errs.UserNotFound)
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) DeActivateUser(userID string) error {
	column := []string{"is_active"}
	query := config.UpdateQuery("users", "id", "", column)

	_, err := repo.db.Exec(query, false, userID)
	return err
}

func (repo *UserRepository) GetSecurityAnswerByEmail(userEmail string) (*string, error) {
	column := []string{"security_answer"}
	query := config.SelectQuery("users", "email", "", column)

	row := repo.db.QueryRow(query, userEmail)
	var securityAnswer string
	err := row.Scan(&securityAnswer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(errs.UserNotFound)
		}
		return nil, err
	}
	return &securityAnswer, nil
}

func (repo *UserRepository) UpdatePassword(userEmail, updatedPassword string) error {
	column := []string{"password"}
	query := config.UpdateQuery("users", "email", "", column)

	result, err := repo.db.Exec(query, updatedPassword, userEmail)
	if err != nil {
		return err
	}

	// Check affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New(errs.UserNotFound)
	}

	return nil
}
