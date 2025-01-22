package service_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"serviceNest/model"
	"serviceNest/service"
	"serviceNest/tests/mocks"
	"testing"
)

func TestViewReports(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)

	adminService := service.NewAdminService(mockServiceRepo, mockServiceRequestRepo, mockUserRepo, mockProviderRepo)

	serviceRequests := []model.ServiceRequest{
		{ID: "request1"},
		{ID: "request2"},
	}

	mockServiceRequestRepo.EXPECT().
		GetAllServiceRequests().
		Return(serviceRequests, nil)

	result, err := adminService.ViewReports()
	assert.NoError(t, err)
	assert.Equal(t, serviceRequests, result)
}
func TestDeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)

	adminService := service.NewAdminService(mockServiceRepo, mockServiceRequestRepo, mockUserRepo, mockProviderRepo)

	serviceID := "service1"

	mockServiceRepo.EXPECT().
		RemoveService(serviceID).
		Return(nil)

	err := adminService.DeleteService(serviceID)
	assert.NoError(t, err)
}
func TestDeactivateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)

	adminService := service.NewAdminService(mockServiceRepo, mockServiceRequestRepo, mockUserRepo, mockProviderRepo)

	userID := "provider1"
	provider := &model.ServiceProvider{
		User:     model.User{ID: userID},
		IsActive: true,
	}

	mockProviderRepo.EXPECT().
		GetProviderByID(userID).
		Return(provider, nil)

	mockProviderRepo.EXPECT().
		UpdateServiceProvider(provider).
		Do(func(p *model.ServiceProvider) {
			assert.False(t, p.IsActive)
		}).
		Return(nil)
	mockUserRepo.EXPECT().DeActivateUser(userID).Return(nil)
	err := adminService.DeactivateAccount(userID)
	assert.NoError(t, err)
}
func TestGetAllService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)

	adminService := service.NewAdminService(mockServiceRepo, mockServiceRequestRepo, mockUserRepo, mockProviderRepo)

	services := []model.Service{
		{ID: "service1", Name: "Service 1"},
		{ID: "service2", Name: "Service 2"},
	}

	mockServiceRepo.EXPECT().
		GetAllServices().
		Return(services, nil)

	result, err := adminService.GetAllService()
	assert.NoError(t, err)
	assert.Equal(t, services, result)
}
