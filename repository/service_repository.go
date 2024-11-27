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

type ServiceRepository struct {
	db *sql.DB
}

// NewServiceRepository creates a new instance of ServiceRepository for MySQL
func NewServiceRepository(client *sql.DB) interfaces.ServiceRepository {
	return &ServiceRepository{db: client}
}

func (repo *ServiceRepository) GetAllServices(limit, offset int) ([]model.Service, error) {
	column := []string{"id", "name", "description", "price", "provider_id", "category", "avg_rating", "rating_count"}
	query := config.SelectQueryWithLimit("services", "", "", column, limit, offset)

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var service model.Service
		var providerID sql.NullString

		if err := rows.Scan(&service.ID, &service.Name, &service.Description, &service.Price, &providerID, &service.Category, &service.AvgRating, &service.RatingCount); err != nil {
			return nil, err
		}

		if providerID.Valid {
			service.ProviderID = providerID.String
		} else {
			service.ProviderID = ""
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// GetServiceByID retrieves a service by its ID
func (repo *ServiceRepository) GetServiceByID(serviceID string) (*model.Service, error) {
	column := []string{"id", "name", "description", "price", "provider_id", "category"}
	query := config.SelectQuery("services", "id", "", column)

	var service model.Service
	err := repo.db.QueryRow(query, serviceID).Scan(&service.ID, &service.Name, &service.Description, &service.Price, &service.ProviderID, &service.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errs.ServiceNotFound)
		}
		return nil, err
	}

	return &service, nil
}

// SaveService adds a new service to the MySQL database
func (repo *ServiceRepository) SaveService(service model.Service) error {
	column := []string{"id", "name", "description", "price", "provider_id", "category", "avg_rating", "rating_count"}
	query := config.InsertQuery("services", column)

	var providerID *string
	if service.ProviderID == "" {
		providerID = nil
	} else {
		providerID = &service.ProviderID
	}
	_, err := repo.db.Exec(query, service.ID, service.Name, service.Description, service.Price, providerID, service.Category, service.AvgRating, service.RatingCount)
	return err
}

// RemoveService removes a service from the MySQL database
func (repo *ServiceRepository) RemoveService(serviceID string) error {
	// Begin a transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Ensure the transaction is rolled back in case of an error
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Delete from service_provider_details
	query := config.DeleteQuery("service_provider_details", "service_id", "")
	_, err = tx.Exec(query, serviceID)
	if err != nil {
		return fmt.Errorf("failed to delete from service_provider_details: %v", err)
	}

	// Delete from service_requests
	query = config.DeleteQuery("service_requests", "service_id", "")
	_, err = tx.Exec(query, serviceID)
	if err != nil {
		return fmt.Errorf("failed to delete from service_requests: %v", err)
	}

	// Delete from service_providers_services
	query = config.DeleteQuery("service_providers_services", "service_id", "")
	_, err = tx.Exec(query, serviceID)
	if err != nil {
		return fmt.Errorf("failed to delete from service_providers_services: %v", err)
	}

	// Delete from service_categories (if applicable)
	query = config.DeleteQuery("service_categories", "id", "")
	_, err = tx.Exec(query, serviceID)
	if err != nil {
		return fmt.Errorf("failed to delete from service_categories: %v", err)
	}

	// Finally, delete from services table
	query = config.DeleteQuery("services", "id", "")
	_, err = tx.Exec(query, serviceID)
	if err != nil {
		return fmt.Errorf("failed to delete from services: %v", err)
	}

	return nil
}
func (repo *ServiceRepository) GetServiceIdByCategory(category string) (*string, error) {
	column := []string{"id"}
	query := config.SelectQuery("service_categories", "category_name", "", column)

	var categoryId string
	err := repo.db.QueryRow(query, category).Scan(&categoryId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(errs.ServiceNotFound)
		}
		return nil, err
	}
	return &categoryId, nil
}

// GetServiceByProviderID retrieves a service by its ProviderID
func (repo *ServiceRepository) GetServiceByProviderID(providerID string) ([]model.Service, error) {
	column := []string{"id", "name", "description", "price", "provider_id", "category", "avg_rating", "rating_count"}
	query := config.SelectQuery("services", "provider_id", "", column)

	rows, err := repo.db.Query(query, providerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errs.ServiceNotFound)
		}
		return nil, err
	}
	var services []model.Service
	for rows.Next() {
		var service model.Service
		err := rows.Scan(&service.ID, &service.Name, &service.Description, &service.Price, &service.ProviderID, &service.Category, &service.AvgRating, &service.RatingCount)
		if err != nil {
			return nil, err
		} else {
			services = append(services, service)
		}
	}

	return services, nil
}

func (repo *ServiceRepository) UpdateService(providerID string, updatedService model.Service) error {
	column := []string{"name", "description", "price"}
	query := config.UpdateQuery("services", "provider_id", "id", column)

	result, err := repo.db.Exec(query, updatedService.Name, updatedService.Description, updatedService.Price, providerID, updatedService.ID)
	// Check how many rows were affected
	if err != nil {

		return err
	}
	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.New(errs.ServiceIdNotExists)
	}
	return nil
}

func (repo *ServiceRepository) RemoveServiceByProviderID(providerID string, serviceID string) error {
	// Begin a transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := config.DeleteQuery("services", "id", "provider_id")
	result, err := tx.Exec(query, serviceID, providerID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New(errs.InvalidServiceId)
	}

	query2 := config.DeleteQuery("service_providers_services", "service_id", "service_provider_id")
	_, err = tx.Exec(query2, serviceID, providerID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ServiceRepository) CategoryExists(categoryName string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM service_categories WHERE category_name = ?)", categoryName).Scan(&exists)
	return exists, err
}

func (repo *ServiceRepository) GetServicesByCategory(category string) ([]model.Service, error) {
	column := []string{"id", "name", "description", "price", "provider_id", "category", "avg_rating"}
	query := config.SelectQuery("services", "category", "", column)

	rows, err := repo.db.Query(query, category)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errs.ServiceNotFound)
		}
		return nil, err
	}
	var services []model.Service
	for rows.Next() {
		var service model.Service
		err := rows.Scan(&service.ID, &service.Name, &service.Description, &service.Price, &service.ProviderID, &service.Category, &service.AvgRating)
		if err != nil {
			return nil, err
		} else {
			services = append(services, service)
		}
	}

	return services, nil
}

func (repo *ServiceRepository) GetAllCategory() ([]model.Category, error) {
	column := []string{"id", "category_name", "description"}
	query := config.SelectQuery("service_categories", "", "", column)

	rows, err := repo.db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errs.ServiceNotFound)
		}
		return nil, err
	}
	var categories []model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (repo *ServiceRepository) AddCategory(service *model.Category) error {
	column := []string{"id", "category_name", "description"}
	query := config.InsertQuery("service_categories", column)

	_, err := repo.db.Exec(query, service.ID, service.Name, service.Description)
	return err

}
