package repository_test

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"serviceNest/model"
	"serviceNest/repository"
	"serviceNest/tests/mocks"
	"testing"
	"time"
)

func TestServiceProviderRepository_AddReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	review := model.Review{
		ProviderID:    "provider123",
		HouseholderID: "householder123",
		Comments:      "Great service!",
		Rating:        4.5,
	}

	// Expectations
	mockRepo.EXPECT().
		AddReview(review).
		Return(nil)

	// Execute
	err := mockRepo.AddReview(review)

	// Verify
	assert.NoError(t, err)
}

func TestServiceProviderRepository_GetProviderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	providerID := "provider123"
	expectedProvider := &model.ServiceProvider{
		User: model.User{
			ID:   providerID,
			Name: "Test Provider",
		},
	}

	// Expectations
	mockRepo.EXPECT().
		GetProviderByID(providerID).
		Return(expectedProvider, nil)

	// Execute
	result, err := mockRepo.GetProviderByID(providerID)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedProvider, result)
}

func TestServiceProviderRepository_GetProviderByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	providerID := "invalid-provider-id"
	expectedError := errors.New("provider not found")

	// Expectations
	mockRepo.EXPECT().
		GetProviderByID(providerID).
		Return(nil, expectedError)

	// Execute
	result, err := mockRepo.GetProviderByID(providerID)

	// Verify
	assert.Nil(t, result)
	assert.EqualError(t, err, expectedError.Error())
}

func TestServiceProviderRepository_GetProvidersByServiceType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	serviceType := "Plumbing"
	expectedProviders := []model.ServiceProvider{
		{
			User: model.User{
				ID:   "provider123",
				Name: "Provider 1",
			},
		},
		{
			User: model.User{
				ID:   "provider456",
				Name: "Provider 2",
			},
		},
	}

	// Expectations
	mockRepo.EXPECT().
		GetProvidersByServiceType(serviceType).
		Return(expectedProviders, nil)

	// Execute
	result, err := mockRepo.GetProvidersByServiceType(serviceType)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedProviders, result)
}

func TestServiceProviderRepository_SaveServiceProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	provider := model.ServiceProvider{
		User: model.User{
			ID:   "provider123",
			Name: "Test Provider",
		},
	}

	// Expectations
	mockRepo.EXPECT().
		SaveServiceProvider(provider).
		Return(nil)

	// Execute
	err := mockRepo.SaveServiceProvider(provider)

	// Verify
	assert.NoError(t, err)
}

