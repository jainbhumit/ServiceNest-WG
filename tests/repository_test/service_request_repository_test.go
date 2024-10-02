package repository_test

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"serviceNest/model"
	"serviceNest/repository"
	"testing"
	"time"
)

func convertNullString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func TestSaveServiceRequest(t *testing.T) {
	// Create a mock database and defer its closure
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error while opening stub database connection: %v", err)
	}
	defer db.Close()

	// Create a sample service request
	request := model.ServiceRequest{
		ID:                 "1",
		ServiceName:        "ServiceName",
		HouseholderID:      convertNullString(sql.NullString{String: "1001", Valid: true}),
		HouseholderName:    "John Doe",
		HouseholderAddress: convertNullString(sql.NullString{String: "123 Main St", Valid: true}),
		ServiceID:          "SVC001",
		RequestedTime:      time.Now(),
		ScheduledTime:      time.Now().Add(24 * time.Hour),
		Status:             "Pending",
		ApproveStatus:      false,
	}

	// Expecting a query that inserts the service request into the database
	mock.ExpectExec("INSERT INTO service_requests").
		WithArgs(request.ID, request.HouseholderID, request.HouseholderName, request.HouseholderAddress, request.ServiceID, request.RequestedTime, request.ScheduledTime, request.Status, request.ApproveStatus, request.ServiceName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Initialize the repository and call SaveServiceRequest
	repo := repository.NewServiceRequestRepository(db)
	err = repo.SaveServiceRequest(request)
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
func TestGetServiceRequestByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	requestID := "request123"
	householderID := "householder123"
	householderAddress := "123 Street"

	row := sqlmock.NewRows([]string{"id", "householder_id", "householder_name", "householder_address", "service_id",
		"requested_time", "scheduled_time", "status", "approve_status", "service_name"}).
		AddRow(requestID, &householderID, "John Doe", &householderAddress, "Service A",
			time.Now(), time.Now().Add(24*time.Hour), "Pending", false, "service123")

	mock.ExpectQuery("SELECT service_requests.id, service_requests.householder_id, service_requests.householder_name, service_requests.householder_address, service_requests.service_id, service_requests.requested_time, service_requests.scheduled_time, service_requests.status, service_requests.approve_status, services.name FROM service_requests INNER JOIN services ON service_requests.service_id = services.id WHERE service_requests.id = ?").
		WithArgs(requestID).
		WillReturnRows(row)

	request, err := repo.GetServiceRequestByID(requestID)

	assert.NoError(t, err)
	assert.Equal(t, "request123", request.ID)
	assert.NotNil(t, request.HouseholderID)
	assert.Equal(t, "householder123", *request.HouseholderID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetServiceRequestsByHouseholderID(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	// Initialize the repository with the mock db
	repo := repository.NewServiceRequestRepository(db)

	// Mock rows returned by the query
	rows := sqlmock.NewRows([]string{
		"id", "householder_id", "householder_name", "householder_address", "service_id",
		"requested_time", "scheduled_time", "status", "approve_status", "service_name", "service_provider_id",
		"name", "contact", "address", "price", "rating", "approve",
	}).AddRow(1, "1001", "John Doe", "123 Main St", "2001", []byte("2024-09-01 12:00:00"),
		[]byte("2024-09-02 14:00:00"), "Pending", true, "carpenter", "3001", "Provider 1",
		"1234567890", "456 Provider St", "100.00", 4.5, true)

	// Set up expectation for the query, use the full query and escape special characters
	mock.ExpectQuery(`SELECT service_requests\.id, service_requests\.householder_id, service_requests\.householder_name, service_requests\.householder_address, service_requests\.service_id, service_requests\.requested_time, service_requests\.scheduled_time, service_requests\.status, service_requests\.approve_status, service_requests\.service_name, service_provider_details\.service_provider_id, service_provider_details\.name, service_provider_details\.contact, service_provider_details\.address, service_provider_details\.price, service_provider_details\.rating, service_provider_details\.approve FROM service_requests LEFT JOIN service_provider_details ON service_requests\.id = service_provider_details\.service_request_id WHERE service_requests\.householder_id = \?`).
		WithArgs("1001").
		WillReturnRows(rows)

	// Call the method
	requests, err := repo.GetServiceRequestsByHouseholderID("1001")

	// Assert that no error occurred
	assert.NoError(t, err)
	assert.Len(t, requests, 1)
	assert.Equal(t, "John Doe", requests[0].HouseholderName)
	assert.Equal(t, "Provider 1", requests[0].ProviderDetails[0].Name)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateServiceRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	householderID := "householder123"
	householderAddress := "123 Street"
	// Create a sample service request for testing
	request := &model.ServiceRequest{
		ID:                 "1001",
		HouseholderID:      &householderID,
		HouseholderName:    "John Doe",
		HouseholderAddress: &householderAddress,
		ServiceID:          "service-456",
		RequestedTime:      time.Now(),
		ScheduledTime:      time.Now().Add(24 * time.Hour),
		Status:             "Pending",
		ApproveStatus:      false,
	}

	// Mock the Exec query for the update operation
	mock.ExpectExec("UPDATE service_requests").
		WithArgs(request.HouseholderID, request.HouseholderName, request.HouseholderAddress,
			request.ServiceID, request.RequestedTime, request.ScheduledTime, request.Status,
			request.ApproveStatus, request.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err = repo.UpdateServiceRequest(request)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllServiceRequests(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	// Mock rows
	rows := sqlmock.NewRows([]string{
		"id", "householder_id", "householder_name", "householder_address", "service_id",
		"requested_time", "scheduled_time", "status", "approve_status", "service_name",
		"service_provider_id", "name", "contact", "address", "price", "rating", "approve",
	}).
		AddRow(1, "1001", "John Doe", "123 Main St", "2001",
			[]byte("2024-09-01 12:00:00"), []byte("2024-09-02 14:00:00"),
			"Pending", true, "carpenter", "3001", "Provider 1", "1234567890",
			"456 Provider St", "100.00", 4.5, true)

	// Expect the query to return the rows
	mock.ExpectQuery("SELECT service_requests.id").
		WillReturnRows(rows)

	// Call the method
	requests, err := repo.GetAllServiceRequests()

	// Assert no error and check result
	assert.NoError(t, err)
	assert.Len(t, requests, 1)
	assert.Equal(t, "John Doe", requests[0].HouseholderName)
	assert.Equal(t, "Provider 1", requests[0].ProviderDetails[0].Name)

	// Ensure expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetServiceRequestsByProviderID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	// Mock rows
	rows := sqlmock.NewRows([]string{
		"id", "householder_id", "householder_name", "householder_address", "service_id",
		"requested_time", "scheduled_time", "status", "approve_status",
		"service_provider_id", "name", "contact", "address", "price", "rating", "approve",
	}).
		AddRow(1, "1001", "John Doe", "123 Main St", "2001",
			[]byte("2024-09-01 12:00:00"), []byte("2024-09-02 14:00:00"),
			"Pending", true, "3001", "Provider 1", "1234567890",
			"456 Provider St", "100.00", 4.5, true)

	// Expect the query to return the rows
	mock.ExpectQuery("SELECT service_requests.id").
		WithArgs("3001").
		WillReturnRows(rows)

	// Call the method
	requests, err := repo.GetServiceRequestsByProviderID("3001")

	// Assert no error and check result
	assert.NoError(t, err)
	assert.Len(t, requests, 1)
	assert.Equal(t, "John Doe", requests[0].HouseholderName)
	assert.Equal(t, "Provider 1", requests[0].ProviderDetails[0].Name)

	// Ensure expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetServiceProviderByRequestID(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	// Initialize the repository with the mock db
	repo := repository.NewServiceRequestRepository(db)

	// Mock the row returned by the query
	rows := sqlmock.NewRows([]string{
		"id", "householder_id", "householder_name", "householder_address", "service_id",
		"requested_time", "scheduled_time", "status", "approve_status",
		"service_provider_id", "name", "contact", "address", "price", "rating", "approve",
	}).
		AddRow("1001", "householder-123", "John Doe", "123 Main St", "service-456",
			[]byte("2024-09-01 12:00:00"), []byte("2024-09-02 14:00:00"),
			"Pending", true, "provider-789", "Provider Name", "1234567890",
			"456 Provider St", "100.00", 4.5, true)

	// Set up expectation for the query
	mock.ExpectQuery("SELECT service_requests.id, service_requests.householder_id").
		WithArgs("provider-789", "1001").
		WillReturnRows(rows)

	// Call the method
	result, err := repo.GetServiceProviderByRequestID("1001", "provider-789")

	// Assert that no error occurred
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Validate returned data
	assert.Equal(t, "John Doe", result.HouseholderName)
	assert.Equal(t, "provider-789", result.ProviderDetails[0].ServiceProviderID)
	assert.Equal(t, "Provider Name", result.ProviderDetails[0].Name)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetServiceProviderByRequestID_NoRows(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	// Initialize the repository with the mock db
	repo := repository.NewServiceRequestRepository(db)

	// Set up expectation for no rows returned
	mock.ExpectQuery("SELECT service_requests.id, service_requests.householder_id").
		WithArgs("provider-789", "1001").
		WillReturnRows(sqlmock.NewRows(nil)) // Empty result set

	// Call the method
	result, err := repo.GetServiceProviderByRequestID("1001", "provider-789")

	// Assert that error occurred due to no rows found
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "no service request found for request ID: 1001 and provider ID: provider-789")

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetServiceProviderByRequestID_ParseRequestedTimeError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	// Mock the row with invalid time format
	rows := sqlmock.NewRows([]string{
		"id", "householder_id", "householder_name", "householder_address", "service_id",
		"requested_time", "scheduled_time", "status", "approve_status",
		"service_provider_id", "name", "contact", "address", "price", "rating", "approve",
	}).AddRow("1001", "householder-123", "John Doe", "123 Main St", "service-456",
		[]byte("invalid-time"), []byte("2024-09-02 14:00:00"),
		"Pending", true, "provider-789", "Provider Name", "1234567890",
		"456 Provider St", "100.00", 4.5, true)

	mock.ExpectQuery("SELECT service_requests.id, service_requests.householder_id").
		WithArgs("provider-789", "1001").
		WillReturnRows(rows)

	result, err := repo.GetServiceProviderByRequestID("1001", "provider-789")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error parsing requested_time")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetServiceProviderByRequestID_ParseScheduledTimeError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	// Mock the row with invalid time format
	rows := sqlmock.NewRows([]string{
		"id", "householder_id", "householder_name", "householder_address", "service_id",
		"requested_time", "scheduled_time", "status", "approve_status",
		"service_provider_id", "name", "contact", "address", "price", "rating", "approve",
	}).AddRow("1001", "householder-123", "John Doe", "123 Main St", "service-456",
		[]byte("2024-09-01 12:00:00"), []byte("invalid-time"),
		"Pending", true, "provider-789", "Provider Name", "1234567890",
		"456 Provider St", "100.00", 4.5, true)

	mock.ExpectQuery("SELECT service_requests.id, service_requests.householder_id").
		WithArgs("provider-789", "1001").
		WillReturnRows(rows)

	result, err := repo.GetServiceProviderByRequestID("1001", "provider-789")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error parsing scheduled_time")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetServiceProviderByRequestID_NoRowsReturned(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	// Mock no rows returned
	mock.ExpectQuery("SELECT service_requests.id, service_requests.householder_id").
		WithArgs("provider-789", "1001").
		WillReturnRows(sqlmock.NewRows([]string{}))

	result, err := repo.GetServiceProviderByRequestID("1001", "provider-789")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no service request found for request ID")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetServiceProviderByRequestID_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	// Mock query error
	mock.ExpectQuery("SELECT service_requests.id, service_requests.householder_id").
		WithArgs("provider-789", "1001").
		WillReturnError(fmt.Errorf("query error"))

	result, err := repo.GetServiceProviderByRequestID("1001", "provider-789")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "query error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetServiceRequestsByProviderID_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	// Mock a query error
	mock.ExpectQuery("SELECT service_requests.id").
		WithArgs("3001").
		WillReturnError(fmt.Errorf("query error"))

	// Call the method
	requests, err := repo.GetServiceRequestsByProviderID("3001")

	// Assert error and check no results
	assert.Error(t, err)
	assert.Nil(t, requests)
	assert.Contains(t, err.Error(), "query error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetServiceRequestsByProviderID_ParseRequestedTimeError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	// Mock row with invalid time format
	rows := sqlmock.NewRows([]string{
		"id", "householder_id", "householder_name", "householder_address", "service_id",
		"requested_time", "scheduled_time", "status", "approve_status",
		"service_provider_id", "name", "contact", "address", "price", "rating", "approve",
	}).AddRow(1, "1001", "John Doe", "123 Main St", "2001",
		[]byte("invalid-time"), []byte("2024-09-02 14:00:00"),
		"Pending", true, "3001", "Provider 1", "1234567890",
		"456 Provider St", "100.00", 4.5, true)

	mock.ExpectQuery("SELECT service_requests.id").
		WithArgs("3001").
		WillReturnRows(rows)

	// Call the method
	requests, err := repo.GetServiceRequestsByProviderID("3001")

	// Assert error and check no results
	assert.Error(t, err)
	assert.Nil(t, requests)
	assert.Contains(t, err.Error(), "error parsing requested_time")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetServiceRequestsByProviderID_ParseScheduledTimeError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	// Mock row with invalid time format
	rows := sqlmock.NewRows([]string{
		"id", "householder_id", "householder_name", "householder_address", "service_id",
		"requested_time", "scheduled_time", "status", "approve_status",
		"service_provider_id", "name", "contact", "address", "price", "rating", "approve",
	}).AddRow(1, "1001", "John Doe", "123 Main St", "2001",
		[]byte("2024-09-01 12:00:00"), []byte("invalid-time"),
		"Pending", true, "3001", "Provider 1", "1234567890",
		"456 Provider St", "100.00", 4.5, true)

	mock.ExpectQuery("SELECT service_requests.id").
		WithArgs("3001").
		WillReturnRows(rows)

	// Call the method
	requests, err := repo.GetServiceRequestsByProviderID("3001")

	// Assert error and check no results
	assert.Error(t, err)
	assert.Nil(t, requests)
	assert.Contains(t, err.Error(), "error parsing scheduled_time")
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetServiceRequestsByProviderID_NoRowsReturned(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %s", err)
	}
	defer db.Close()

	repo := repository.NewServiceRequestRepository(db)

	// Mock no rows returned
	mock.ExpectQuery("SELECT service_requests.id").
		WithArgs("3001").
		WillReturnRows(sqlmock.NewRows([]string{}))

	// Call the method
	requests, err := repo.GetServiceRequestsByProviderID("3001")

	// Assert no error and check empty result
	assert.NoError(t, err)
	assert.Len(t, requests, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}
