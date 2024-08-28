package service_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"serviceNest/model"
	"serviceNest/service"
	"serviceNest/tests/mocks"
	"testing"
)

func TestAddService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider1"
	newService := model.Service{ID: "service1", Name: "Test Service"}

	mockServiceProviderRepo.EXPECT().
		GetProviderByID(providerID).
		Return(&model.ServiceProvider{User: model.User{ID: providerID}, ServicesOffered: []model.Service{}}, nil)
	mockServiceProviderRepo.EXPECT().
		UpdateServiceProvider(gomock.Any()).
		Return(nil)
	mockServiceRepo.EXPECT().
		SaveService(newService).
		Return(nil)

	err := serviceProviderService.AddService(providerID, newService)
	assert.NoError(t, err)
}

func TestUpdateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider1"
	serviceID := "service1"
	updatedService := model.Service{ID: serviceID, Name: "Updated Service"}

	mockServiceProviderRepo.EXPECT().
		GetProviderByID(providerID).
		Return(&model.ServiceProvider{
			User: model.User{ID: providerID},
			ServicesOffered: []model.Service{
				{ID: serviceID, Name: "Old Service"},
			},
		}, nil)
	mockServiceProviderRepo.EXPECT().
		UpdateServiceProvider(gomock.Any()).
		Return(nil)
	mockServiceRepo.EXPECT().
		SaveService(updatedService).
		Return(nil)

	err := serviceProviderService.UpdateService(providerID, serviceID, updatedService)
	assert.NoError(t, err)
}

func TestRemoveService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider1"
	serviceID := "service1"

	mockServiceProviderRepo.EXPECT().
		GetProviderByID(providerID).
		Return(&model.ServiceProvider{
			User: model.User{ID: providerID},
			ServicesOffered: []model.Service{
				{ID: serviceID, Name: "Test Service"},
			},
		}, nil)
	mockServiceProviderRepo.EXPECT().
		UpdateServiceProvider(gomock.Any()).
		Return(nil)
	mockServiceRepo.EXPECT().
		RemoveService(serviceID).
		Return(nil)

	err := serviceProviderService.RemoveService(providerID, serviceID)
	assert.NoError(t, err)
}

func TestAddReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	serviceID := "service1"
	householderID := "householder1"
	comments := "Great Service!"
	rating := 4.5

	mockServiceProviderRepo.EXPECT().
		GetProviderByID(serviceID).
		Return(&model.ServiceProvider{ServicesOffered: []model.Service{{ID: serviceID}}}, nil)
	mockServiceProviderRepo.EXPECT().
		UpdateServiceProvider(gomock.Any()).
		Return(nil)

	err := serviceProviderService.AddReview(serviceID, householderID, comments, rating)
	assert.NoError(t, err)
}

func TestAcceptServiceRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider1"
	requestID := "request1"

	mockServiceRequestRepo.EXPECT().
		GetServiceRequestByID(requestID).
		Return(&model.ServiceRequest{
			ID:     requestID,
			Status: "Pending",
		}, nil)
	mockServiceProviderRepo.EXPECT().
		GetProviderByID(providerID).
		Return(&model.ServiceProvider{
			User: model.User{ID: providerID},
		}, nil)
	mockServiceRequestRepo.EXPECT().
		UpdateServiceRequest(gomock.Any()).
		Return(nil)

	err := serviceProviderService.AcceptServiceRequest(providerID, requestID)
	assert.NoError(t, err)
}

func TestDeclineServiceRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider1"
	requestID := "request1"

	mockServiceRequestRepo.EXPECT().
		GetServiceRequestByID(requestID).
		Return(&model.ServiceRequest{
			ID:     requestID,
			Status: "Pending",
		}, nil)
	mockServiceRequestRepo.EXPECT().
		UpdateServiceRequest(gomock.Any()).
		Return(nil)

	err := serviceProviderService.DeclineServiceRequest(providerID, requestID)
	assert.NoError(t, err)
}

func TestUpdateAvailability(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider1"
	availability := true

	mockServiceProviderRepo.EXPECT().
		GetProviderByID(providerID).
		Return(&model.ServiceProvider{User: model.User{ID: providerID}}, nil)
	mockServiceProviderRepo.EXPECT().
		UpdateServiceProvider(gomock.Any()).
		Return(nil)

	err := serviceProviderService.UpdateAvailability(providerID, availability)
	assert.NoError(t, err)
}

func TestViewServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider1"
	services := []model.Service{
		{ID: "service1", Name: "Test Service"},
	}

	mockServiceProviderRepo.EXPECT().
		GetProviderByID(providerID).
		Return(&model.ServiceProvider{
			User:            model.User{ID: providerID},
			ServicesOffered: services,
		}, nil)

	result, err := serviceProviderService.ViewServices(providerID)
	assert.NoError(t, err)
	assert.Equal(t, services, result)
}

func TestGetServiceByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	serviceID := "123"
	expectedService := &model.Service{ID: serviceID, Name: "Service Name"}

	mockServiceRepo.EXPECT().
		GetServiceByID(serviceID).
		Return(expectedService, nil)

	result, err := serviceProviderService.GetServiceByID(serviceID)
	assert.NoError(t, err)
	assert.Equal(t, expectedService, result)
}

func TestViewReviews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider123"
	expectedReviews := []*model.Review{
		{ID: "review1", Comments: "Great service!"},
		{ID: "review2", Comments: "Not bad!"},
	}

	mockServiceProviderRepo.EXPECT().
		GetProviderByID(providerID).
		Return(&model.ServiceProvider{
			User:    model.User{ID: providerID},
			Reviews: expectedReviews,
		}, nil)

	result, err := serviceProviderService.ViewReviews(providerID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedReviews, result)
}
