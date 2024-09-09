package service_test

import (
	"errors"
	"fmt"
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
	serviceRequest := &model.ServiceRequest{ // Pointer type here
		ID:            requestID,
		HouseholderID: &householderID,
		Status:        "Accepted",
	}

	// Set up the mock expectations
	mockServiceRequestRepo.EXPECT().
		GetServiceRequestByID(requestID).
		Return(serviceRequest, nil) // Return pointer to serviceRequest

	// We use `Do` to verify the argument passed to `UpdateServiceRequest`
	mockServiceRequestRepo.EXPECT().
		UpdateServiceRequest(gomock.Any()).             // Accept any pointer argument here
		Do(func(updatedRequest *model.ServiceRequest) { // Expect pointer type
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

	// Create mocks
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Create the service object
	householderService := service.NewHouseholderService(nil, mockProviderRepo, mockServiceRepo, nil)

	// Test data
	services := []model.Service{
		{
			ID:         "service1",
			Category:   "Cleaning",
			ProviderID: "provider1",
		},
		{
			ID:         "service2",
			Category:   "Plumbing",
			ProviderID: "provider2",
		},
	}

	provider1 := &model.ServiceProviderDetails{
		ServiceProviderID: "provider1",
		Name:              "John Doe",
		Contact:           "123456789",
		Address:           "123 Street",
		Rating:            4.5,
	}

	// Mock behavior
	mockServiceRepo.EXPECT().
		GetAllServices().
		Return(services, nil)

	mockProviderRepo.EXPECT().
		GetProviderDetailByID("provider1").
		Return(provider1, nil)

	// Define the category to filter by
	category := "Cleaning"

	// Call the function
	filteredServices, err := householderService.GetServicesByCategory(category)

	// Assert no error
	assert.NoError(t, err)

	// Assert the filtered services
	expectedServices := []model.Service{
		{
			ID:              "service1",
			Category:        "Cleaning",
			ProviderID:      "provider1",
			ProviderName:    "John Doe",
			ProviderContact: "123456789",
			ProviderAddress: "123 Street",
		},
	}
	assert.Equal(t, expectedServices, filteredServices)
}

func TestGetServicesByCategory_NoMatchingCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Create the service object
	householderService := service.NewHouseholderService(nil, mockProviderRepo, mockServiceRepo, nil)

	// Test data
	services := []model.Service{
		{
			ID:         "service1",
			Category:   "Cleaning",
			ProviderID: "provider1",
		},
		{
			ID:         "service2",
			Category:   "Plumbing",
			ProviderID: "provider2",
		},
	}

	// Mock behavior
	mockServiceRepo.EXPECT().
		GetAllServices().
		Return(services, nil)

	// Define a category that does not match any service
	category := "Electrical"

	// Call the function
	filteredServices, err := householderService.GetServicesByCategory(category)

	// Assert no error
	assert.NoError(t, err)

	// Assert the filtered services are empty
	assert.Empty(t, filteredServices)
}

func TestGetServicesByCategory_ErrorFetchingProviderDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Create the service object
	householderService := service.NewHouseholderService(nil, mockProviderRepo, mockServiceRepo, nil)

	// Test data
	services := []model.Service{
		{
			ID:         "service1",
			Category:   "Cleaning",
			ProviderID: "provider1",
		},
	}

	// Mock behavior for GetAllServices
	mockServiceRepo.EXPECT().
		GetAllServices().
		Return(services, nil)

	// Simulate an error when fetching provider details
	// Updated to return *model.ServiceProviderDetails with error
	mockProviderRepo.EXPECT().
		GetProviderDetailByID("provider1").
		Return(&model.ServiceProviderDetails{}, errors.New("provider not found"))

	// Define the category to filter by
	category := "Cleaning"

	// Call the function
	_, err := householderService.GetServicesByCategory(category)

	// Assert that the error occurred
	assert.Error(t, err)
	assert.EqualError(t, err, "provider not found")
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
		Do(func(req *model.ServiceRequest) {
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
		Do(func(req *model.ServiceRequest) {
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

func TestApproveServiceRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	requestID := "request123"
	providerID := "provider123"

	// Define test cases
	tests := []struct {
		name                 string
		serviceRequest       *model.ServiceRequest
		repoErr              error
		providerErr          error
		expectedErr          error
		expectUpdateProvider bool
	}{
		{
			name: "Successful approval",
			serviceRequest: &model.ServiceRequest{
				ID:            requestID,
				ApproveStatus: false,
				ProviderDetails: []model.ServiceProviderDetails{
					{ServiceProviderID: providerID, Approve: false},
				},
			},
			repoErr:              nil,
			providerErr:          nil,
			expectedErr:          nil,
			expectUpdateProvider: true,
		},
		{
			name: "Service request already approved",
			serviceRequest: &model.ServiceRequest{
				ID:            requestID,
				ApproveStatus: true,
				ProviderDetails: []model.ServiceProviderDetails{
					{ServiceProviderID: providerID, Approve: false},
				},
			},
			repoErr:              nil,
			providerErr:          nil,
			expectedErr:          errors.New("service request has already been approved"),
			expectUpdateProvider: false,
		},
		{
			name:                 "Error retrieving service request",
			serviceRequest:       nil,
			repoErr:              fmt.Errorf("database error"),
			providerErr:          nil,
			expectedErr:          fmt.Errorf("could not find service request: %v", fmt.Errorf("database error")),
			expectUpdateProvider: false,
		},
		{
			name: "Error updating provider detail",
			serviceRequest: &model.ServiceRequest{
				ID:            requestID,
				ApproveStatus: false,
				ProviderDetails: []model.ServiceProviderDetails{
					{ServiceProviderID: providerID, Approve: false},
				},
			},
			repoErr:              nil,
			providerErr:          fmt.Errorf("provider update error"),
			expectedErr:          fmt.Errorf("could not update service provider detail"),
			expectUpdateProvider: true,
		},
		{
			name: "Error updating service request",
			serviceRequest: &model.ServiceRequest{
				ID:            requestID,
				ApproveStatus: false,
				ProviderDetails: []model.ServiceProviderDetails{
					{ServiceProviderID: providerID, Approve: true},
				},
			},
			repoErr:              fmt.Errorf("service request update error"),
			providerErr:          nil,
			expectedErr:          fmt.Errorf("could not find service request: %v", fmt.Errorf("service request update error")),
			expectUpdateProvider: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Expectation for retrieving service provider details by request ID
			mockServiceRequestRepo.EXPECT().
				GetServiceProviderByRequestID(requestID, providerID).
				Return(tt.serviceRequest, tt.repoErr).
				Times(1)

			if tt.repoErr == nil && tt.expectUpdateProvider {
				// Expectation for updating provider details
				mockProviderRepo.EXPECT().
					UpdateServiceProviderDetailByRequestID(gomock.Any(), requestID).
					Return(tt.providerErr).
					Times(1)
			}

			// Expectation for updating the service request if no repo error occurred
			if tt.repoErr == nil && tt.providerErr == nil && tt.expectUpdateProvider {
				mockServiceRequestRepo.EXPECT().
					UpdateServiceRequest(gomock.Any()).
					Return(tt.repoErr).
					Times(1)
			}

			// Call the function under test
			err := service.ApproveServiceRequest(requestID, providerID)

			// Assert the result
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAddReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	var originalGenerateUniqueID = service.GetUniqueID
	defer func() { service.GetUniqueID = originalGenerateUniqueID }()
	service.GetUniqueID = func() string {
		return "uniqueID"
	}

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	// Replace util.GenerateUniqueID with a mockable function if necessary

	providerID := "provider123"
	householderID := "householder123"
	serviceID := "service123"
	comments := "Great service!"
	rating := 4.5

	// Define test cases
	tests := []struct {
		name            string
		addReviewErr    error
		updateRatingErr error
		expectedErr     error
	}{
		{
			name:            "Successful review addition and rating update",
			addReviewErr:    nil,
			updateRatingErr: nil,
			expectedErr:     nil,
		},
		{
			name:            "Error adding review",
			addReviewErr:    errors.New("add review error"),
			updateRatingErr: nil,
			expectedErr:     errors.New("add review error"),
		},
		{
			name:            "Error updating provider rating",
			addReviewErr:    nil,
			updateRatingErr: errors.New("update rating error"),
			expectedErr:     errors.New("failed to update provider rating"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up expectations
			mockProviderRepo.EXPECT().
				AddReview(gomock.Any()). // gomock.Any() is used to match any Review object
				Return(tt.addReviewErr).
				Times(1)

			if tt.addReviewErr == nil {
				mockProviderRepo.EXPECT().
					UpdateProviderRating(providerID).
					Return(tt.updateRatingErr).
					Times(1)
			}

			// Call the method under test
			err := service.AddReview(providerID, householderID, serviceID, comments, rating)

			// Assert results
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestRequestService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	// Mock unique ID generation
	originalGenerateUniqueID := service.GetUniqueID
	service.GetUniqueID = func() string {
		return "uniqueID"
	}
	defer func() { service.GetUniqueID = originalGenerateUniqueID }()

	householderService := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	householder := &model.Householder{
		User: model.User{
			ID:      "householderID",
			Name:    "John Doe",
			Address: "123 Main St",
		},
	}

	serviceName := "Test Service"
	scheduledTime := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name              string
		getServiceErr     error
		service           *model.Service
		saveServiceErr    error
		saveRequestErr    error
		expectedRequestID string
		expectedErr       error
	}{

		{
			name:              "Service exists, successful request creation",
			getServiceErr:     nil,
			service:           &model.Service{ID: "existingServiceID"},
			saveServiceErr:    nil,
			saveRequestErr:    nil,
			expectedRequestID: "uniqueID",
			expectedErr:       nil,
		},
		{
			name:              "Error fetching service",
			getServiceErr:     errors.New("service fetch error"),
			service:           nil,
			saveServiceErr:    nil,
			saveRequestErr:    nil,
			expectedRequestID: "",
			expectedErr:       errors.New("service fetch error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the call to GetServiceByName based on the test case
			mockServiceRepo.EXPECT().
				GetServiceByName(serviceName).
				Return(tt.service, tt.getServiceErr).
				Times(1)

			// If there's an error fetching the service, SaveService and SaveServiceRequest should NOT be called
			if tt.getServiceErr != nil {
				mockServiceRepo.EXPECT().
					SaveService(gomock.Any()).
					Times(0)

				mockServiceRequestRepo.EXPECT().
					SaveServiceRequest(gomock.Any()).
					Times(0)
			} else {
				// If the service is not found, expect SaveService to be called
				if tt.service == nil {
					mockServiceRepo.EXPECT().
						SaveService(gomock.Any()).
						Return(tt.saveServiceErr).
						Times(1)
				} else {
					// Ensure SaveService is not called if the service exists
					mockServiceRepo.EXPECT().
						SaveService(gomock.Any()).
						Times(0)
				}

				// Mock the call to SaveServiceRequest, which should always be called
				mockServiceRequestRepo.EXPECT().
					SaveServiceRequest(gomock.Any()).
					Return(tt.saveRequestErr).
					Times(1)
			}

			// Call the method under test
			requestID, err := householderService.RequestService(householder, serviceName, &scheduledTime)

			// Assert results
			assert.Equal(t, tt.expectedRequestID, requestID)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetAvailableServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	service := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	services := []model.Service{
		{
			Name:        "service1",
			Description: "plumber",
			Price:       500,
		},
	}
	mockServiceRepo.EXPECT().GetAllServices().Return(services, nil)

	availableServices, err := service.GetAvailableServices()
	assert.Equal(t, services, availableServices)
	assert.Nil(t, err)
}
func TestRequestService_ALL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHouseholderRepo := mocks.NewMockHouseholderRepository(ctrl)
	mockProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	// Mock unique ID generation
	originalGenerateUniqueID := service.GetUniqueID
	service.GetUniqueID = func() string {
		return "uniqueID"
	}
	defer func() { service.GetUniqueID = originalGenerateUniqueID }()

	householderService := service.NewHouseholderService(mockHouseholderRepo, mockProviderRepo, mockServiceRepo, mockServiceRequestRepo)

	householder := &model.Householder{
		User: model.User{
			ID:      "householderID",
			Name:    "John Doe",
			Address: "123 Main St",
		},
	}

	serviceName := "Test Service"
	scheduledTime := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name              string
		getServiceErr     error
		service           *model.Service
		saveServiceErr    error
		saveRequestErr    error
		expectedRequestID string
		expectedErr       error
	}{
		{
			name:              "Service exists, successful request creation",
			getServiceErr:     nil,
			service:           &model.Service{ID: "existingServiceID"},
			saveServiceErr:    nil,
			saveRequestErr:    nil,
			expectedRequestID: "uniqueID",
			expectedErr:       nil,
		},
		{
			name:              "Error fetching service",
			getServiceErr:     errors.New("service fetch error"),
			service:           nil,
			saveServiceErr:    nil,
			saveRequestErr:    nil,
			expectedRequestID: "",
			expectedErr:       errors.New("service fetch error"),
		},
		{
			name:              "Service does not exist, successful custom service creation",
			getServiceErr:     nil,
			service:           nil,
			saveServiceErr:    nil,
			saveRequestErr:    nil,
			expectedRequestID: "uniqueID",
			expectedErr:       nil,
		},
		{
			name:              "Service does not exist, error creating custom service",
			getServiceErr:     nil,
			service:           nil,
			saveServiceErr:    errors.New("save service error"),
			saveRequestErr:    nil,
			expectedRequestID: "",
			expectedErr:       errors.New("save service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the call to GetServiceByName based on the test case
			mockServiceRepo.EXPECT().
				GetServiceByName(serviceName).
				Return(tt.service, tt.getServiceErr).
				Times(1)

			// If there's an error fetching the service, SaveService and SaveServiceRequest should NOT be called
			if tt.getServiceErr != nil {
				mockServiceRepo.EXPECT().
					SaveService(gomock.Any()).
					Times(0)

				mockServiceRequestRepo.EXPECT().
					SaveServiceRequest(gomock.Any()).
					Times(0)
			} else {
				// If the service is not found, expect SaveService to be called
				if tt.service == nil {
					mockServiceRepo.EXPECT().
						SaveService(gomock.Any()).
						Return(tt.saveServiceErr).
						Times(1)
					// If there's an error saving the custom service, do not expect SaveServiceRequest to be called
					if tt.saveServiceErr != nil {
						mockServiceRequestRepo.EXPECT().
							SaveServiceRequest(gomock.Any()).
							Times(0)
					}
				} else {
					// Ensure SaveService is not called if the service exists
					mockServiceRepo.EXPECT().
						SaveService(gomock.Any()).
						Times(0)
				}

				// Mock the call to SaveServiceRequest
				if tt.saveServiceErr == nil {
					mockServiceRequestRepo.EXPECT().
						SaveServiceRequest(gomock.Any()).
						Return(tt.saveRequestErr).
						Times(1)
				}
			}

			// Call the method under test
			requestID, err := householderService.RequestService(householder, serviceName, &scheduledTime)

			// Assert results
			assert.Equal(t, tt.expectedRequestID, requestID)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
