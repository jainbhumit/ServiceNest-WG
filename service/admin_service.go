package service

import (
	"errors"
	"serviceNest/interfaces"
	"serviceNest/model"
	"serviceNest/util"
)

type AdminService struct {
	serviceRepo        interfaces.ServiceRepository
	userRepo           interfaces.UserRepository
	householderRepo    interfaces.HouseholderRepository
	providerRepo       interfaces.ServiceProviderRepository
	serviceRequestRepo interfaces.ServiceRequestRepository
}

func NewAdminService(serviceRepo interfaces.ServiceRepository, serviceRequestRepo interfaces.ServiceRequestRepository, userRepo interfaces.UserRepository, providerRepo interfaces.ServiceProviderRepository) interfaces.AdminService {
	return &AdminService{
		serviceRepo:        serviceRepo,
		userRepo:           userRepo,
		providerRepo:       providerRepo,
		serviceRequestRepo: serviceRequestRepo,
	}
}

func (s *AdminService) ViewReports(limit, offset int) ([]model.ServiceRequest, error) {

	return s.serviceRequestRepo.GetAllServiceRequests(limit, offset)

}
func (s *AdminService) DeleteService(serviceID string) error {
	return s.serviceRepo.RemoveService(serviceID)
}

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
	err = s.userRepo.DeActivateUser(userID)
	if err != nil {
		return err
	}

	// Delete all services associated with the user
	err = s.providerRepo.DeleteServicesByProviderID(userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminService) GetAllService(limit, offset int) ([]model.Service, error) {
	return s.serviceRepo.GetAllServices(limit, offset)
}

func (s *AdminService) AddService(name, description string) error {
	service := &model.Category{
		Name:        name,
		Description: description,
		ID:          util.GenerateUniqueID(),
	}
	err := s.serviceRepo.AddCategory(service)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminService) GetUserByEmail(email string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.Role != "Householder" {
		return nil, errors.New("You are not a Householder")
	}
	return user, nil
}
