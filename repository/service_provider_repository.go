package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"serviceNest/model"
)

type ServiceProviderRepository struct {
	filePath string
}

// NewServiceProviderRepository initializes a new ServiceProviderRepository
func NewServiceProviderRepository(filePath string) *ServiceProviderRepository {
	return &ServiceProviderRepository{filePath: filePath}
}

// SaveServiceProvider saves a new service provider to the file
func (repo *ServiceProviderRepository) SaveServiceProvider(provider model.ServiceProvider) error {
	providers, err := repo.loadServiceProviders()
	if err != nil {
		return err
	}

	providers = append(providers, provider)

	return repo.saveServiceProviders(providers)
}

// GetProviderByID retrieves a service provider by their ID
func (repo *ServiceProviderRepository) GetProviderByID(id string) (*model.ServiceProvider, error) {
	providers, err := repo.loadServiceProviders()
	if err != nil {
		return nil, err
	}

	for _, provider := range providers {
		if provider.ID == id {
			return &provider, nil
		}
	}

	return nil, fmt.Errorf("service provider not found")
}

// GetProvidersByServiceType retrieves all service providers that offer a specific service type
func (repo *ServiceProviderRepository) GetProvidersByServiceType(serviceType string) ([]model.ServiceProvider, error) {
	providers, err := repo.loadServiceProviders()
	if err != nil {
		return nil, err
	}

	var matchingProviders []model.ServiceProvider
	for _, provider := range providers {
		for _, service := range provider.ServicesOffered {
			if service.Name == serviceType {
				matchingProviders = append(matchingProviders, provider)
				break
			}
		}
	}

	return matchingProviders, nil
}

// GetProviderByServiceID retrieves all service providers that offer a specific service type
func (repo *ServiceProviderRepository) GetProviderByServiceID(serviceID string) (*model.ServiceProvider, error) {
	providers, err := repo.loadServiceProviders()
	if err != nil {
		return nil, err
	}

	var matchingProviders model.ServiceProvider
	for _, provider := range providers {
		for _, service := range provider.ServicesOffered {
			if service.ID == serviceID {
				matchingProviders = provider
				break
			}
		}
	}

	return &matchingProviders, nil
}

// UpdateServiceProvider updates an existing service provider
func (repo *ServiceProviderRepository) UpdateServiceProvider(updatedProvider *model.ServiceProvider) error {
	providers, err := repo.loadServiceProviders()
	if err != nil {
		return err
	}

	// Find and update the service provider
	for i, provider := range providers {
		if provider.User.ID == updatedProvider.User.ID {
			providers[i] = *updatedProvider
			break
		}
	}

	// Save the updated list of service providers back to the file
	return repo.saveServiceProviders(providers)
}

// AddReview adds a review and rating for a service provider
func (repo *ServiceProviderRepository) AddReview(providerID, householderID, review string, rating float64) error {
	providers, err := repo.loadServiceProviders()
	if err != nil {
		return err
	}

	for i, provider := range providers {
		if provider.ID == providerID {
			newReview := model.Review{
				HouseholderID: householderID,
				Comments:      review,
				Rating:        rating,
			}
			providers[i].Reviews = append(providers[i].Reviews, &newReview)

			return repo.saveServiceProviders(providers)
		}
	}

	return fmt.Errorf("service provider not found")
}

// Private helper methods for loading and saving service providers
func (repo *ServiceProviderRepository) loadServiceProviders() ([]model.ServiceProvider, error) {
	var providers []model.ServiceProvider

	file, err := ioutil.ReadFile(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return providers, nil
		}
		return nil, err
	}

	err = json.Unmarshal(file, &providers)
	if err != nil {
		return nil, err
	}

	return providers, nil
}

func (repo *ServiceProviderRepository) saveServiceProviders(providers []model.ServiceProvider) error {
	data, err := json.MarshalIndent(providers, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(repo.filePath, data, 0644)
}
