package repository_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"serviceNest/model"
	"serviceNest/repository"
	"testing"
)

func TestSaveUser(t *testing.T) {
	// Step 1: Create a new mock database connection using sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err) // Ensure no error in mock creation
	defer db.Close()

	// Step 2: Initialize the UserRepository with the mocked db
	repo := repository.NewUserRepository(db)

	// Step 3: Define the mock behavior for the INSERT INTO users query
	mock.ExpectExec("INSERT INTO users").
		WithArgs("123", "John Doe", "john@example.com", "hashed_password", "Householder", "123 Main St", "1234567890").
		WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate successful insert

	// Step 4: Define the user object that we want to save
	user := &model.User{
		ID:       "123",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
		Role:     "Householder",
		Address:  "123 Main St",
		Contact:  "1234567890",
	}

	// Step 5: Call the SaveUser method from the repository
	err = repo.SaveUser(user)

	// Step 6: Assert that no error occurred
	assert.NoError(t, err)

	// Step 7: Ensure all expectations set by sqlmock were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize the repository
	repo := repository.NewUserRepository(db)

	// Mock row returned by query
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "address", "contact"}).
		AddRow("123", "John Doe", "john@example.com", "hashed_password", "user", "123 Main St", "1234567890")

	// Expect the query with the provided email
	mock.ExpectQuery("SELECT id, name, email, password").
		WithArgs("john@example.com").
		WillReturnRows(rows)

	// Call GetUserByEmail and assert no error
	user, err := repo.GetUserByEmail("john@example.com")
	assert.NoError(t, err)

	// Assert the user details are as expected
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "hashed_password", user.Password)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_UserNotFound(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize the repository
	repo := repository.NewUserRepository(db)

	// Expect the query but return no rows
	mock.ExpectQuery("SELECT id, name, email, password").
		WithArgs("john@example.com").
		WillReturnError(sql.ErrNoRows)

	// Call GetUserByEmail and assert the error
	user, err := repo.GetUserByEmail("john@example.com")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "user not found")

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize the repository
	repo := repository.NewUserRepository(db)

	// Mock row for existing user check
	existingUserRows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "address", "contact", "latitude", "longitude"}).
		AddRow("123", "John Doe", "john@example.com", "hashed_password", "user", "123 Main St", "1234567890", 12.34, 56.78)

	// Expect GetUserByEmail query and return existing user
	mock.ExpectQuery("SELECT id, name, email, password").
		WithArgs("john@example.com").
		WillReturnRows(existingUserRows)

	// Mock Exec for updating user
	mock.ExpectExec("UPDATE users").
		WithArgs("John Updated", "john@example.com", "new_hashed_password", "admin", "123 New St", "0987654321", "123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create updated user
	updatedUser := &model.User{
		ID:       "123",
		Name:     "John Updated",
		Email:    "john@example.com",
		Password: "new_hashed_password",
		Role:     "admin",
		Address:  "123 New St",
		Contact:  "0987654321",
	}

	// Call UpdateUser and assert no error
	err = repo.UpdateUser(updatedUser)
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize the repository
	repo := repository.NewUserRepository(db)

	// Mock row returned by query
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "address", "contact"}).
		AddRow("123", "John Doe", "john@example.com", "hashed_password", "user", "123 Main St", "1234567890")

	// Expect the query with the provided user ID
	mock.ExpectQuery("SELECT id, name, email, password").
		WithArgs("123").
		WillReturnRows(rows)

	// Call GetUserByID and assert no error
	user, err := repo.GetUserByID("123")
	assert.NoError(t, err)

	// Assert the user details are as expected
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "hashed_password", user.Password)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_UserNotFound(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize the repository
	repo := repository.NewUserRepository(db)

	// Expect the query but return no rows
	mock.ExpectQuery("SELECT id, name, email, password").
		WithArgs("123").
		WillReturnError(sql.ErrNoRows)

	// Call GetUserByID and assert the error
	user, err := repo.GetUserByID("123")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "user not found")

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
