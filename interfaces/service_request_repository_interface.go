package interfaces

import "serviceNest/model"

type ServiceRequestRepository interface {
	//SaveAllServiceRequests(serviceRequests []model.ServiceRequest) error
	GetAllServiceRequests(limit, offset int) ([]model.ServiceRequest, error)
	UpdateServiceRequest(updatedRequest *model.ServiceRequest) error
	GetServiceRequestsByHouseholderID(householderID string, limit, offset int, status string) ([]model.ServiceRequest, error)
	GetServiceRequestByID(requestID string) (*model.ServiceRequest, error)
	SaveServiceRequest(request model.ServiceRequest) error
	GetServiceRequestsByProviderID(providerID string, limit, offset int, sortOrder string) ([]model.ServiceRequest, error)
	GetServiceProviderByRequestID(requestID, providerID string) (*model.ServiceRequest, error)
	GetApproveServiceRequestsByHouseholderID(householderID string, limit, offset int, sortOrder string) ([]model.ServiceRequest, error)
	GetAllPendingRequestsByProvider(providerId string, serviceID string, limit, offset int) ([]model.ServiceRequest, error)
}
