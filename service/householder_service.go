package service

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
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

func NewHouseholderService(householderRepo interfaces.HouseholderRepository, providerRepo interfaces.ServiceProviderRepository, serviceRepo interfaces.ServiceRepository, serviceRequestRepo interfaces.ServiceRequestRepository) *HouseholderService {
	return &HouseholderService{
		householderRepo:    householderRepo,
		providerRepo:       providerRepo,
		serviceRepo:        serviceRepo,
		serviceRequestRepo: serviceRequestRepo,
	}
}
func (s *HouseholderService) ViewStatus(serviceRequestRepo *HouseholderService, householder *model.Householder) ([]model.ServiceRequest, error) {
	// Fetch all service requests for the householder
	requests, err := s.serviceRequestRepo.GetServiceRequestsByHouseholderID(householder.ID)
	if err != nil {
		color.Red("Error fetching service requests: %v", err)
		return nil, err
	}
	return requests, nil
}

// CancelAcceptedRequest allows a householder to cancel a request that has been accepted by a service_test provider
func (s *HouseholderService) CancelAcceptedRequest(requestID, householderID string) error {
	// Fetch the service_test request by ID
	serviceRequest, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	// Ensure the service_test request belongs to the householder
	if serviceRequest.HouseholderID == nil || *serviceRequest.HouseholderID != householderID {
		return errors.New("service request does not belong to the householder")
	}

	// Check if the service_test request is in "Accepted" status
	if serviceRequest.Status != "Accepted" {
		return errors.New("only accepted service requests can be canceled")
	}

	// Update the status to "Cancelled"
	serviceRequest.Status = "Cancelled"

	// Save the updated service_test request
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
	services, err := s.serviceRepo.GetAllServices()
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold the filtered services
	var filteredServices []model.Service

	// Iterate over each service_test and filter by category
	for _, service := range services {
		if service.Category == category {
			//// Fetch the service_test provider details using the ProviderID from the service_test
			//provider, err := s.getProviderDetails(service_test.ProviderID)
			//if err != nil {
			//	return nil, err
			//}
			//
			//// Attach the provider details to the service_test object
			//service_test.ProviderName = provider.Name
			//service_test.ProviderContact = provider.Contact
			//service_test.ProviderAddress = provider.Address
			//service_test.ProviderRating = provider.Rating

			// Add the service_test with the provider details to the filtered services slice
			provider, err := s.providerRepo.GetProviderDetailByID(service.ProviderID)
			if err != nil {
				return nil, err
			}
			service.ProviderName = provider.Name
			service.ProviderContact = provider.Contact
			service.ProviderAddress = provider.Address
			filteredServices = append(filteredServices, service)
		}
	}

	return filteredServices, nil
}

// getProviderDetails is a helper method to fetch provider details by ProviderID
//func (s *HouseholderService) getProviderDetails(providerID string) (*model.ServiceProvider, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	// Assume serviceProviderRepo is an instance of ServiceProviderRepository
//	var provider model.ServiceProvider
//	err := s.providerRepo.Collection.FindOne(ctx, bson.M{"id": providerID}).Decode(&provider)
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return nil, errors.New("service_test provider not found")
//		}
//		return nil, err
//	}
//
//	return &provider, nil
//}

// RequestService allows the householder to request a service_test from a provider
func (s *HouseholderService) RequestService(householder *model.Householder, serviceName string, scheduleTime *time.Time) (string, error) {
	// Check if the service already exists
	service, err := s.serviceRepo.GetServiceByName(serviceName)
	if err != nil && err.Error() != "service not found" {
		return "", err
	}

	var serviceID string
	if service == nil {
		// Service does not exist, create a custom service entry
		customServiceID := util.GenerateUniqueID() // Function to generate a unique ID
		customService := model.Service{
			ID:          customServiceID,
			Name:        serviceName,
			Description: "Custom Service Request",
			Price:       0.0,      // Placeholder price
			ProviderID:  "",       // No provider assigned
			Category:    "Custom", // Assign a category if needed
		}
		// Save the custom service to the repository
		err := s.serviceRepo.SaveService(customService)
		if err != nil {
			return "", err
		}
		serviceID = customServiceID
	} else {
		serviceID = service.ID
	}

	// Generate a unique ID for the service request
	requestID := GetUniqueID()

	// Create the service request
	serviceRequest := model.ServiceRequest{
		ID:                 requestID,
		HouseholderName:    householder.Name,
		HouseholderID:      &householder.User.ID,
		HouseholderAddress: &householder.Address,
		ServiceID:          serviceID,
		RequestedTime:      time.Now(),
		ScheduledTime:      *scheduleTime,
		Status:             "Pending",
		ApproveStatus:      false,
	}

	// Save the service request to the repository
	err = s.serviceRequestRepo.SaveServiceRequest(serviceRequest)
	if err != nil {
		return "", err
	}

	return serviceRequest.ID, nil
}

// ViewBookingHistory returns the booking history for a householder
func (s *HouseholderService) ViewBookingHistory(householderID string) ([]model.ServiceRequest, error) {
	return s.serviceRequestRepo.GetServiceRequestsByHouseholderID(householderID)
}

// ReviewServiceProvider allows the householder to leave a review for a service_test provider
//func (s *HouseholderService) ReviewServiceProvider(householderID, providerID, review string, rating float64) error {
//	return s.providerRepo.AddReview(providerID, householderID, review, rating)
//}

// Helper function to determine if a provider is nearby
func (s *HouseholderService) isNearby(householder *model.Householder, provider *model.ServiceProvider) bool {
	// Implement proximity logic here (e.g., based on distance between coordinates)
	// For simplicity, this could return true or false based on some criteria
	return true
}

