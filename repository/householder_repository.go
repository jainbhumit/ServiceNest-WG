package repository

import (
	"database/sql"
	"serviceNest/interfaces"
	"serviceNest/model"
)

type MySQLHouseholderRepository struct {
	db *sql.DB
}

// NewHouseholderRepository creates a new instance of MySQLHouseholderRepository
func NewHouseholderRepository(client *sql.DB) interfaces.HouseholderRepository {
	return &MySQLHouseholderRepository{
		db: client,
	}
}

func (repo *MySQLHouseholderRepository) SaveHouseholder(householder *model.Householder) error {
	query := "INSERT INTO users (id, name, email, password, role, address, contact, latitude, longitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := repo.db.Exec(query, householder.ID, householder.Name, householder.Email, householder.Password, householder.Role, householder.Address, householder.Contact, householder.Latitude, householder.Longitude)
	return err
}

func (repo *MySQLHouseholderRepository) GetHouseholderByID(id string) (*model.Householder, error) {
	query := "SELECT id, name, email, password, role, address, contact, latitude, longitude FROM users WHERE id = ?"
	row := repo.db.QueryRow(query, id)

	var householder model.Householder
	err := row.Scan(&householder.ID, &householder.Name, &householder.Email, &householder.Password, &householder.Role, &householder.Address, &householder.Contact, &householder.Latitude, &householder.Longitude)
	if err != nil {
		return nil, err
	}

	return &householder, nil
}
