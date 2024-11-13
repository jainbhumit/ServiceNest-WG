package interfaces

import "serviceNest/model"

type ServiceProviderService interface {
	GetReviews(providerID string, limit, offset int, serviceID string) ([]model.Review, error)
	ViewApprovedRequestsByProvider(providerID string, limit, offset int, sortOrder string) ([]model.ServiceRequest, error)
	GetServiceByID(serviceID string) (*model.Service, error)
	ViewServices(providerID string) ([]model.Service, error)
	UpdateAvailability(providerID string, availability bool) error
	DeclineServiceRequest(providerID, requestID string) error
	GetServiceRequestByID(requestID string) (*model.ServiceRequest, error)
	AcceptServiceRequest(providerID, requestID string, estimatedPrice string) error
	RemoveService(providerID, serviceID string) error
	GetAllServiceRequests(providerId string, serviceID string, limit, offset int) ([]model.ServiceRequest, error)
	UpdateService(providerID, serviceID string, updatedService model.Service) error
	AddService(providerID string, newService model.Service) (string, error)
}
