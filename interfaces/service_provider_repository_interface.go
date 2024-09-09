package interfaces

import "serviceNest/model"

type ServiceProviderRepository interface {
	UpdateServiceProvider(provider *model.ServiceProvider) error
	GetProviderByServiceID(serviceID string) (*model.ServiceProvider, error)
	GetProvidersByServiceType(serviceType string) ([]model.ServiceProvider, error)
	GetProviderByID(providerID string) (*model.ServiceProvider, error)
	SaveServiceProvider(provider model.ServiceProvider) error
	GetProviderDetailByID(providerID string) (*model.ServiceProviderDetails, error)
	SaveServiceProviderDetail(provider *model.ServiceProviderDetails, requestID string) error
	UpdateServiceProviderDetailByRequestID(provider *model.ServiceProviderDetails, requestID string) error
	IsProviderApproved(providerID string) (bool, error)
	AddReview(review model.Review) error
	UpdateProviderRating(providerID string) error
	GetReviewsByProviderID(providerID string) ([]model.Review, error)
}
