package repository_test

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"serviceNest/model"
	"serviceNest/repository"
	"testing"
)

func TestGetAllServices(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "provider_id", "category"}).
		AddRow("1", "Service A", "Description A", 100.0, "provider123", "Category A").
		AddRow("2", "Service B", "Description B", 200.0, nil, "Category B")

	mock.ExpectQuery("SELECT id, name, description, price, provider_id, category FROM services").
		WillReturnRows(rows)

	services, err := repo.GetAllServices()

	assert.NoError(t, err)
	assert.Len(t, services, 2)
	assert.Equal(t, "Service A", services[0].Name)
	assert.Equal(t, "", services[1].ProviderID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetServiceByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	serviceID := "service123"
	row := sqlmock.NewRows([]string{"id", "name", "description", "price", "provider_id", "category"}).
		AddRow(serviceID, "Service A", "Description A", 100.0, "provider123", "Category A")

	mock.ExpectQuery("SELECT id, name, description, price, provider_id, category FROM services WHERE id = ?").
		WithArgs(serviceID).
		WillReturnRows(row)

	service, err := repo.GetServiceByID(serviceID)

	assert.NoError(t, err)
	assert.Equal(t, "Service A", service.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveService(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	service := model.Service{
		ID:          "service123",
		Name:        "Service A",
		Description: "Description A",
		Price:       100.0,
		ProviderID:  "provider123",
		Category:    "Category A",
	}

	mock.ExpectExec("INSERT INTO services").
		WithArgs(service.ID, service.Name, service.Description, service.Price, service.ProviderID, service.Category).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveService(service)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestRemoveService(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	serviceID := "service123"

	mock.ExpectExec("DELETE FROM services WHERE id = ?").
		WithArgs(serviceID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.RemoveService(serviceID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetServiceByProviderID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	providerID := "provider123"
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "provider_id", "category"}).
		AddRow("1", "Service A", "Description A", 100.0, providerID, "Category A").
		AddRow("2", "Service B", "Description B", 200.0, providerID, "Category B")

	mock.ExpectQuery("SELECT id, name, description, price, provider_id, category FROM services WHERE provider_id = ?").
		WithArgs(providerID).
		WillReturnRows(rows)

	services, err := repo.GetServiceByProviderID(providerID)

	assert.NoError(t, err)
	assert.Len(t, services, 2)
	assert.Equal(t, "Service A", services[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestSaveAllServices_Success(t *testing.T) {
	// Mock the database and create a new ServiceRepository instance
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	// Mock services data
	services := []model.Service{
		{ID: "service1", Name: "Plumbing", Description: "Fix plumbing issues", Price: 100.0, ProviderID: "provider1", Category: "Home"},
		{ID: "service2", Name: "Electrical", Description: "Fix electrical issues", Price: 200.0, ProviderID: "provider2", Category: "Maintenance"},
	}

	// Mock transaction and prepared statement
	mock.ExpectBegin()

	query := regexp.QuoteMeta("INSERT INTO services (id, name, description, price, provider_id, category) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description), price=VALUES(price), provider_id=VALUES(provider_id), category=VALUES(category)")
	prep := mock.ExpectPrepare(query)

	for _, service := range services {
		prep.ExpectExec().
			WithArgs(service.ID, service.Name, service.Description, service.Price, service.ProviderID, service.Category).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	// Call the function
	err = repo.SaveAllServices(services)
	assert.NoError(t, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSaveAllServices_FailureOnPrepare(t *testing.T) {
	// Mock the database and create a new ServiceRepository instance
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	// Mock services data
	services := []model.Service{
		{ID: "service1", Name: "Plumbing", Description: "Fix plumbing issues", Price: 100.0, ProviderID: "provider1", Category: "Home"},
	}

	// Mock transaction and failure in Prepare statement
	mock.ExpectBegin()

	query := regexp.QuoteMeta("INSERT INTO services (id, name, description, price, provider_id, category) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description), price=VALUES(price), provider_id=VALUES(provider_id), category=VALUES(category)")
	mock.ExpectPrepare(query).WillReturnError(sql.ErrConnDone)

	mock.ExpectRollback()

	// Call the function
	err = repo.SaveAllServices(services)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSaveAllServices_FailureOnExec(t *testing.T) {
	// Mock the database and create a new ServiceRepository instance
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	// Mock services data
	services := []model.Service{
		{ID: "service1", Name: "Plumbing", Description: "Fix plumbing issues", Price: 100.0, ProviderID: "provider1", Category: "Home"},
	}

	// Mock transaction and prepared statement
	mock.ExpectBegin()

	query := regexp.QuoteMeta("INSERT INTO services (id, name, description, price, provider_id, category) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description), price=VALUES(price), provider_id=VALUES(provider_id), category=VALUES(category)")
	prep := mock.ExpectPrepare(query)

	// Mock Exec failure
	prep.ExpectExec().WithArgs("service1", "Plumbing", "Fix plumbing issues", 100.0, "provider1", "Home").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectRollback()

	// Call the function
	err = repo.SaveAllServices(services)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func TestGetServiceByName_Success(t *testing.T) {
	// Mock the database and create a new ServiceRepository instance
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	// Mock service data
	expectedService := model.Service{
		ID:          "service123",
		Name:        "Plumbing",
		Description: "Fix water issues",
		Price:       100.0,
		ProviderID:  "provider123",
		Category:    "Home",
	}

	// Mock the query result
	query := regexp.QuoteMeta("SELECT id, name, description, price, provider_id, category FROM services WHERE name = ?")
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "provider_id", "category"}).
		AddRow(expectedService.ID, expectedService.Name, expectedService.Description, expectedService.Price, expectedService.ProviderID, expectedService.Category)
	mock.ExpectQuery(query).WithArgs("Plumbing").WillReturnRows(rows)

	// Call the function
	actualService, err := repo.GetServiceByName("Plumbing")
	assert.NoError(t, err)
	assert.NotNil(t, actualService)

	// Assert the returned service matches the expected service
	assert.Equal(t, &expectedService, actualService)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetServiceByName_NoRows(t *testing.T) {
	// Mock the database and create a new ServiceRepository instance
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	// Mock the query with no matching rows
	query := regexp.QuoteMeta("SELECT id, name, description, price, provider_id, category FROM services WHERE name = ?")
	mock.ExpectQuery(query).WithArgs("NonExistingService").WillReturnError(sql.ErrNoRows)

	// Call the function
	actualService, err := repo.GetServiceByName("NonExistingService")

	// Assert that the service was not found and the error is correct
	assert.Nil(t, actualService)
	assert.EqualError(t, err, "service not found")

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetServiceByName_QueryError(t *testing.T) {
	// Mock the database and create a new ServiceRepository instance
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	// Mock a query error
	query := regexp.QuoteMeta("SELECT id, name, description, price, provider_id, category FROM services WHERE name = ?")
	mock.ExpectQuery(query).WithArgs("Plumbing").WillReturnError(errors.New("database error"))

	// Call the function
	actualService, err := repo.GetServiceByName("Plumbing")

	// Assert that an error was returned
	assert.Nil(t, actualService)
	assert.EqualError(t, err, "database error")

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func TestUpdateService_Success(t *testing.T) {
	// Create a new SQL mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock database
	repo := repository.NewServiceRepository(db)

	// Prepare mock expectation
	query := regexp.QuoteMeta("UPDATE services SET name = ?, description = ?, price = ? WHERE provider_id= ? AND id=?")
	mock.ExpectExec(query).
		WithArgs("ServiceName", "ServiceDescription", 100.0, "provider1", "service1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create the service to update
	service := model.Service{
		ID:          "service1",
		Name:        "ServiceName",
		Description: "ServiceDescription",
		Price:       100.0,
	}

	// Call the method
	err = repo.UpdateService("provider1", service)

	// Assert the expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateService_NoRowsAffected(t *testing.T) {
	// Create a new SQL mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock database
	repo := repository.NewServiceRepository(db)

	// Prepare mock expectation
	query := regexp.QuoteMeta("UPDATE services SET name = ?, description = ?, price = ? WHERE provider_id= ? AND id=?")
	mock.ExpectExec(query).
		WithArgs("ServiceName", "ServiceDescription", 100.0, "provider1", "service1").
		WillReturnResult(sqlmock.NewResult(1, 0))

	// Create the service to update
	service := model.Service{
		ID:          "service1",
		Name:        "ServiceName",
		Description: "ServiceDescription",
		Price:       100.0,
	}

	// Call the method
	err = repo.UpdateService("provider1", service)

	// Assert the expectations
	assert.EqualError(t, err, "The service ID may not exist.")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateService_Error(t *testing.T) {
	// Create a new SQL mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock database
	repo := repository.NewServiceRepository(db)
	// Prepare mock expectation
	query := regexp.QuoteMeta("UPDATE services SET name = ?, description = ?, price = ? WHERE provider_id= ? AND id=?")
	mock.ExpectExec(query).
		WithArgs("ServiceName", "ServiceDescription", 100.0, "provider1", "service1").
		WillReturnError(errors.New("some database error"))

	// Create the service to update
	service := model.Service{
		ID:          "service1",
		Name:        "ServiceName",
		Description: "ServiceDescription",
		Price:       100.0,
	}

	// Call the method
	err = repo.UpdateService("provider1", service)

	// Assert the expectations
	assert.EqualError(t, err, "some database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestRemoveServiceByProviderID_Success(t *testing.T) {
	// Create a new SQL mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock database
	repo := repository.NewServiceRepository(db)
	query := regexp.QuoteMeta("DELETE FROM services WHERE id = ? AND provider_id = ?")
	// Prepare mock expectation
	mock.ExpectExec(query).
		WithArgs("service1", "provider1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err = repo.RemoveServiceByProviderID("provider1", "service1")

	// Assert the expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveServiceByProviderID_NoRowsAffected(t *testing.T) {
	// Create a new SQL mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock database
	repo := repository.NewServiceRepository(db)

	query := regexp.QuoteMeta("DELETE FROM services WHERE id = ? AND provider_id = ?")
	// Prepare mock expectation
	mock.ExpectExec(query).
		WithArgs("service1", "provider1").
		WillReturnResult(sqlmock.NewResult(1, 0))

	// Call the method
	err = repo.RemoveServiceByProviderID("provider1", "service1")

	// Assert the expectations
	assert.EqualError(t, err, "Invalid service ID")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveServiceByProviderID_Error(t *testing.T) {
	// Create a new SQL mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock database
	repo := repository.NewServiceRepository(db)
	query := regexp.QuoteMeta("DELETE FROM services WHERE id = ? AND provider_id = ?")
	// Prepare mock expectation
	mock.ExpectExec(query).
		WithArgs("service1", "provider1").
		WillReturnError(errors.New("some database error"))

	// Call the method
	err = repo.RemoveServiceByProviderID("provider1", "service1")

	// Assert the expectations
	assert.EqualError(t, err, "some database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetServiceByProviderIDAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRepository(db)

	providerID := "provider123"

	t.Run("Successful Query Execution", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "provider_id", "category"}).
			AddRow("1", "Service A", "Description A", 100.0, providerID, "Category A").
			AddRow("2", "Service B", "Description B", 200.0, providerID, "Category B")

		query := "SELECT id, name, description, price, provider_id, category FROM services WHERE provider_id = ?"
		mock.ExpectQuery(query).
			WithArgs(providerID).
			WillReturnRows(rows)

		services, err := repo.GetServiceByProviderID(providerID)

		assert.NoError(t, err)
		assert.Len(t, services, 2)
		assert.Equal(t, "Service A", services[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database Query Error", func(t *testing.T) {
		query := "SELECT id, name, description, price, provider_id, category FROM services WHERE provider_id = ?"
		mock.ExpectQuery(query).
			WithArgs(providerID).
			WillReturnError(errors.New("query error"))

		services, err := repo.GetServiceByProviderID(providerID)

		assert.Error(t, err)
		assert.Nil(t, services)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No Rows Returned", func(t *testing.T) {
		query := "SELECT id, name, description, price, provider_id, category FROM services WHERE provider_id = ?"
		mock.ExpectQuery(query).
			WithArgs(providerID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "provider_id", "category"}))

		services, err := repo.GetServiceByProviderID(providerID)

		assert.NoError(t, err)
		assert.Nil(t, services)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
