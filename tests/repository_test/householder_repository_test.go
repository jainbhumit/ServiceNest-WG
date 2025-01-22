package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"serviceNest/model"
	"serviceNest/repository"
)

func TestSaveHouseholder(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a new instance of the repository
	repo := repository.NewHouseholderRepository(db)

	// Create a test householder
	householder := &model.Householder{
		User: model.User{
			ID:       "1",
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
			Role:     "householder",
			Address:  "123 Main St",
			Contact:  "1234567890",
		},
	}

	// Set up the expectation for the INSERT query
	mock.ExpectExec("INSERT INTO users").
		WithArgs(householder.ID, householder.Name, householder.Email, householder.Password, householder.Role, householder.Address, householder.Contact).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Return result as if one row was inserted

	// Call the method under test
	err = repo.SaveHouseholder(householder)

	// Assert that there was no error and expectations were met
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet()) // Ensure that all expectations are met
}

func TestGetHouseholderByID(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a new instance of the repository
	repo := repository.NewHouseholderRepository(db)

	// Create the expected householder
	expectedHouseholder := &model.Householder{
		User: model.User{
			ID:       "1",
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
			Role:     "householder",
			Address:  "123 Main St",
			Contact:  "1234567890",
		},
	}

	// Set up the expectation for the SELECT query
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "address", "contact"}).
		AddRow(expectedHouseholder.ID, expectedHouseholder.Name, expectedHouseholder.Email, expectedHouseholder.Password, expectedHouseholder.Role, expectedHouseholder.Address, expectedHouseholder.Contact)

	mock.ExpectQuery("SELECT id, name, email, password, role, address, contact FROM users WHERE id = ?").
		WithArgs(expectedHouseholder.ID).
		WillReturnRows(rows)

	// Call the method under test
	householder, err := repo.GetHouseholderByID("1")

	// Assert that there was no error and the householder matches the expected result
	assert.NoError(t, err)
	assert.Equal(t, expectedHouseholder, householder)
	assert.NoError(t, mock.ExpectationsWereMet()) // Ensure that all expectations are met
}
