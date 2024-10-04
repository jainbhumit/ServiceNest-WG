package interfaces

import (
	"serviceNest/model"
	"time"
)

type HouseholderService interface {
	ViewApprovedRequests(householderID string) ([]model.ServiceRequest, error)
	ApproveServiceRequest(requestID string, providerID string, householderID string) error
	AddReview(providerID, householderID, serviceID, comments string, rating float64) error
	ViewServiceRequestStatus(requestID string) (string, error)
	RescheduleServiceRequest(requestID string, newTime time.Time, householderID string) error
	CancelServiceRequest(requestID string, householderID string) error
	GetAvailableServices() ([]model.Service, error)
	ViewBookingHistory(householderID string) ([]model.ServiceRequest, error)
	RequestService(householderID string, serviceName string, scheduleTime *time.Time) (string, error)
	GetServicesByCategory(category string) ([]model.Service, error)
	SearchService(householder *model.Householder, serviceType string) ([]model.ServiceProvider, error)
	CancelAcceptedRequest(requestID, householderID string) error
	ViewStatus(householderID string) ([]model.ServiceRequest, error)
}
