package service

import (
	"serviceNest/model"
	"serviceNest/repository"
)

type AdminService struct {
	adminRepo   *repository.AdminRepository
	serviceRepo *repository.ServiceRepository
	userRepo    *repository.UserRepository
	//serviceAreaRepo *repository.ServiceAreaRepository
	providerRepo       *repository.ServiceProviderRepository
	serviceRequestRepo *repository.ServiceRequestRepository
}

func NewAdminService(adminRepo *repository.AdminRepository, serviceRepo *repository.ServiceRepository, serviceRequestRepo *repository.ServiceRequestRepository, userRepo *repository.UserRepository, providerRepo *repository.ServiceProviderRepository) *AdminService {
	return &AdminService{
		adminRepo:   adminRepo,
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

//// Add service area
//func (s *AdminService) AddServiceArea(area model.ServiceArea) error {
//	service, err := s.serviceRepo.GetServiceByID(serviceID)
//	if err != nil {
//		return err
//	}
//
//	service.Areas = append(service.Areas, area)
//	return s.serviceRepo.UpdateService(service)
//
//	return nil
//}
//
//// Remove service area
//func (s *AdminService) RemoveServiceArea(areaID string) error {
//	service, err := s.serviceRepo.GetServiceByID(serviceID)
//	if err != nil {
//		return err
//	}
//
//	// Filter out the area to remove it
//	var updatedAreas []string
//	for _, a := range service.Areas {
//		if a != area {
//			updatedAreas = append(updatedAreas, a)
//		}
//	}
//
//	service.Areas = updatedAreas
//	return s.serviceRepo.UpdateService(service)
//	return nil
//}
//
//// Update service area
//func (s *AdminService) UpdateServiceArea(areaID string, newArea model.ServiceArea) error {
//	service, err := s.serviceRepo.GetServiceByID(serviceID)
//	if err != nil {
//		return err
//	}
//
//	for i, a := range service.Areas {
//		if a == oldArea {
//			service.Areas[i] = newArea
//			break
//		}
//	}
//
//	return s.serviceRepo.UpdateService(service)
//	return nil
//}
