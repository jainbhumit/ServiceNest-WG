package interfaces

import "serviceNest/model"

type AdminService interface {
	GetAllService(limit, offset int) ([]model.Service, error)
	DeactivateAccount(userID string) error
	DeleteService(serviceID string) error
	ViewReports(limit, offset int) ([]model.ServiceRequest, error)
	AddService(name, description string) error
	GetUserByEmail(userEmail string) (*model.User, error)
}
