package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"serviceNest/interfaces"
	"serviceNest/model"
	"serviceNest/util"
)

type ServiceProviderRepository struct {
	Collection *sql.DB
}

// NewServiceProviderRepository initializes a new ServiceProviderRepository
func NewServiceProviderRepository(collection *sql.DB) interfaces.ServiceProviderRepository {
	return &ServiceProviderRepository{Collection: collection}
}

func (repo *ServiceProviderRepository) SaveServiceProvider(provider model.ServiceProvider) error {
	query := "INSERT INTO service_providers (user_id, rating, availability, is_active) VALUES (?, ?, ?, ?)"
	_, err := repo.Collection.Exec(query, provider.User.ID, provider.Rating, provider.Availability, provider.IsActive)
	return err
}

func (repo *ServiceProviderRepository) GetProviderByID(providerID string) (*model.ServiceProvider, error) {
	query := "SELECT user_id, rating, availability, is_active FROM service_providers WHERE user_id = ?"
	row := repo.Collection.QueryRow(query, providerID)

	var provider model.ServiceProvider
	err := row.Scan(&provider.User.ID, &provider.Rating, &provider.Availability, &provider.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("provider not found")
		}
		return nil, err
	}

	return &provider, nil
}

func (repo *ServiceProviderRepository) GetProvidersByServiceType(serviceType string) ([]model.ServiceProvider, error) {
	query := `
	SELECT sp.user_id, sp.rating, sp.availability, sp.is_active 
	FROM service_providers sp
	INNER JOIN service_providers_services sps ON sp.user_id = sps.service_provider_id
	INNER JOIN services s ON sps.service_id = s.id
	WHERE s.name = ?
	`
	rows, err := repo.Collection.Query(query, serviceType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []model.ServiceProvider
	for rows.Next() {
		var provider model.ServiceProvider
		err := rows.Scan(&provider.User.ID, &provider.Rating, &provider.Availability, &provider.IsActive)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return providers, nil
}

func (repo *ServiceProviderRepository) GetProviderByServiceID(serviceID string) (*model.ServiceProvider, error) {
	query := `
	SELECT sp.user_id, sp.rating, sp.availability, sp.is_active 
	FROM service_providers sp
	INNER JOIN service_providers_services sps ON sp.user_id = sps.service_provider_id
	WHERE sps.service_id = ?
	`
	row := repo.Collection.QueryRow(query, serviceID)

	var provider model.ServiceProvider
	err := row.Scan(&provider.User.ID, &provider.Rating, &provider.Availability, &provider.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("provider not found")
		}
		return nil, err
	}

	return &provider, nil
}

func (repo *ServiceProviderRepository) UpdateServiceProvider(provider *model.ServiceProvider) error {
	query := `
	UPDATE service_providers
	SET rating = ?, availability = ?, is_active = ?
	WHERE user_id = ?
	`
	_, err := repo.Collection.Exec(query, provider.Rating, provider.Availability, provider.IsActive, provider.ID)
	return err
}

func (repo *ServiceProviderRepository) GetProviderDetailByID(providerID string) (*model.ServiceProviderDetails, error) {
	query := "SELECT name, address, contact,rating FROM users INNER JOIN service_providers ON id=user_id WHERE id = ?"
	row := repo.Collection.QueryRow(query, providerID)

	var provider model.ServiceProviderDetails
	err := row.Scan(&provider.Name, &provider.Address, &provider.Contact, &provider.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("provider not found")
		}
		return nil, err
	}

	return &provider, nil
}

