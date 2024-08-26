package service

//
//import (
//	"github.com/stretchr/testify/assert"
//	"serviceNest/model"
//	"serviceNest/repository"
//	service2 "serviceNest/service"
//	"testing"
//)
//
//func TestAdminService_ManageServices_Add(t *testing.T) {
//	adminRepo := repository.NewAdminRepository("test_admins.json")
//	serviceRepo := repository.NewServiceRepository("test_services.json")
//	serviceRequestRepo := repository.NewServiceRequestRepository("test_service_requests.json")
//	userRepo := repository.NewUserRepository("test_users.json")
//	providerRepo := repository.NewServiceProviderRepository("test_providers.json")
//
//	err1 := service2.NewAdminService(adminRepo, serviceRepo, serviceRequestRepo, userRepo, providerRepo)
//	if err1 != nil {
//		t.Error(err1)
//	}
//	service := model.Service{
//		ID:          "1",
//		Name:        "Plumbing",
//		Description: "Fixing pipes",
//	}
//
//	savedService, err := serviceRepo.GetServiceByID(service.ID)
//	assert.NoError(t, err, "Expected no error when fetching service")
//	assert.Equal(t, service.Name, savedService.Name, "Expected service name to match")
//}
//
//func TestAdminService_ViewReports(t *testing.T) {
//	adminRepo := repository.NewAdminRepository("test_admins.json")
//	serviceRepo := repository.NewServiceRepository("test_services.json")
//	serviceRequestRepo := repository.NewServiceRequestRepository("test_service_requests.json")
//	userRepo := repository.NewUserRepository("test_users.json")
//	providerRepo := repository.NewServiceProviderRepository("test_providers.json")
//
//	adminService := service2.NewAdminService(adminRepo, serviceRepo, serviceRequestRepo, userRepo, providerRepo)
//
//	reports, err := adminService.ViewReports()
//	assert.NoError(t, err, "Expected no error when viewing reports")
//	assert.Greater(t, len(reports), 0, "Expected at least one report")
//}
