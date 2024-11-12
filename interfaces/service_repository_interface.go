package interfaces

import "serviceNest/model"

type ServiceRepository interface {
	RemoveService(serviceID string) error
	SaveService(service model.Service) error
	GetAllServices(limit, offset int) ([]model.Service, error)
	GetServicesByCategory(category string) ([]model.Service, error)
	GetServiceByID(serviceID string) (*model.Service, error)
	GetServiceIdByCategory(category string) (*string, error)
	GetServiceByProviderID(providerID string) ([]model.Service, error)
	UpdateService(providerID string, updatedService model.Service) error
	RemoveServiceByProviderID(providerID string, serviceID string) error
	CategoryExists(categoryName string) (bool, error)
	GetAllCategory() ([]model.Category, error)
	AddCategory(category *model.Category) error
}
