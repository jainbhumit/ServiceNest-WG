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
	// Fetch all service_test requests for the householder
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
	err = s.serviceRequestRepo.UpdateServiceRequest(*serviceRequest)
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

//	func (s *HouseholderService) GetServicesByCategory(category string) ([]model.Service, error) {
//		services, err := s.serviceRepo.GetAllServices()
//		if err != nil {
//			return nil, err
//		}
//
//		var filteredServices []model.Service
//		for _, service_test := range services {
//			if service_test.Category == category {
//				filteredServices = append(filteredServices, *service_test)
//			}
//		}
//
//		return filteredServices, nil
//	}
func (s *HouseholderService) GetServicesByCategory(category string) ([]*model.Service, error) {
	// Fetch all services from the service_test repository_test
	services, err := s.serviceRepo.GetAllServices()
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold the filtered services
	var filteredServices []*model.Service

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
func (s *HouseholderService) RequestService(householder *model.Householder, serviceID string) (string, error) {
	// Generate a unique ID for the service_test request
	requestID := util.GenerateUniqueID() // Function to generate a unique ID
	// Create the service_test request
	serviceRequest := model.ServiceRequest{
		ID:                 requestID,
		HouseholderName:    householder.Name,
		HouseholderAddress: &householder.Address,
		ServiceID:          serviceID,
		RequestedTime:      time.Now(),
		Status:             "Pending", // Initial status
		ApproveStatus:      false,
	}

	// Save the service_test request to the repository_test
	err := s.serviceRequestRepo.SaveServiceRequest(serviceRequest)
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
func (s *HouseholderService) GetAvailableServices() ([]*model.Service, error) {
	return s.serviceRepo.GetAllServices()
}

// CancelServiceRequest allows the householder to cancel a service_test request
func (s *HouseholderService) CancelServiceRequest(requestID string) error {
	request, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status == "Cancelled" {
		return fmt.Errorf("service_test request is already cancelled")
	}

	request.Status = "Cancelled"
	return s.serviceRequestRepo.UpdateServiceRequest(*request)
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
	return s.serviceRequestRepo.UpdateServiceRequest(*request)
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
func (s *HouseholderService) AddReview(householderID, serviceID, comments string, rating float64) error {
	// Fetch the service_test provider associated with the service_test
	service, err := s.serviceRepo.GetServiceByID(serviceID)
	if err != nil {
		return errors.New("service not found")
	}

	// Fetch the service_test provider offering this service_test
	provider, err := s.providerRepo.GetProviderByServiceID(service.ID)
	if err != nil {
		return errors.New("service provider not found")
	}

	// Create the review
	review := &model.Review{
		ID:            util.GenerateUniqueID(),
		ServiceID:     serviceID,
		HouseholderID: householderID,
		Rating:        rating,
		Comments:      comments,
		ReviewDate:    time.Now(),
	}

	// Append the review to the service_test provider's list of reviews
	provider.Reviews = append(provider.Reviews, review)

	// Recalculate the service_test provider's rating
	totalRating := 0.0
	for _, rev := range provider.Reviews {
		totalRating += rev.Rating
	}
	provider.Rating = totalRating / float64(len(provider.Reviews))

	// Save the updated service_test provider data
	err = s.providerRepo.UpdateServiceProvider(provider)
	if err != nil {
		return errors.New("failed to save review")
	}

	return nil
}

//func (s *HouseholderService) GetServiceDetails(serviceID string) (*model.Service, error) {
//	return s.serviceRepo.GetServiceByID(serviceID)
//}
