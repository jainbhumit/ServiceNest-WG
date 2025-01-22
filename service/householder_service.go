package service

import (
	"errors"
	"fmt"
	"log"
	"serviceNest/errs"
	"serviceNest/interfaces"
	"serviceNest/model"
	"serviceNest/util"
	"time"
)

var GetUniqueID = util.GenerateUniqueID

type HouseholderService struct {
	householderRepo    interfaces.HouseholderRepository
	providerRepo       interfaces.ServiceProviderRepository
	serviceRepo        interfaces.ServiceRepository
	serviceRequestRepo interfaces.ServiceRequestRepository
}

func NewHouseholderService(householderRepo interfaces.HouseholderRepository, providerRepo interfaces.ServiceProviderRepository, serviceRepo interfaces.ServiceRepository, serviceRequestRepo interfaces.ServiceRequestRepository) interfaces.HouseholderService {
	return &HouseholderService{
		householderRepo:    householderRepo,
		providerRepo:       providerRepo,
		serviceRepo:        serviceRepo,
		serviceRequestRepo: serviceRequestRepo,
	}
}
func (s *HouseholderService) ViewStatus(householderID string, limit, offset int, status string) ([]model.ServiceRequest, error) {
	// Fetch all service requests for the householder
	requests, err := s.serviceRequestRepo.GetServiceRequestsByHouseholderID(householderID, limit, offset, status)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// CancelAcceptedRequest allows a householder to cancel a request that has been accepted by a service provider
func (s *HouseholderService) CancelAcceptedRequest(requestID, householderID string) error {
	// Fetch the service_test request by ID
	serviceRequest, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	// Ensure the service_test request belongs to the householder
	if serviceRequest.HouseholderID == nil || *serviceRequest.HouseholderID != householderID {
		return errors.New(errs.RequestNotBelongToHouseholder)
	}

	// Check if the service request is in "Accepted" status
	if serviceRequest.Status != "Accepted" {
		return errors.New(errs.OnlyAcceptedRequestCancelled)
	}

	// Update the status to "Cancelled"
	serviceRequest.Status = "Cancelled"

	// Save the updated service request
	err = s.serviceRequestRepo.UpdateServiceRequest(serviceRequest)
	if err != nil {
		return err
	}

	return nil
}

// SearchService searches for available service_test providers based on service_test type and proximity
func (s *HouseholderService) SearchService(householder *model.Householder, serviceType string) ([]model.ServiceProvider, error) {
	providers, err := s.providerRepo.GetProvidersByServiceType(serviceType)
	if err != nil {
		return nil, err
	}

	// Example logic: Filter providers by proximity to householder
	nearbyProviders := []model.ServiceProvider{}
	for _, provider := range providers {
		if s.isNearby(householder, &provider) {
			nearbyProviders = append(nearbyProviders, provider)
		}
	}

	return nearbyProviders, nil

}

func (s *HouseholderService) GetServicesByCategory(category string) ([]model.Service, error) {
	// Fetch all services from the service_test repository_test
	services, err := s.serviceRepo.GetServicesByCategory(category)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold the filtered services
	var filteredServices []model.Service

	// Iterate over each service and get providerDetail
	for _, service := range services {

		provider, err := s.providerRepo.GetProviderDetailByID(service.ProviderID, service.ID)
		if err != nil {
			return nil, err
		}
		service.ProviderName = provider.Name
		service.ProviderContact = provider.Contact
		service.ProviderAddress = provider.Address
		service.RatingCount = provider.RatingCount
		filteredServices = append(filteredServices, service)

	}

	return filteredServices, nil
}

// RequestService allows the householder to request a service_test from a provider
func (s *HouseholderService) RequestService(householderID string, serviceName string, category string, description string, scheduleTime *time.Time) (string, error) {
	// Check if a service exists for the given category ID
	serviceId, err := s.serviceRepo.GetServiceIdByCategory(category)
	if err != nil {
		return "", err
	}
	fmt.Println("Service Id is ", serviceId)
	if serviceId == nil {
		return "", errors.New("service category does not exist")
	}

	// Generate a unique ID for the service request
	requestID := GetUniqueID()

	// Fetch householder details from Db
	householder, err := s.householderRepo.GetHouseholderByID(householderID)
	if err != nil {
		return "", err
	}
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Printf("Failed to load location: %v", err)
	}
	fmt.Println(time.Now().In(location))
	// Create the service request
	serviceRequest := model.ServiceRequest{
		ID:                 requestID,
		ServiceName:        serviceName,
		HouseholderName:    householder.Name,
		HouseholderID:      &householder.User.ID,
		HouseholderAddress: &householder.Address,
		HouseholderContact: householder.Contact,
		Description:        description,
		ServiceID:          *serviceId,
		RequestedTime:      time.Now().In(location),
		ScheduledTime:      *scheduleTime,
		Status:             "Pending",
		ApproveStatus:      false,
	}
	fmt.Print(serviceRequest)
	// Save the service request to the repository
	err = s.serviceRequestRepo.SaveServiceRequest(serviceRequest)
	if err != nil {
		return "", err
	}

	return serviceRequest.ID, nil

}

// ViewBookingHistory returns the booking history for a householder
func (s *HouseholderService) ViewBookingHistory(householderID string, limit, offset int, status string) ([]model.ServiceRequest, error) {
	return s.serviceRequestRepo.GetServiceRequestsByHouseholderID(householderID, limit, offset, status)
}

