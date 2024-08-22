package repository

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"serviceNest/model"
)

type ServiceRepository struct {
	filePath string
}

// NewServiceRepository creates a new instance of ServiceRepository
func NewServiceRepository(filePath string) *ServiceRepository {
	return &ServiceRepository{filePath: filePath}
}

// GetAllServices fetches all available services from the file-based database
func (repo *ServiceRepository) GetAllServices() ([]model.Service, error) {
	var services []model.Service

	// Load services from the file
	file, err := ioutil.ReadFile(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, return an empty list and no error
			return services, nil
		}
		return nil, err
	}

	// Unmarshal the JSON data into the services slice
	err = json.Unmarshal(file, &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

// GetServiceByID retrieves a service by its ID
func (r *ServiceRepository) GetServiceByID(serviceID string) (*model.Service, error) {
	// Open the services file
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the data into a slice of services
	var services []model.Service
	err = json.Unmarshal(byteValue, &services)
	if err != nil {
		return nil, err
	}

	// Search for the service by ID
	for _, service := range services {
		if service.Name == serviceID {
			return &service, nil
		}
	}

	// Return an error if the service is not found
	return nil, errors.New("service not found")
}

// SaveService adds a new service to the repository
func (repo *ServiceRepository) SaveService(service model.Service) error {
	services, err := repo.GetAllServices()
	if err != nil {
		return err
	}

	services = append(services, service)

	// Marshal the updated services slice back to JSON
	data, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		return err
	}

	// Write the updated JSON data back to the file
	return ioutil.WriteFile(repo.filePath, data, 0644)
}

// SaveAllServices saves the entire list of services to the file
func (r *ServiceRepository) SaveAllServices(services []model.Service) error {
	// Convert the services slice to JSON
	data, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	err = ioutil.WriteFile(r.filePath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
func (r *ServiceRepository) RemoveService(serviceID string) error {
	// Load all services
	services, err := r.GetAllServices()
	if err != nil {
		return err
	}

	// Filter out the service to be removed
	var updatedServices []model.Service
	for _, service := range services {
		if service.ID != serviceID {
			updatedServices = append(updatedServices, service)
		}
	}

	// Save the updated services back to the file
	return r.SaveAllServices(updatedServices)
}
