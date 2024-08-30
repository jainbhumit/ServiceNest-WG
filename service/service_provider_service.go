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

type ServiceProviderService struct {
	serviceProviderRepo interfaces.ServiceProviderRepository
	serviceRequestRepo  interfaces.ServiceRequestRepository
	serviceRepo         interfaces.ServiceRepository
}

// NewServiceProviderService initializes a new ServiceProviderService
func NewServiceProviderService(serviceProviderRepo interfaces.ServiceProviderRepository, serviceRequestRepo interfaces.ServiceRequestRepository, serviceRepo interfaces.ServiceRepository) *ServiceProviderService {
	return &ServiceProviderService{
		serviceProviderRepo: serviceProviderRepo,
		serviceRequestRepo:  serviceRequestRepo,
		serviceRepo:         serviceRepo,
	}
}

// AddService adds a new service_test to the provider's list of offered services
func (s *ServiceProviderService) AddService(providerID string, newService model.Service) error {
	// Get the service_test provider
	provider, err := s.serviceProviderRepo.GetProviderByID(providerID)
	if err != nil {
		return err
	}

	// Add the new service_test to the provider's list
	provider.ServicesOffered = append(provider.ServicesOffered, newService)

	// Save the updated service_test provider information
	err = s.serviceProviderRepo.UpdateServiceProvider(provider)
	if err != nil {
		return err
	}

	// Save the new service_test to the service_test repository_test
	return s.serviceRepo.SaveService(newService)
}

// UpdateService updates an existing service_test offered by the provider
func (s *ServiceProviderService) UpdateService(providerID, serviceID string, updatedService model.Service) error {
	// Get the service_test provider
	provider, err := s.serviceProviderRepo.GetProviderByID(providerID)
	if err != nil {
		return err
	}

	// Update the service_test in the provider's list
	for i, service := range provider.ServicesOffered {
		if service.ID == serviceID {
			provider.ServicesOffered[i] = updatedService
			break
		}
	}

	// Save the updated service_test provider information
	err = s.serviceProviderRepo.UpdateServiceProvider(provider)
	if err != nil {
		return err
	}

	// Update the service_test in the service_test repository_test
	return s.serviceRepo.SaveService(updatedService)
}

func (s *ServiceProviderService) GetAllServiceRequests() ([]model.ServiceRequest, error) {
	return s.serviceRequestRepo.GetAllServiceRequests()
}

