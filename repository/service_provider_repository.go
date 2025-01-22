package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"serviceNest/config"
	"serviceNest/errs"
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
	column := []string{"user_id", "rating", "availability", "is_active"}
	query := config.InsertQuery("service_providers", column)

	_, err := repo.Collection.Exec(query, provider.User.ID, provider.Rating, provider.Availability, provider.IsActive)
	return err
}

func (repo *ServiceProviderRepository) GetProviderByID(providerID string) (*model.ServiceProvider, error) {
	column := []string{"user_id", "rating", "availability", "is_active"}
	query := config.SelectQuery("service_providers", "user_id", "", column)

	row := repo.Collection.QueryRow(query, providerID)

	var provider model.ServiceProvider
	err := row.Scan(&provider.User.ID, &provider.Rating, &provider.Availability, &provider.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errs.ProviderNotFound)
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
	firstTableColumn := []string{"user_id", "rating", "availability", "is_active"}
	secondTableColumn := []string{}
	query := config.SelectInnerJoinQuery("service_providers", "service_providers_services", "service_providers.user_id = service_providers_services.service_provider_id", "service_providers_services.service_id", firstTableColumn, secondTableColumn)

	row := repo.Collection.QueryRow(query, serviceID)

	var provider model.ServiceProvider
	err := row.Scan(&provider.User.ID, &provider.Rating, &provider.Availability, &provider.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errs.ProviderNotFound)
		}
		return nil, err
	}

	return &provider, nil
}

func (repo *ServiceProviderRepository) UpdateServiceProvider(provider *model.ServiceProvider) error {
	column := []string{"rating", "availability", "is_active"}
	query := config.UpdateQuery("service_providers", "user_id", "", column)

	_, err := repo.Collection.Exec(query, provider.Rating, provider.Availability, provider.IsActive, provider.ID)
	return err

}

func (repo *ServiceProviderRepository) GetProviderDetailByID(providerID string, serviceId string) (*model.ServiceProviderDetails, error) {
	firstTableColumn := []string{"name", "address", "contact"}
	secondTableColumn := []string{"avg_rating", "rating_count"}
	query := config.SelectInnerJoinQuery("users", "services", "users.id = services.provider_id", "services.id = ? and services.provider_id", firstTableColumn, secondTableColumn)
	fmt.Println(query)
	row := repo.Collection.QueryRow(query, serviceId, providerID)
	fmt.Println(row)
	var provider model.ServiceProviderDetails
	err := row.Scan(&provider.Name, &provider.Address, &provider.Contact, &provider.Rating, &provider.RatingCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errs.ProviderNotFound)
		}
		return nil, err
	}

	return &provider, nil
}

func (repo *ServiceProviderRepository) SaveServiceProviderDetail(provider *model.ServiceProviderDetails, requestID string, serviceID string) error {
	// Check if the service provider exists in the service_providers table
	existsQuery := config.SelectCountQuery("service_providers", "user_id")
	//existsQuery := "SELECT COUNT(*) FROM service_providers WHERE user_id = ?"
	var count int
	err := repo.Collection.QueryRow(existsQuery, provider.ServiceProviderID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf(errs.ProviderNotExists)
	}

	// Proceed with the insertion
	id := util.GenerateUniqueID()
	column := []string{"id", "service_request_id", "service_provider_id", "name", "contact", "address", "price", "rating", "approve", "service_id"}
	query := config.InsertQuery("service_provider_details", column)

	_, err = repo.Collection.Exec(query, id, requestID, provider.ServiceProviderID, provider.Name, provider.Contact, provider.Address, provider.Price, provider.Rating, provider.Approve, serviceID)
	return err
}

func (repo *ServiceProviderRepository) UpdateServiceProviderDetailByRequestID(provider *model.ServiceProviderDetails, requestID string) error {
	column := []string{"approve"}
	query := config.UpdateQuery("service_provider_details", "service_provider_id", "service_request_id", column)

	_, err := repo.Collection.Exec(query, provider.Approve, provider.ServiceProviderID, requestID)
	return err
}

func (repo *ServiceProviderRepository) IsProviderApproved(providerID string) (bool, error) {
	var approveStatus bool
	column := []string{"approve"}
	query := config.SelectQuery("service_provider_details", "service_provider_id", "approve", column)

	err := repo.Collection.QueryRow(query, providerID).Scan(&approveStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New(errs.ProviderNotFound)
		}
		return false, err
	}
	return approveStatus, nil
}