// Helper function to determine if a provider is nearby
func (s *HouseholderService) isNearby(householder *model.Householder, provider *model.ServiceProvider) bool {
	// TODO Implement proximity logic  based on distance between coordinates

	return true
}

// GetAvailableServices fetches all available services from the repository_test
func (s *HouseholderService) GetAvailableServices(limit, offset int) ([]model.Service, error) {
	return s.serviceRepo.GetAllServices(limit, offset)
}

// CancelServiceRequest allows the householder to cancel a service_test request
func (s *HouseholderService) CancelServiceRequest(requestID string, householderID string) error {
	request, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	if *request.HouseholderID != householderID {
		return errors.New(errs.RequestNotBelongToHouseholder)
	}

	// Check if the scheduled time is less than 4 hours away
	currentTime := time.Now()
	fmt.Println(request.ScheduledTime, " ", currentTime, " ", request.ScheduledTime.Sub(currentTime))
	if request.ScheduledTime.Sub(currentTime) < 4*time.Hour {
		return fmt.Errorf(errs.RequestCancellationTooLate)
	}
	if request.Status == "Cancelled" {
		return fmt.Errorf(errs.RequestAlreadyCancelled)
	}

	request.Status = "Cancelled"
	return s.serviceRequestRepo.UpdateServiceRequest(request)
}

// RescheduleServiceRequest allows the householder to reschedule a service_test request
func (s *HouseholderService) RescheduleServiceRequest(requestID string, newTime time.Time, householderID string) error {
	request, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	if *request.HouseholderID != householderID {
		return errors.New(errs.RequestNotBelongToHouseholder)
	}
	if request.Status != "Pending" && request.Status != "Accepted" {
		return fmt.Errorf(errs.OnlyPendingRequestRescheduled)
	}

	request.ScheduledTime = newTime
	return s.serviceRequestRepo.UpdateServiceRequest(request)
}

// ViewServiceRequestStatus returns the status of a specific service_test request
func (s *HouseholderService) ViewServiceRequestStatus(requestID string) (string, error) {
	request, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return "", err
	}
	return request.Status, nil
}

func (s *HouseholderService) AddReview(providerID, householderID, serviceID, comments string, rating float64) error {
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Printf("Failed to load location: %v", err)
	}

	// Create the review object
	review := model.Review{
		ID:            GetUniqueID(),
		ProviderID:    providerID,
		ServiceID:     serviceID,
		HouseholderID: householderID,
		Rating:        rating,
		Comments:      comments,
		ReviewDate:    time.Now().In(location),
	}

	// Save the review in the repository
	err = s.providerRepo.AddReview(review)
	if err != nil {
		return err
	}

	// Recalculate and update the provider's rating
	err = s.providerRepo.UpdateProviderRating(providerID, serviceID, review.Rating)
	if err != nil {
		return errors.New(errs.FailUpdateRating)
	}

	return nil
}

func (s *HouseholderService) ApproveServiceRequest(requestID string, providerID string, householderID string) error {
	// Retrieve the service request by ID
	serviceRequest, err := s.serviceRequestRepo.GetServiceProviderByRequestID(requestID, providerID)
	if err != nil {
		return fmt.Errorf("%v: %v", errs.ServiceRequestNotFound, err)
	}

	if *serviceRequest.HouseholderID != householderID {
		return errors.New(errs.RequestNotBelongToHouseholder)
	}
	// Check if the request has already been approved
	if serviceRequest.ApproveStatus {
		return errors.New(errs.RequestAlreadyApproved)
	}

	// Set the approval status to true
	serviceRequest.ApproveStatus = true
	serviceRequest.Status = "Approved"
	for _, provider := range serviceRequest.ProviderDetails {
		if provider.ServiceProviderID == providerID {
			provider.Approve = 1
			if err := s.providerRepo.UpdateServiceProviderDetailByRequestID(&provider, requestID); err != nil {
				return fmt.Errorf(errs.NotUpdateProviderDetails)
			}
			break
		}
	}
	// Update the service request in the repository
	if err := s.serviceRequestRepo.UpdateServiceRequest(serviceRequest); err != nil {
		return fmt.Errorf("%v: %v", errs.NotUpdateRequest, err)
	}

	return nil
}
func (s *HouseholderService) ViewApprovedRequests(householderID string, limit, offset int, sortOrder string) ([]model.ServiceRequest, error) {
	// Retrieve all service requests for the householder
	serviceRequests, err := s.serviceRequestRepo.GetApproveServiceRequestsByHouseholderID(householderID, limit, offset, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errs.NotRetrieveRequest, err)
	}

	// Filter to only include approved requests
	var approvedRequests []model.ServiceRequest
	for _, req := range serviceRequests {
		if req.ApproveStatus && req.Status != "Cancelled" {
			approvedRequests = append(approvedRequests, req)
		}
	}

	if len(approvedRequests) == 0 {
		return nil, errors.New(errs.NoApproveRequestFound)
	}

	return approvedRequests, nil
}

func (s *HouseholderService) GetAllServiceCategory() ([]model.Category, error) {
	categories, err := s.serviceRepo.GetAllCategory()
	if err != nil {
		return nil, err
	}
	return categories, nil
}
