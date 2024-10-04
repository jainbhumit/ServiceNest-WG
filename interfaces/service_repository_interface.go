package interfaces

import "serviceNest/model"

type ServiceRepository interface {
	RemoveService(serviceID string) error
	SaveService(service model.Service) error
	GetAllServices() ([]model.Service, error)
	GetServiceByID(serviceID string) (*model.Service, error)
	GetServiceByName(serviceName string) (*model.Service, error)
	GetServiceByProviderID(providerID string) ([]model.Service, error)
	UpdateService(providerID string, updatedService model.Service) error
	RemoveServiceByProviderID(providerID string, serviceID string) error
}
