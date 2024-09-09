package repository

import (
	"database/sql"
	"errors"
	"log"
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

func (repo *ServiceRepository) GetAllServices() ([]model.Service, error) {
	query := "SELECT id, name, description, price, provider_id, category FROM services"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var service model.Service
		var providerID sql.NullString

		if err := rows.Scan(&service.ID, &service.Name, &service.Description, &service.Price, &providerID, &service.Category); err != nil {
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
	query := "SELECT id, name, description, price, provider_id, category FROM services WHERE id = ?"
	var service model.Service
	err := repo.db.QueryRow(query, serviceID).Scan(&service.ID, &service.Name, &service.Description, &service.Price, &service.ProviderID, &service.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("service not found")
		}
		return nil, err
	}

	return &service, nil
}

// SaveService adds a new service to the MySQL database
func (repo *ServiceRepository) SaveService(service model.Service) error {
	query := "INSERT INTO services (id, name, description, price, provider_id, category) VALUES (?, ?, ?, ?, ?, ?)"
	var providerID *string
	if service.ProviderID == "" {
		providerID = nil
	} else {
		providerID = &service.ProviderID
	}
	_, err := repo.db.Exec(query, service.ID, service.Name, service.Description, service.Price, providerID, service.Category)
	return err
}

// SaveAllServices saves the entire list of services to the MySQL database
func (repo *ServiceRepository) SaveAllServices(services []model.Service) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO services (id, name, description, price, provider_id, category) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description), price=VALUES(price), provider_id=VALUES(provider_id), category=VALUES(category)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, service := range services {
		if _, err := stmt.Exec(service.ID, service.Name, service.Description, service.Price, service.ProviderID, service.Category); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// RemoveService removes a service from the MySQL database
func (repo *ServiceRepository) RemoveService(serviceID string) error {
	query := "DELETE FROM services WHERE id = ?"
	_, err := repo.db.Exec(query, serviceID)
	return err
}
func (repo *ServiceRepository) GetServiceByName(serviceName string) (*model.Service, error) {
	query := "SELECT id, name, description, price, provider_id, category FROM services WHERE name = ?"
	var service model.Service
	err := repo.db.QueryRow(query, serviceName).Scan(&service.ID, &service.Name, &service.Description, &service.Price, &service.ProviderID, &service.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("service not found")
		}
		return nil, err
	}
	return &service, nil
}

// GetServiceByProviderID retrieves a service by its ProviderID
func (repo *ServiceRepository) GetServiceByProviderID(providerID string) ([]model.Service, error) {
	query := "SELECT id, name, description, price, provider_id, category FROM services WHERE provider_id = ?"
	rows, err := repo.db.Query(query, providerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("service not found")
		}
		return nil, err
	}
	var services []model.Service
	for rows.Next() {
		var service model.Service
		err := rows.Scan(&service.ID, &service.Name, &service.Description, &service.Price, &service.ProviderID, &service.Category)
		if err != nil {
			return nil, err
		} else {
			services = append(services, service)
		}
	}

	return services, nil
}

func (repo *ServiceRepository) UpdateService(providerID string, updatedService model.Service) error {
	query := "UPDATE services SET name = ?, description = ?, price = ? WHERE provider_id= ? AND id=?;"
	result, err := repo.db.Exec(query, updatedService.Name, updatedService.Description, updatedService.Price, providerID, updatedService.ID)
	// Check how many rows were affected
	if err != nil {
		log.Println("Error executing update query:", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
	}

	if rowsAffected == 0 {
		return errors.New("The service ID may not exist.")
	}
	return nil
}

func (repo *ServiceRepository) RemoveServiceByProviderID(providerID string, serviceID string) error {
	query := "DELETE FROM services WHERE id = ? AND provider_id = ?"
	result, err := repo.db.Exec(query, serviceID, providerID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Return a custom error or handle it as needed
		return errors.New("Invalid service ID")
	}

	return nil
}
