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
//func TestServiceProviderService_AddService(t *testing.T) {
//	providerRepo := repository.NewServiceProviderRepository("test_providers.json")
//	serviceRepo := repository.NewServiceRepository("test_services.json")
//	serviceRequestRepo := repository.NewServiceRequestRepository("test_service_requests.json")
//	serviceProviderService := service2.NewServiceProviderService(providerRepo, serviceRequestRepo, serviceRepo)
//
//	service := model.Service{
//		ID:          "1",
//		Name:        "Plumbing",
//		Description: "Fixing pipes",
//	}
//
//	err := serviceProviderService.AddService("providerID", service)
//	assert.NoError(t, err, "Expected no error when adding service")
//
//	savedService, err := serviceRepo.GetServiceByID(service.ID)
//	assert.NoError(t, err, "Expected no error when fetching service")
//	assert.Equal(t, service.Name, savedService.Name, "Expected service name to match")
//}
//
//func TestServiceProviderService_RemoveService(t *testing.T) {
//	providerRepo := repository.NewServiceProviderRepository("test_providers.json")
//	serviceRepo := repository.NewServiceRepository("test_services.json")
//	serviceRequestRepo := repository.NewServiceRequestRepository("test_service_requests.json")
//	serviceProviderService := service2.NewServiceProviderService(providerRepo, serviceRequestRepo, serviceRepo)
//
//	err := serviceProviderService.RemoveService("providerID", "serviceID")
//	assert.NoError(t, err, "Expected no error when removing service")
//
//	_, err = serviceRepo.GetServiceByID("serviceID")
//	assert.Error(t, err, "Expected an error when fetching removed service")
//}
