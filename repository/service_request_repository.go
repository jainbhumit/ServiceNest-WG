package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"serviceNest/model"
)

type ServiceRequestRepository struct {
	filePath string
}

// NewServiceRequestRepository initializes a new ServiceRequestRepository
func NewServiceRequestRepository(filePath string) *ServiceRequestRepository {
	return &ServiceRequestRepository{filePath: filePath}
}

// SaveServiceRequest saves a service request to the file
func (repo *ServiceRequestRepository) SaveServiceRequest(request model.ServiceRequest) error {
	requests, err := repo.loadServiceRequests()
	if err != nil {
		return err
	}

	requests = append(requests, request)

	return repo.saveServiceRequests(requests)
}

// GetServiceRequestByID retrieves a service request by its ID
func (repo *ServiceRequestRepository) GetServiceRequestByID(requestID string) (*model.ServiceRequest, error) {
	requests, err := repo.loadServiceRequests()
	if err != nil {
		return nil, err
	}

	for _, request := range requests {
		if request.ID == requestID {
			return &request, nil
		}
	}

	return nil, fmt.Errorf("service request with ID %s not found", requestID)
}

// GetServiceRequestsByHouseholderID retrieves all service requests made by a specific householder
func (repo *ServiceRequestRepository) GetServiceRequestsByHouseholderID(householderID string) ([]model.ServiceRequest, error) {
	requests, err := repo.loadServiceRequests()
	if err != nil {
		return nil, err
	}

	var householderRequests []model.ServiceRequest
	for _, request := range requests {
		if *request.HouseholderID == householderID {
			householderRequests = append(householderRequests, request)
		}
	}

	return householderRequests, nil
}

// UpdateServiceRequest updates an existing service request
func (repo *ServiceRequestRepository) UpdateServiceRequest(updatedRequest model.ServiceRequest) error {
	requests, err := repo.loadServiceRequests()
	if err != nil {
		return err
	}

	for i, request := range requests {
		if request.ID == updatedRequest.ID {
			requests[i] = updatedRequest
			break
		}
	}

	return repo.saveServiceRequests(requests)
}

// GetAllServiceRequests retrieves all service requests from the file
func (repo *ServiceRequestRepository) GetAllServiceRequests() ([]model.ServiceRequest, error) {
	return repo.loadServiceRequests()
}

// SaveAllServiceRequests saves all service requests to the file
func (r *ServiceRequestRepository) SaveAllServiceRequests(serviceRequests []model.ServiceRequest) error {
	// Convert the service requests slice to JSON
	data, err := json.MarshalIndent(serviceRequests, "", "  ")
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

// Private helper methods for loading and saving service requests
func (repo *ServiceRequestRepository) loadServiceRequests() ([]model.ServiceRequest, error) {
	var serviceRequests []model.ServiceRequest

	// Check if the file exists
	if _, err := os.Stat(repo.filePath); os.IsNotExist(err) {
		// File does not exist, return an empty slice
		return serviceRequests, nil
	}

	// File exists, proceed to read it
	file, err := ioutil.ReadFile(repo.filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	err = json.Unmarshal(file, &serviceRequests)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal file: %v", err)
	}

	return serviceRequests, nil
}

func (repo *ServiceRequestRepository) saveServiceRequests(requests []model.ServiceRequest) error {
	data, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(repo.filePath, data, 0644)
}
