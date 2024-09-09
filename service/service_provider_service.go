package service

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"serviceNest/interfaces"
	"serviceNest/model"
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
	// Save the updated service provider information
	err := s.serviceRepo.UpdateService(providerID, updatedService)
	if err != nil {
		return err
	}

	// Update the service in the service repository
	return nil
}
func (s *ServiceProviderService) GetAllServiceRequests() ([]model.ServiceRequest, error) {
	return s.serviceRequestRepo.GetAllServiceRequests()
}

func (s *ServiceProviderService) RemoveService(providerID, serviceID string) error {
	err := s.serviceRepo.RemoveServiceByProviderID(providerID, serviceID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceProviderService) AcceptServiceRequest(providerID, requestID string) error {
	serviceRequest, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	if serviceRequest.ApproveStatus {
		return fmt.Errorf("service request has already been approved")
	}

	// Update the service request status to "Accepted"
	serviceRequest.Status = "Accepted"

	// Get the ServiceProvider details
	provider, err := s.serviceProviderRepo.GetProviderDetailByID(providerID)
	if err != nil {
		return err
	}
	providerReviews, err := s.serviceProviderRepo.GetReviewsByProviderID(providerID)
	if err != nil {
		return err
	}
	provider.ServiceProviderID = providerID
	var estimatedPrice string
	color.Cyan("Enter the Price for service:")
	fmt.Scanln(&estimatedPrice)

	provider.Price = estimatedPrice

	// Add ServiceProvider details to the ServiceRequest
	serviceRequest.ProviderDetails = append(serviceRequest.ProviderDetails, model.ServiceProviderDetails{
		ServiceProviderID: providerID,
		Name:              provider.Name,
		Contact:           provider.Contact,
		Address:           provider.Address,
		Price:             estimatedPrice,
		Rating:            provider.Rating,
		Reviews:           providerReviews,
	})

	// Save the updated service request
	err = s.serviceRequestRepo.UpdateServiceRequest(serviceRequest)

	err = s.serviceProviderRepo.SaveServiceProviderDetail(provider, requestID)
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
	// Get the service request
	request, err := s.serviceRequestRepo.GetServiceRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status != "Pending" {
		return fmt.Errorf("service request is not pending")
	}

	// Decline the service_test request
	request.Status = "Declined"
	return s.serviceRequestRepo.UpdateServiceRequest(request)
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

//// ViewServices returns all services offered by a specific service_test provider
//func (s *ServiceProviderService) ViewServices(providerID string) ([]model.Service, error) {
//	provider, err := s.serviceProviderRepo.GetProviderByID(providerID)
//	if err != nil {
//		return nil, err
//	}
//
//	return provider.ServicesOffered, nil
//}

func (s *ServiceProviderService) ViewServices(providerID string) ([]model.Service, error) {
	providerService, err := s.serviceRepo.GetServiceByProviderID(providerID)
	if err != nil {
		return nil, err
	}

	return providerService, nil
}
func (s *ServiceProviderService) GetServiceByID(serviceID string) (*model.Service, error) {
	return s.serviceRepo.GetServiceByID(serviceID)
}

func (s *ServiceProviderService) ViewApprovedRequestsByHouseholder(providerID string) ([]model.ServiceRequest, error) {
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

func (s *ServiceProviderService) GetReviews(providerID string) ([]model.Review, error) {
	reviews, err := s.serviceProviderRepo.GetReviewsByProviderID(providerID)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
