package service_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"serviceNest/model"
	"serviceNest/service"
	"serviceNest/tests/mocks"
	"time"

	"testing"
)

func TestViewStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	householder := &model.Householder{User: model.User{ID: "householder1"}}
	requests := []model.ServiceRequest{
		{ID: "request1", HouseholderID: &householder.ID, Status: "Pending"},
	}

	mockServiceRequestRepo.EXPECT().
		GetServiceRequestsByHouseholderID(householder.ID).
		Return(requests, nil)

	result, err := service.ViewStatus(service, householder)
	assert.NoError(t, err)
	assert.Equal(t, requests, result)
}

func TestCancelAcceptedRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	requestID := "request1"
	householderID := "householder1"
	serviceRequest := &model.ServiceRequest{
		ID:            requestID,
		HouseholderID: &householderID,
		Status:        "Accepted",
	}

	// Set up the mock expectations
	mockServiceRequestRepo.EXPECT().
		GetServiceRequestByID(requestID).
		Return(serviceRequest, nil)

	// We use `Do` to verify the argument passed to `UpdateServiceRequest`
	mockServiceRequestRepo.EXPECT().
		UpdateServiceRequest(gomock.Any()). // Accept any argument here
		Do(func(updatedRequest model.ServiceRequest) {
			assert.Equal(t, "Cancelled", updatedRequest.Status)
			assert.Equal(t, requestID, updatedRequest.ID)
			assert.Equal(t, householderID, *updatedRequest.HouseholderID)
		}).
		Return(nil)

	err := service.CancelAcceptedRequest(requestID, householderID)
	assert.NoError(t, err)
}

func TestSearchService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	householder := &model.Householder{User: model.User{ID: "householder1", Latitude: 10, Longitude: 10}}
	providers := []model.ServiceProvider{
		{User: model.User{ID: "provider1", Latitude: 10, Longitude: 10}, ServicesOffered: []model.Service{{ID: "service1", Name: "Cleaning"}}},
	}

	mockProviderRepo.EXPECT().
		GetProvidersByServiceType("Cleaning").
		Return(providers, nil)

	result, err := service.SearchService(householder, "Cleaning")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "provider1", result[0].ID)
}

func TestGetServicesByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	category := "Home"
	services := []*model.Service{
		{ID: "service1", Category: category},
		{ID: "service2", Category: "Office"},
	}

	mockServiceRepo.EXPECT().
		GetAllServices().
		Return(services, nil)

	result, err := service.GetServicesByCategory(category)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "service1", result[0].ID)
}

func TestRequestService_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	householder := &model.Householder{
		User: model.User{
			ID:      "householder1",
			Name:    "John Doe",
			Address: "123 Main St",
		},
	}
	serviceID := "service1"
	scheduleTime := time.Now()

	// Expected error to be returned by the repository mock
	expectedError := errors.New("failed to save service request")

	// Adjusting the service request to align with the original function
	mockServiceRequestRepo.EXPECT().
		SaveServiceRequest(gomock.Any()).
		Return(expectedError)

	// Call the RequestService method
	resultID, err := service.RequestService(householder, serviceID, &scheduleTime)

	// Asserting that the error is returned and the result ID is empty
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, resultID)
}

func TestViewBookingHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	householderID := "householder1"
	requests := []model.ServiceRequest{
		{ID: "request1", HouseholderID: &householderID},
	}

	mockServiceRequestRepo.EXPECT().
		GetServiceRequestsByHouseholderID(householderID).
		Return(requests, nil)

	result, err := service.ViewBookingHistory(householderID)
	assert.NoError(t, err)
	assert.Equal(t, requests, result)
}

func TestCancelServiceRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	requestID := "request1"
	serviceRequest := &model.ServiceRequest{
		ID:     requestID,
		Status: "Pending",
	}

	mockServiceRequestRepo.EXPECT().
		GetServiceRequestByID(requestID).
		Return(serviceRequest, nil)
	mockServiceRequestRepo.EXPECT().
		UpdateServiceRequest(gomock.Any()).
		Do(func(req model.ServiceRequest) {
			// Ensure that the Status is updated correctly.
			assert.Equal(t, "Cancelled", req.Status)
		}).
		Return(nil)

	err := service.CancelServiceRequest(requestID)
	assert.NoError(t, err)
	assert.Equal(t, "Cancelled", serviceRequest.Status)
}

func TestRescheduleServiceRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	requestID := "request1"
	newTime := time.Now().Add(time.Hour * 24)
	serviceRequest := &model.ServiceRequest{
		ID:            requestID,
		Status:        "Pending",
		ScheduledTime: time.Now(),
	}

	mockServiceRequestRepo.EXPECT().
		GetServiceRequestByID(requestID).
		Return(serviceRequest, nil)
	mockServiceRequestRepo.EXPECT().
		UpdateServiceRequest(gomock.Any()).
		Do(func(req model.ServiceRequest) {
			// Ensure that the ScheduledTime is updated correctly.
			assert.Equal(t, newTime, req.ScheduledTime)
		}).
		Return(nil)

	err := service.RescheduleServiceRequest(requestID, newTime)
	assert.NoError(t, err)
	assert.Equal(t, newTime, serviceRequest.ScheduledTime)
}

func TestViewServiceRequestStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	requestID := "request1"
	status := "Accepted"
	serviceRequest := &model.ServiceRequest{
		ID:     requestID,
		Status: status,
	}

	mockServiceRequestRepo.EXPECT().
		GetServiceRequestByID(requestID).
		Return(serviceRequest, nil)

	result, err := service.ViewServiceRequestStatus(requestID)
	assert.NoError(t, err)
	assert.Equal(t, status, result)
}
func TestViewApprovedRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	householderID := "householder1"

	serviceRequests := []model.ServiceRequest{
		{
			ID:            "request1",
			HouseholderID: &householderID,
			ApproveStatus: true,
		},
		{
			ID:            "request2",
			HouseholderID: &householderID,
			ApproveStatus: false,
		},
	}

	mockServiceRequestRepo.EXPECT().
		GetServiceRequestsByHouseholderID(householderID).
		Return(serviceRequests, nil)

	approvedRequests, err := service.ViewApprovedRequests(householderID)
	assert.NoError(t, err)
	assert.Len(t, approvedRequests, 1)
	assert.Equal(t, "request1", approvedRequests[0].ID)
}

func TestViewApprovedRequests_NoApprovedRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	householderID := "householder1"

	serviceRequests := []model.ServiceRequest{
		{
			ID:            "request1",
			HouseholderID: &householderID,
			ApproveStatus: false,
		},
	}

	mockServiceRequestRepo.EXPECT().
		GetServiceRequestsByHouseholderID(householderID).
		Return(serviceRequests, nil)

	approvedRequests, err := service.ViewApprovedRequests(householderID)
	assert.Error(t, err)
	assert.Nil(t, approvedRequests)
	assert.EqualError(t, err, "no approved service requests found")
}
