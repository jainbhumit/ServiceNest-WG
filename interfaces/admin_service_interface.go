package interfaces

import "serviceNest/model"

type AdminService interface {
	GetAllService() ([]model.Service, error)
	DeactivateAccount(userID string) error
	DeleteService(serviceID string) error
	ViewReports() ([]model.ServiceRequest, error)
}