// AddReview adds a review to the reviews table
func (repo *ServiceProviderRepository) AddReview(review model.Review) error {
	tx, err := repo.Collection.Begin()
	if err != nil {
		return err
	}

	columns := []string{"id", "provider_id", "service_id", "householder_id", "rating", "comments", "review_date"}
	reviewQuery := config.InsertQuery("reviews", columns)

	checkQuery := config.CountReviewAddedQuery()
	var count int
	err = tx.QueryRow(checkQuery, review.ProviderID, review.ServiceID, review.HouseholderID).Scan(&count)
	if err != nil {
		tx.Rollback()
		return err
	}

	if count > 0 {

		tx.Rollback()
		return fmt.Errorf("review already exists for this provider, service, and householder")
	}

	_, err = tx.Exec(reviewQuery, review.ID, review.ProviderID, review.ServiceID, review.HouseholderID, review.Rating, review.Comments, review.ReviewDate)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// UpdateProviderRating recalculates and updates the provider's average rating
func (repo *ServiceProviderRepository) UpdateProviderRating(providerID string, serviceId string, rating float64) error {
	// Calculate the  rating from the service table
	column := []string{"avg_rating", "rating_count"}
	ratingQuery := config.SelectQuery("services", "id", "provider_id", column)

	var avgRating float64
	var ratingCount int64
	err := repo.Collection.QueryRow(ratingQuery, serviceId, providerID).Scan(&avgRating, &ratingCount)
	if err != nil {
		return fmt.Errorf("%v: %v", errs.FailCalculateRating, err)
	}
	updatedRating := util.CalculateRating(avgRating, ratingCount, rating)
	// Update the rating in the service_providers table
	providerColumn := []string{"rating"}
	updateServiceProviderQuery := config.UpdateQuery("service_providers", "user_id", "", providerColumn)

	_, err = repo.Collection.Exec(updateServiceProviderQuery, updatedRating, providerID)
	if err != nil {
		return fmt.Errorf("%ve: %v", errs.FailUpdateRating, err)
	}

	// Update the rating in the service_provider_details table
	providerDetailColumn := []string{"rating"}
	updateServiceProviderDetailsQuery := config.UpdateQuery("service_provider_details", "service_provider_id", "service_id", providerDetailColumn)

	_, err = repo.Collection.Exec(updateServiceProviderDetailsQuery, updatedRating, providerID, serviceId)
	if err != nil {
		return fmt.Errorf("%v %v", errs.FailUpdateRating, err)
	}
	// Update the rating in the services table
	serviceDetailColumn := []string{"avg_rating", "rating_count"}
	updateServiceDetailsQuery := config.UpdateQuery("services", "id", "provider_id", serviceDetailColumn)

	_, err = repo.Collection.Exec(updateServiceDetailsQuery, updatedRating, ratingCount+1, serviceId, providerID)
	if err != nil {
		return fmt.Errorf("%v %v", errs.FailUpdateRating, err)
	}
	return nil
}

func (repo *ServiceProviderRepository) GetReviewsByProviderID(providerID string, limit, offset int, serviceID string) ([]model.Review, error) {
	column := []string{"id", "provider_id", "service_id", "householder_id", "rating", "comments", "review_date"}
	var query string
	var rows *sql.Rows
	var err error

	// Determine query and arguments based on provided filters
	if providerID != "" && serviceID != "" {
		// Both filters provided
		query = config.SelectQueryWithLimit("reviews", "provider_id", "service_id", column, limit, offset)
		rows, err = repo.Collection.Query(query, providerID, serviceID)
	} else if providerID != "" {
		// Only provider filter provided
		query = config.SelectQueryWithLimit("reviews", "provider_id", "", column, limit, offset)
		rows, err = repo.Collection.Query(query, providerID)
	}
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
	if len(reviews) == 0 {
		return nil, fmt.Errorf("no reviews found")
	}
	return reviews, nil
}
func (repo *ServiceProviderRepository) AddServiceToProvider(providerID, serviceID string) error {
	var existingProviderID string
	coloumn1 := []string{"user_id"}
	checkProviderQuery := config.SelectQuery("service_providers", "user_id", "", coloumn1)
	err := repo.Collection.QueryRow(checkProviderQuery, providerID).Scan(&existingProviderID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If the provider does not exist, insert a new provider
			coloumn2 := []string{"user_id", "rating", "availability", "is_active"}
			insertProviderQuery := config.InsertQuery("service_providers", coloumn2)

			_, err := repo.Collection.Exec(insertProviderQuery, providerID, 0, 1, 1)
			if err != nil {
				return fmt.Errorf("failed to insert new provider: %v", err)
			}
		} else {
			// Handle other potential errors in checking the provider
			return fmt.Errorf("failed to check provider existence: %v", err)
		}
	}

	columns := []string{"service_provider_id", "service_id", "avg_rating", "rating_count"}
	query := config.InsertQuery("service_providers_services", columns)

	// Initializing avg_rating and rating_count for a new service
	_, err = repo.Collection.Exec(query, providerID, serviceID, 0.0, 0)
	if err != nil {
		return fmt.Errorf("failed to add service to provider: %v", err)
	}
	return nil
}

func (repo *ServiceProviderRepository) DeleteServicesByProviderID(providerID string) error {
	query := config.DeleteQuery("services", "provider_id", "")
	_, err := repo.Collection.Exec(query, providerID)
	return err
}