// RemoveService removes a service_test from the provider's list of offered services
func (s *ServiceProviderService) RemoveService(providerID, serviceID string) error {
	// Get the service_test provider
	provider, err := s.serviceProviderRepo.GetProviderByID(providerID)
	if err != nil {
		return err
	}

	// Filter out the service_test from the provider's list
	var updatedServices []model.Service
	serviceExists := false
	for _, service := range provider.ServicesOffered {
		if service.ID != serviceID {
			updatedServices = append(updatedServices, service)
		} else {
			serviceExists = true
		}
	}

	// If the service_test was not found in the provider's list, return an error
	if !serviceExists {
		return fmt.Errorf("service with ID %s not found for provider %s", serviceID, providerID)
	}

	provider.ServicesOffered = updatedServices

	// Save the updated service_test provider information
	err = s.serviceProviderRepo.UpdateServiceProvider(provider)
	if err != nil {
		return err
	}

	// Remove the service_test from the main services repository_test
	err = s.serviceRepo.RemoveService(serviceID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceProviderService) AddReview(serviceID, householderID, comments string, rating float64) error {
	// Generate a new Review ID
	reviewID := util.GenerateUniqueID()

	// Get the service_test provider by the serviceID
	provider, err := s.serviceProviderRepo.GetProviderByID(serviceID)
	if err != nil {
		return err
	}

	// Create a new review
	review := &model.Review{
		ID:            reviewID,
		ServiceID:     serviceID,
		HouseholderID: householderID,
		Rating:        rating,
		Comments:      comments,
		ReviewDate:    time.Now(),
	}

	// Add the review to the provider's list of reviews
	provider.Reviews = append(provider.Reviews, review)

	// Update the provider's rating (you can implement a method to calculate the average rating if needed)
	s.updateProviderRating(provider)

	// Save the updated provider information
	err = s.serviceProviderRepo.UpdateServiceProvider(provider)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceProviderService) updateProviderRating(provider *model.ServiceProvider) {
	totalRating := 0.0
	for _, review := range provider.Reviews {
		totalRating += review.Rating
	}
	provider.Rating = totalRating / float64(len(provider.Reviews))
}

// AcceptServiceRequest allows the provider to accept a service_test request
func (s *ServiceProviderService) AcceptServiceRequest(providerID, requestID string) error {
	serviceRequest, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	if serviceRequest.ApproveStatus != false {
		return fmt.Errorf("service request is not pending")
	}

	// Update the service_test request status to "Accepted"
	serviceRequest.Status = "Accepted"

	// Get the ServiceProvider details
	provider, err := s.serviceProviderRepo.GetProviderByID(providerID)
	if err != nil {
		return err
	}
	var estimatedPrice string
	color.Cyan("Enter the Price for service")
	fmt.Scanln(&estimatedPrice)
	// Add ServiceProvider details to the ServiceRequest

	serviceRequest.ProviderDetails = append(serviceRequest.ProviderDetails, &model.ServiceProviderDetails{
		ServiceProviderID: providerID,
		Name:              provider.Name,
		Contact:           provider.Contact,
		Address:           provider.Address,
		Price:             estimatedPrice,
		Rating:            provider.Rating,
		Reviews:           provider.Reviews,
	})

	// Save the updated service_test request
	err = s.serviceRequestRepo.UpdateServiceRequest(*serviceRequest)
	if err != nil {
		return err
	}

	return nil
}
func (s *ServiceProviderService) GetServiceRequestByID(requestID string) (*model.ServiceRequest, error) {
	return s.serviceRequestRepo.GetServiceRequestByID(requestID)
}

// DeclineServiceRequest allows the provider to decline a service_test request
func (s *ServiceProviderService) DeclineServiceRequest(providerID, requestID string) error {
	// Get the service_test request
	request, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status != "Pending" {
		return fmt.Errorf("service request is not pending")
	}

	// Decline the service_test request
	request.Status = "Declined"
	return s.serviceRequestRepo.UpdateServiceRequest(*request)
}

// UpdateAvailability updates the provider's availability status
func (s *ServiceProviderService) UpdateAvailability(providerID string, availability bool) error {
	// Get the service_test provider
	provider, err := s.serviceProviderRepo.GetProviderByID(providerID)
	if err != nil {
		return err
	}

	// Update the availability status
	provider.Availability = availability
	return s.serviceProviderRepo.UpdateServiceProvider(provider)
}

// ViewServices returns all services offered by a specific service_test provider
func (s *ServiceProviderService) ViewServices(providerID string) ([]model.Service, error) {
	provider, err := s.serviceProviderRepo.GetProviderByID(providerID)
	if err != nil {
		return nil, err
	}

	return provider.ServicesOffered, nil
}
func (s *ServiceProviderService) GetServiceByID(serviceID string) (*model.Service, error) {
	return s.serviceRepo.GetServiceByID(serviceID)
}

// ViewReviews retrieves all reviews for a specific service_test provider
func (s *ServiceProviderService) ViewReviews(providerID string) ([]*model.Review, error) {
	// Get the service_test provider
	provider, err := s.serviceProviderRepo.GetProviderByID(providerID)
	if err != nil {
		return nil, err
	}

	// Return all reviews for this provider
	return provider.Reviews, nil
}
func (s *ServiceProviderService) ViewApprovedRequests(providerID string) ([]model.ServiceRequest, error) {
	// Fetch all service requests related to the provider
	serviceRequests, err := s.serviceRequestRepo.GetServiceRequestsByProviderID(providerID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service requests: %v", err)
	}

	// Filter out only the approved requests
	var approvedRequests []model.ServiceRequest
	for _, req := range serviceRequests {
		if req.ApproveStatus {
			for _, providerDetail := range req.ProviderDetails {
				if providerDetail.ServiceProviderID == providerID && providerDetail.Approve {
					approvedRequests = append(approvedRequests, req)
					break
				}
			}
		}
	}

	if len(approvedRequests) == 0 {
		return nil, errors.New("no approved requests found for this provider")
	}

	return approvedRequests, nil
}