// GetAvailableServices fetches all available services from the repository_test
func (s *HouseholderService) GetAvailableServices() ([]model.Service, error) {
	return s.serviceRepo.GetAllServices()
}

// CancelServiceRequest allows the householder to cancel a service_test request
func (s *HouseholderService) CancelServiceRequest(requestID string) error {
	request, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status == "Cancelled" {
		return fmt.Errorf("service request is already cancelled")
	}

	request.Status = "Cancelled"
	return s.serviceRequestRepo.UpdateServiceRequest(request)
}

// RescheduleServiceRequest allows the householder to reschedule a service_test request
func (s *HouseholderService) RescheduleServiceRequest(requestID string, newTime time.Time) error {
	request, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status != "Pending" && request.Status != "Accepted" {
		return fmt.Errorf("only pending or accepted requests can be rescheduled")
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

// AddReview allows the householder to add a review for a service_test provided by a service_test provider
//func (s *HouseholderService) AddReview(householderID, serviceID, comments string, rating float64) error {
//	// Fetch the service_test provider associated with the service_test
//	service, err := s.serviceRepo.GetServiceByID(serviceID)
//	if err != nil {
//		return errors.New("service not found")
//	}
//
//	// Fetch the service_test provider offering this service_test
//	provider, err := s.providerRepo.GetProviderByServiceID(service.ID)
//	if err != nil {
//		return errors.New("service provider not found")
//	}
//
//	// Create the review
//	review := &model.Review{
//		ID:            util.GenerateUniqueID(),
//		ServiceID:     serviceID,
//		HouseholderID: householderID,
//		Rating:        rating,
//		Comments:      comments,
//		ReviewDate:    time.Now(),
//	}
//
//	// Append the review to the service_test provider's list of reviews
//	provider.Reviews = append(provider.Reviews, review)
//
//	// Recalculate the service_test provider's rating
//	totalRating := 0.0
//	for _, rev := range provider.Reviews {
//		totalRating += rev.Rating
//	}
//	provider.Rating = totalRating / float64(len(provider.Reviews))
//
//	// Save the updated service_test provider data
//	err = s.providerRepo.UpdateServiceProvider(provider)
//	if err != nil {
//		return errors.New("failed to save review")
//	}
//
//	return nil
//}

func (s *HouseholderService) AddReview(providerID, householderID, serviceID, comments string, rating float64) error {
	// Create the review object
	review := model.Review{
		ID:            GetUniqueID(),
		ProviderID:    providerID, // Include providerID
		ServiceID:     serviceID,
		HouseholderID: householderID,
		Rating:        rating,
		Comments:      comments,
		ReviewDate:    time.Now(),
	}

	// Save the review in the repository
	err := s.providerRepo.AddReview(review)
	if err != nil {
		return err
	}

	// Recalculate and update the provider's rating
	err = s.providerRepo.UpdateProviderRating(providerID)
	if err != nil {
		return errors.New("failed to update provider rating")
	}

	return nil
}

// ApproveServiceRequest allows the householder to approve a service request.
//
//	func (s *HouseholderService) ApproveServiceRequest(requestID string, providerID string) error {
//		// Retrieve the service request by ID
//		serviceRequest, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
//		if err != nil {
//			return fmt.Errorf("could not find service request: %v", err)
//		}
//
//		// Check if the request has already been approved
//		if serviceRequest.ApproveStatus {
//			return errors.New("service request has already been approved")
//		}
//
//
//		// Set the approval status to true
//		serviceRequest.ApproveStatus = true
//		for _, provider := range serviceRequest.ProviderDetails {
//			if provider.ServiceProviderID == providerID {
//				provider.Approve = true
//				break
//			}
//		}
//		// Update the service request in the repository
//		if err := s.serviceRequestRepo.UpdateServiceRequest(*serviceRequest); err != nil {
//			return fmt.Errorf("could not update service request: %v", err)
//		}
//
//		return nil
//	}
func (s *HouseholderService) ApproveServiceRequest(requestID string, providerID string) error {
	// Retrieve the service request by ID
	serviceRequest, err := s.serviceRequestRepo.GetServiceProviderByRequestID(requestID, providerID)
	if err != nil {
		return fmt.Errorf("could not find service request: %v", err)
	}

	// Check if the request has already been approved
	if serviceRequest.ApproveStatus {
		return errors.New("service request has already been approved")
	}

	// Set the approval status to true
	serviceRequest.ApproveStatus = true
	for _, provider := range serviceRequest.ProviderDetails {
		if provider.ServiceProviderID == providerID {
			provider.Approve = true
			if err := s.providerRepo.UpdateServiceProviderDetailByRequestID(&provider, requestID); err != nil {
				return fmt.Errorf("could not update service provider detail")
			}
			break
		}
	}
	// Update the service request in the repository
	if err := s.serviceRequestRepo.UpdateServiceRequest(serviceRequest); err != nil {
		return fmt.Errorf("could not update service request: %v", err)
	}

	return nil
}
func (s *HouseholderService) ViewApprovedRequests(householderID string) ([]model.ServiceRequest, error) {
	// Retrieve all service requests for the householder
	serviceRequests, err := s.serviceRequestRepo.GetServiceRequestsByHouseholderID(householderID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service requests: %v", err)
	}

	// Filter to only include approved requests
	var approvedRequests []model.ServiceRequest
	for _, req := range serviceRequests {
		if req.ApproveStatus {
			approvedRequests = append(approvedRequests, req)
		}
	}

	if len(approvedRequests) == 0 {
		return nil, errors.New("no approved service requests found")
	}

	return approvedRequests, nil
}
