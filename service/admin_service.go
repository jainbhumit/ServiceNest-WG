package service

import (
	"serviceNest/interfaces"
	"serviceNest/model"
)

type AdminService struct {
	serviceRepo     interfaces.ServiceRepository
	userRepo        interfaces.UserRepository
	householderRepo interfaces.HouseholderRepository
	//serviceAreaRepo *repository_test.ServiceAreaRepository
	providerRepo       interfaces.ServiceProviderRepository
	serviceRequestRepo interfaces.ServiceRequestRepository
}

func NewAdminService(serviceRepo interfaces.ServiceRepository, serviceRequestRepo interfaces.ServiceRequestRepository, userRepo interfaces.UserRepository, providerRepo interfaces.ServiceProviderRepository) interfaces.AdminService {
	return &AdminService{
		serviceRepo: serviceRepo,
		userRepo:    userRepo,
		//serviceAreaRepo: serviceAreaRepo,
		providerRepo:       providerRepo,
		serviceRequestRepo: serviceRequestRepo,
	}
}

// View reports
func (s *AdminService) ViewReports() ([]model.ServiceRequest, error) {

	return s.serviceRequestRepo.GetAllServiceRequests()

}
func (s *AdminService) DeleteService(serviceID string) error {
	return s.serviceRepo.RemoveService(serviceID)
}

// Deactivate account
func (s *AdminService) DeactivateAccount(userID string) error {
	provider, err := s.providerRepo.GetProviderByID(userID)
	if err != nil {
		return err
	}

	provider.IsActive = false
	err = s.providerRepo.UpdateServiceProvider(provider)
	if err != nil {
		return err
	}
	return s.userRepo.DeActivateUser(userID)
}

func (s *AdminService) GetAllService() ([]model.Service, error) {
	return s.serviceRepo.GetAllServices()
}
