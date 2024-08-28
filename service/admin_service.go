package service

import (
	"serviceNest/interfaces"
	"serviceNest/model"
)

type AdminService struct {
	serviceRepo interfaces.ServiceRepository
	userRepo    interfaces.UserRepository
	//serviceAreaRepo *repository_test.ServiceAreaRepository
	providerRepo       interfaces.ServiceProviderRepository
	serviceRequestRepo interfaces.ServiceRequestRepository
}

func NewAdminService(serviceRepo interfaces.ServiceRepository, serviceRequestRepo interfaces.ServiceRequestRepository, userRepo interfaces.UserRepository, providerRepo interfaces.ServiceProviderRepository) *AdminService {
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
	return s.providerRepo.UpdateServiceProvider(provider)
}

func (s *AdminService) GetAllService() ([]*model.Service, error) {
	return s.serviceRepo.GetAllServices()
}