func TestServiceProviderRepository_UpdateServiceProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	provider := &model.ServiceProvider{
		User: model.User{
			ID:   "provider123",
			Name: "Updated Provider",
		},
	}

	// Expectations
	mockRepo.EXPECT().
		UpdateServiceProvider(provider).
		Return(nil)

	// Execute
	err := mockRepo.UpdateServiceProvider(provider)

	// Verify
	assert.NoError(t, err)
}
func TestGetProviderDetailByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	providerID := "provider123"
	expectedProvider := &model.ServiceProviderDetails{
		Name:    "John Doe",
		Address: "123 Main St",
		Contact: "1234567890",
		Rating:  4.5,
	}

	rows := sqlmock.NewRows([]string{"name", "address", "contact", "rating"}).
		AddRow(expectedProvider.Name, expectedProvider.Address, expectedProvider.Contact, expectedProvider.Rating)

	// Correct the query expectation by matching it exactly with the repository's query
	mock.ExpectQuery("SELECT users.name, users.address, users.contact, service_providers.rating FROM users INNER JOIN service_providers ON users.id = service_providers.user_id WHERE users.id = ?").
		WithArgs(providerID).
		WillReturnRows(rows)

	provider, err := repo.GetProviderDetailByID(providerID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProvider, provider)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveServiceProviderDetail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	provider := &model.ServiceProviderDetails{
		ServiceProviderID: "provider123",
		Name:              "John Doe",
		Contact:           "1234567890",
		Address:           "123 Main St",
		Price:             "100",
		Rating:            4.5,
		Approve:           true,
	}
	requestID := "request123"

	// Mock the check for service provider existence
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM service_providers WHERE user_id = ?").
		WithArgs(provider.ServiceProviderID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock the insert into service_provider_details
	mock.ExpectExec("INSERT INTO service_provider_details").
		WithArgs(sqlmock.AnyArg(), requestID, provider.ServiceProviderID, provider.Name, provider.Contact, provider.Address, provider.Price, provider.Rating, provider.Approve).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveServiceProviderDetail(provider, requestID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestUpdateServiceProviderDetailByRequestID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	provider := &model.ServiceProviderDetails{
		ServiceProviderID: "provider123",
		Approve:           true,
	}
	requestID := "request123"

	// Mock the update
	mock.ExpectExec("UPDATE service_provider_details").
		WithArgs(provider.Approve, provider.ServiceProviderID, requestID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateServiceProviderDetailByRequestID(provider, requestID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestIsProviderApproved(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	providerID := "provider123"

	// Mock the query that checks if provider is approved
	mock.ExpectQuery("SELECT approve FROM service_provider_details").
		WithArgs(providerID).
		WillReturnRows(sqlmock.NewRows([]string{"approve"}).AddRow(true))

	isApproved, err := repo.IsProviderApproved(providerID)

	assert.NoError(t, err)
	assert.True(t, isApproved)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestAddReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	review := model.Review{
		ID:            "review123",
		ProviderID:    "provider123",
		ServiceID:     "service123",
		HouseholderID: "householder123",
		Rating:        4,
		Comments:      "Great service",
		ReviewDate:    time.Now(),
	}

	// Mock the transaction and insert review
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO reviews").
		WithArgs(review.ID, review.ProviderID, review.ServiceID, review.HouseholderID, review.Rating, review.Comments, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.AddReview(review)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

//func TestUpdateProviderRating(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	assert.NoError(t, err)
//	defer db.Close()
//
//	repo := repository.NewServiceProviderRepository(db)
//
//	providerID := "provider123"
//	avgRating := 4.5
//
//	// Mock the query that calculates the average rating
//	mock.ExpectQuery("SELECT AVG(rating) FROM reviews WHERE provider_id = ?").
//		WithArgs(providerID).
//		WillReturnRows(sqlmock.NewRows([]string{"AVG(rating)"}).AddRow(avgRating))
//
//	// Mock the update for service_providers
//	mock.ExpectExec("UPDATE service_providers").
//		WithArgs(avgRating, providerID).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//
//	// Mock the update for service_provider_details
//	mock.ExpectExec("UPDATE service_provider_details").
//		WithArgs(avgRating, providerID).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//
//	err = repo.UpdateProviderRating(providerID)
//
//	assert.NoError(t, err)
//	assert.NoError(t, mock.ExpectationsWereMet())
//}

func TestUpdateProviderRating(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	providerID := "provider123"
	avgRating := 4.5

	// Mock the query that calculates the average rating
	mock.ExpectQuery("(?i)SELECT\\s+AVG\\(rating\\)\\s+FROM\\s+reviews\\s+WHERE\\s+provider_id\\s*=\\s*\\?").
		WithArgs(providerID).
		WillReturnRows(sqlmock.NewRows([]string{"AVG(rating)"}).AddRow(avgRating))

	// Mock the update for service_providers
	mock.ExpectExec("UPDATE service_providers").
		WithArgs(avgRating, providerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock the update for service_provider_details
	mock.ExpectExec("UPDATE service_provider_details").
		WithArgs(avgRating, providerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateProviderRating(providerID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetReviewsByProviderID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	providerID := "provider123"
	reviewDate := time.Now().Format(time.RFC3339) // RFC3339 format

	// Mock the query that fetches reviews
	rows := sqlmock.NewRows([]string{"id", "provider_id", "service_id", "householder_id", "rating", "comments", "review_date"}).
		AddRow("review123", providerID, "service123", "householder123", 4, "Great service", reviewDate)

	mock.ExpectQuery("SELECT id, provider_id, service_id, householder_id, rating, comments, review_date FROM reviews").
		WithArgs(providerID).
		WillReturnRows(rows)

	reviews, err := repo.GetReviewsByProviderID(providerID)

	assert.NoError(t, err)
	assert.Len(t, reviews, 1)
	assert.Equal(t, "review123", reviews[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestSaveServiceProvider_Success(t *testing.T) {
	// Mock the database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	// Mock provider data
	provider := model.ServiceProvider{
		User:         model.User{ID: "user123"},
		Rating:       4.5,
		Availability: true,
		IsActive:     true,
	}

	// Expect the SQL statement to be executed
	query := regexp.QuoteMeta("INSERT INTO service_providers (user_id, rating, availability, is_active) VALUES (?, ?, ?, ?)")
	mock.ExpectExec(query).WithArgs(provider.User.ID, provider.Rating, provider.Availability, provider.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function
	err = repo.SaveServiceProvider(provider)
	assert.NoError(t, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetProviderByID_Success(t *testing.T) {
	// Mock the database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	// Mock provider data
	expectedProvider := model.ServiceProvider{
		User:         model.User{ID: "user123"},
		Rating:       4.5,
		Availability: true,
		IsActive:     true,
	}

	// Expect the query to be executed and return the result
	query := regexp.QuoteMeta("SELECT user_id, rating, availability, is_active FROM service_providers WHERE user_id = ?")
	rows := sqlmock.NewRows([]string{"user_id", "rating", "availability", "is_active"}).
		AddRow(expectedProvider.User.ID, expectedProvider.Rating, expectedProvider.Availability, expectedProvider.IsActive)
	mock.ExpectQuery(query).WithArgs("user123").WillReturnRows(rows)

	// Call the function
	provider, err := repo.GetProviderByID("user123")
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, expectedProvider, *provider)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetProviderByID_NotFound(t *testing.T) {
	// Mock the database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	// Expect the query to return no rows
	query := regexp.QuoteMeta("SELECT user_id, rating, availability, is_active FROM service_providers WHERE user_id = ?")
	mock.ExpectQuery(query).WithArgs("user123").WillReturnError(sql.ErrNoRows)

	// Call the function
	provider, err := repo.GetProviderByID("user123")

	// Assert the provider was not found
	assert.Nil(t, provider)
	assert.EqualError(t, err, "provider not found")

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetProvidersByServiceType_Success(t *testing.T) {
	// Mock the database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	// Mock provider data
	expectedProviders := []model.ServiceProvider{
		{User: model.User{ID: "provider1"}, Rating: 4.5, Availability: true, IsActive: true},
		{User: model.User{ID: "provider2"}, Rating: 4.0, Availability: false, IsActive: true},
	}

	// Expect the query and return rows
	query := regexp.QuoteMeta(`
		SELECT sp.user_id, sp.rating, sp.availability, sp.is_active
		FROM service_providers sp
		INNER JOIN service_providers_services sps ON sp.user_id = sps.service_provider_id
		INNER JOIN services s ON sps.service_id = s.id
		WHERE s.name = ?`)
	rows := sqlmock.NewRows([]string{"user_id", "rating", "availability", "is_active"}).
		AddRow("provider1", 4.5, true, true).
		AddRow("provider2", 4.0, false, true)
	mock.ExpectQuery(query).WithArgs("Plumbing").WillReturnRows(rows)

	// Call the function
	providers, err := repo.GetProvidersByServiceType("Plumbing")
	assert.NoError(t, err)
	assert.Len(t, providers, 2)
	assert.Equal(t, expectedProviders, providers)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateServiceProvider_Success(t *testing.T) {
	// Mock the database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	// Mock provider data
	provider := &model.ServiceProvider{
		User:         model.User{ID: "provider123"},
		Rating:       4.5,
		Availability: true,
		IsActive:     true,
	}

	// Expect the update query
	query := regexp.QuoteMeta(`
		UPDATE service_providers
		SET rating = ?, availability = ?, is_active = ?
		WHERE user_id = ?`)
	mock.ExpectExec(query).WithArgs(provider.Rating, provider.Availability, provider.IsActive, provider.User.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the function
	err = repo.UpdateServiceProvider(provider)
	assert.NoError(t, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetProviderByServiceID_Success(t *testing.T) {
	// Create a new SQL mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock database
	repo := repository.NewServiceProviderRepository(db)

	// Define expected rows and mock the query
	query := `
	SELECT service_providers.user_id, service_providers.rating, service_providers.availability, service_providers.is_active
	FROM service_providers 
	INNER JOIN service_providers_services  ON service_providers.user_id = service_providers_services.service_provider_id
	WHERE service_providers_services.service_id = ?
	`
	rows := sqlmock.NewRows([]string{"user_id", "rating", "availability", "is_active"}).
		AddRow("provider1", 4.5, true, true)
	mock.ExpectQuery(query).WithArgs("service1").WillReturnRows(rows)

	// Call the method
	provider, err := repo.GetProviderByServiceID("service1")

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, "provider1", provider.User.ID)
	assert.Equal(t, 4.5, provider.Rating)
	assert.Equal(t, true, provider.Availability)
	assert.Equal(t, true, provider.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProviderByServiceID_ProviderNotFound(t *testing.T) {
	// Create a new SQL mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock database
	repo := repository.NewServiceProviderRepository(db)

	// Define the mock query to return no rows
	query := `
	SELECT service_providers.user_id, service_providers.rating, service_providers.availability, service_providers.is_active
	FROM service_providers 
	INNER JOIN service_providers_services  ON service_providers.user_id = service_providers_services.service_provider_id
	WHERE service_providers_services.service_id = ?
	`
	mock.ExpectQuery(query).WithArgs("service1").WillReturnRows(sqlmock.NewRows([]string{}))

	// Call the method
	provider, err := repo.GetProviderByServiceID("service1")

	// Assertions
	assert.Nil(t, provider)
	assert.EqualError(t, err, "provider not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProviderByServiceID_QueryError(t *testing.T) {
	// Create a new SQL mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock database
	repo := repository.NewServiceProviderRepository(db)

	// Define the mock query to return an error
	query := `
	SELECT service_providers.user_id, service_providers.rating, service_providers.availability, service_providers.is_active
	FROM service_providers 
	INNER JOIN service_providers_services  ON service_providers.user_id = service_providers_services.service_provider_id
	WHERE service_providers_services.service_id = ?
	`
	mock.ExpectQuery(query).WithArgs("service1").WillReturnError(errors.New("query error"))

	// Call the method
	provider, err := repo.GetProviderByServiceID("service1")

	// Assertions
	assert.Nil(t, provider)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProviderDetailByID_ProviderNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	providerID := "provider123"
	query := "SELECT users.name, users.address, users.contact, service_providers.rating FROM users INNER JOIN service_providers ON users.id = service_providers.user_id WHERE users.id = ?"
	mock.ExpectQuery(query).
		WithArgs(providerID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	provider, err := repo.GetProviderDetailByID(providerID)
	assert.Nil(t, provider)
	assert.EqualError(t, err, "provider not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProviderDetailByID_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	providerID := "provider123"
	query := "SELECT users.name, users.address, users.contact, service_providers.rating FROM users INNER JOIN service_providers ON users.id = service_providers.user_id WHERE users.id = ?"
	mock.ExpectQuery(query).
		WithArgs(providerID).
		WillReturnError(errors.New("query error"))

	provider, err := repo.GetProviderDetailByID(providerID)

	assert.Nil(t, provider)
	assert.EqualError(t, err, "query error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestAddReviewAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewServiceProviderRepository(db)

	review := model.Review{
		ID:            "review123",
		ProviderID:    "provider123",
		ServiceID:     "service123",
		HouseholderID: "householder123",
		Rating:        4,
		Comments:      "Great service",
		ReviewDate:    time.Now(),
	}

	t.Run("Successful Review Addition", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO reviews").
			WithArgs(review.ID, review.ProviderID, review.ServiceID, review.HouseholderID, review.Rating, review.Comments, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.AddReview(review)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Transaction Begin Error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errors.New("begin error"))

		err = repo.AddReview(review)

		assert.Error(t, err)
		assert.Equal(t, "begin error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SQL Execution Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO reviews").
			WithArgs(review.ID, review.ProviderID, review.ServiceID, review.HouseholderID, review.Rating, review.Comments, sqlmock.AnyArg()).
			WillReturnError(errors.New("execution error"))
		mock.ExpectRollback()

		err = repo.AddReview(review)

		assert.Error(t, err)
		assert.Equal(t, "execution error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Transaction Commit Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO reviews").
			WithArgs(review.ID, review.ProviderID, review.ServiceID, review.HouseholderID, review.Rating, review.Comments, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(errors.New("commit error"))

		err = repo.AddReview(review)

		assert.Error(t, err)
		assert.Equal(t, "commit error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
