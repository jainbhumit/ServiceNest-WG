package interfaces

import "serviceNest/model"

type ServiceProviderRepository interface {
	AddReview(providerID, householderID, review string, rating float64) error
	UpdateServiceProvider(provider *model.ServiceProvider) error
	GetProviderByServiceID(serviceID string) (*model.ServiceProvider, error)
	GetProvidersByServiceType(serviceType string) ([]model.ServiceProvider, error)
	GetProviderByID(providerID string) (*model.ServiceProvider, error)
	SaveServiceProvider(provider model.ServiceProvider) error
}
