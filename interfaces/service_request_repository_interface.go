package interfaces

import "serviceNest/model"

type ServiceRequestRepository interface {
	SaveAllServiceRequests(serviceRequests []model.ServiceRequest) error
	GetAllServiceRequests() ([]model.ServiceRequest, error)
	UpdateServiceRequest(updatedRequest model.ServiceRequest) error
	GetServiceRequestsByHouseholderID(householderID string) ([]model.ServiceRequest, error)
	GetServiceRequestByID(requestID string) (*model.ServiceRequest, error)
	SaveServiceRequest(request model.ServiceRequest) error
}
