package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"serviceNest/config"
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
	column := []string{"id", "name", "email", "password", "role", "address", "contact", "latitude", "longitude"}
	query := config.InsertQuery("users", column)
	//query := `INSERT INTO users (id, name, email, password, role, address, contact, latitude, longitude)
	//          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := repo.db.Exec(query, user.ID, user.Name, user.Email, user.Password, user.Role, user.Address, user.Contact, user.Latitude, user.Longitude)
	return err
}

func (repo *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	column := []string{"id", "name", "email", "password", "role", "address", "contact", "latitude", "longitude"}
	query := config.SelectQuery("users", "email", "", column)
	//query := `SELECT id, name, email, password, role, address, contact, latitude, longitude FROM users WHERE email = ?`
	row := repo.db.QueryRow(query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Address, &user.Contact, &user.Latitude, &user.Longitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdateUser(updatedUser *model.User) error {
	// Ensure the new email doesn't already exist in the system
	existingUser, err := repo.GetUserByEmail(updatedUser.Email)
	if err == nil && existingUser.ID != updatedUser.ID {
		return fmt.Errorf("email already in use")
	}
	column := []string{"name", "email", "password", "role", "address", "contact", "latitude", "longitude"}
	query := config.UpdateQuery("users", "email", "", column)
	//query := `UPDATE users SET name=?, email=?, password=?, role=?, address=?, contact=?, latitude=?, longitude=? WHERE id=?`
	_, err = repo.db.Exec(query, updatedUser.Name, updatedUser.Email, updatedUser.Password, updatedUser.Role, updatedUser.Address, updatedUser.Contact, updatedUser.Latitude, updatedUser.Longitude, updatedUser.ID)
	return err
}

func (repo *UserRepository) GetUserByID(userID string) (*model.User, error) {
	column := []string{"id", "name", "email", "password", "role", "address", "contact", "latitude", "longitude"}
	query := config.SelectQuery("users", "id", "", column)
	//query := `SELECT id, name, email, password, role, address, contact, latitude, longitude FROM users WHERE id = ?`
	row := repo.db.QueryRow(query, userID)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Address, &user.Contact, &user.Latitude, &user.Longitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}
