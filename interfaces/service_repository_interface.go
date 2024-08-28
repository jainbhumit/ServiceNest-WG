package interfaces

import "serviceNest/model"

type ServiceRepository interface {
	RemoveService(serviceID string) error
	SaveAllServices(services []model.Service) error
	SaveService(service model.Service) error
	GetAllServices() ([]*model.Service, error)
	GetServiceByID(serviceID string) (*model.Service, error)
}