func (repo *ServiceProviderRepository) SaveServiceProviderDetail(provider *model.ServiceProviderDetails, requestID string) error {
	// Check if the service provider exists in the service_providers table
	existsQuery := "SELECT COUNT(*) FROM service_providers WHERE user_id = ?"
	var count int
	err := repo.Collection.QueryRow(existsQuery, provider.ServiceProviderID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("service provider does not exist")
	}

	// Proceed with the insertion
	id := util.GenerateUniqueID()

	query := "INSERT INTO service_provider_details (id,service_request_id,service_provider_id,name,contact,address,price,rating,approve) VALUES (?, ?, ?, ?,?,?,?,?,?)"
	_, err = repo.Collection.Exec(query, id, requestID, provider.ServiceProviderID, provider.Name, provider.Contact, provider.Address, provider.Price, provider.Rating, provider.Approve)
	return err
}

func (repo *ServiceProviderRepository) UpdateServiceProviderDetailByRequestID(provider *model.ServiceProviderDetails, requestID string) error {
	query := `
	UPDATE service_provider_details
	SET approve= ?
	WHERE service_provider_id = ? and service_request_id=?
	`
	_, err := repo.Collection.Exec(query, provider.Approve, provider.ServiceProviderID, requestID)
	return err
}

func (repo *ServiceProviderRepository) IsProviderApproved(providerID string) (bool, error) {
	var approveStatus bool
	query := `
	SELECT approve FROM service_provider_details
	WHERE service_provider_id = ? AND approve = 1
	`
	err := repo.Collection.QueryRow(query, providerID).Scan(&approveStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("service provider not found")
		}
		return false, err
	}
	return approveStatus, nil
}

// AddReview adds a review to the reviews table
func (repo *ServiceProviderRepository) AddReview(review model.Review) error {
	tx, err := repo.Collection.Begin() // Start a transaction
	if err != nil {
		return err
	}

	// Insert the review into the reviews table with providerID
	reviewQuery := `
	INSERT INTO reviews (id, provider_id, service_id, householder_id, rating, comments, review_date)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(reviewQuery, review.ID, review.ProviderID, review.ServiceID, review.HouseholderID, review.Rating, review.Comments, review.ReviewDate)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit() // Commit the transaction
}

// UpdateProviderRating recalculates and updates the provider's average rating
func (repo *ServiceProviderRepository) UpdateProviderRating(providerID string) error {
	// Calculate the average rating from the reviews table
	ratingQuery := `
	SELECT AVG(r.rating)
	FROM reviews r
	WHERE r.provider_id = ?
	`
	var avgRating float64
	err := repo.Collection.QueryRow(ratingQuery, providerID).Scan(&avgRating)
	if err != nil {
		return fmt.Errorf("failed to calculate average rating: %v", err)
	}

	// Update the rating in the service_providers table
	updateServiceProviderQuery := `
	UPDATE service_providers
	SET rating = ?
	WHERE user_id = ?
	`
	_, err = repo.Collection.Exec(updateServiceProviderQuery, avgRating, providerID)
	if err != nil {
		return fmt.Errorf("failed to update rating in service_providers table: %v", err)
	}

	// Update the rating in the service_provider_details table
	updateServiceProviderDetailsQuery := `
	UPDATE service_provider_details
	SET rating = ?
	WHERE service_provider_id = ?
	`
	_, err = repo.Collection.Exec(updateServiceProviderDetailsQuery, avgRating, providerID)
	if err != nil {
		return fmt.Errorf("failed to update rating in service_provider_details table: %v", err)
	}

	return nil
}

func (repo *ServiceProviderRepository) GetReviewsByProviderID(providerID string) ([]model.Review, error) {
	query := `
	SELECT id, provider_id, service_id, householder_id, rating, comments, review_date
	FROM reviews
	WHERE provider_id = ?
	`
	rows, err := repo.Collection.Query(query, providerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var review model.Review
		var reviewDate []uint8
		err := rows.Scan(&review.ID, &review.ProviderID, &review.ServiceID, &review.HouseholderID, &review.Rating, &review.Comments, &reviewDate)
		if err != nil {
			return nil, err
		}
		parsedDate, err := util.ParseTime(reviewDate)
		if err != nil {
			return nil, err
		}
		review.ReviewDate = parsedDate
		reviews = append(reviews, review)
	}

	return reviews, nil
}
