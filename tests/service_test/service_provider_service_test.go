package service_test

import (
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"serviceNest/model"
	"serviceNest/service"
	"serviceNest/tests/mocks"
	"testing"
	"time"
)

func convertNullString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
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

	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)

	providerID := "provider-123"
	serviceID := "service-456"
	updatedService := model.Service{ID: serviceID, Name: "Updated Service"}

	mockServiceRepo.EXPECT().UpdateService(providerID, updatedService).Return(nil)

	svc := service.NewServiceProviderService(nil, nil, mockServiceRepo)

	err := svc.UpdateService(providerID, serviceID, updatedService)
	assert.NoError(t, err)
}

func TestRemoveService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)

	providerID := "provider-123"
	serviceID := "service-456"

	mockServiceRepo.EXPECT().RemoveServiceByProviderID(providerID, serviceID).Return(nil)

	svc := service.NewServiceProviderService(nil, nil, mockServiceRepo)

	err := svc.RemoveService(providerID, serviceID)
	assert.NoError(t, err)
}

func TestAcceptServiceRequestWithMockInput(t *testing.T) {
	// Mock input to simulate typing "150" for the price
	input := "150\n"

	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Save the original stdin
	oldStdin := os.Stdin
	defer func() {
		os.Stdin = oldStdin
	}()

	// Set os.Stdin to read from the pipe (our simulated input)
	os.Stdin = r

	// Write the input to the pipe
	_, err = w.Write([]byte(input))
	if err != nil {
		t.Fatalf("Failed to write to pipe: %v", err)
	}
	// Close the writer to simulate end of input
	w.Close()

	// Set up mocks and other test structures
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)
	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)

	providerID := "provider-123"
	requestID := "request-456"

	mockServiceRequest := &model.ServiceRequest{
		ID:            requestID,
		Status:        "Pending",
		ApproveStatus: false,
	}

	mockProviderDetails := &model.ServiceProviderDetails{
		ServiceProviderID: providerID,
		Name:              "John's Services",
		Contact:           "1234567890",
		Address:           "123 Service Lane",
		Price:             "150", // This will be overwritten by Scanln mock
		Rating:            4.2,
		Reviews:           []model.Review{},
	}

	mockServiceRequestRepo.EXPECT().GetServiceRequestByID(requestID).Return(mockServiceRequest, nil)
	mockServiceProviderRepo.EXPECT().GetProviderDetailByID(providerID).Return(mockProviderDetails, nil)
	mockServiceProviderRepo.EXPECT().GetReviewsByProviderID(providerID).Return([]model.Review{}, nil)
	mockServiceRequestRepo.EXPECT().UpdateServiceRequest(mockServiceRequest).Return(nil)
	mockServiceProviderRepo.EXPECT().SaveServiceProviderDetail(mockProviderDetails, requestID).Return(nil)

	// Initialize the service with mock repositories
	svc := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, nil)

	// Call the method
	err = svc.AcceptServiceRequest(providerID, requestID)

	// Check for errors
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the price has been set correctly from the mocked input
	if mockServiceRequest.ProviderDetails[0].Price != "150" {
		t.Errorf("expected price to be '150', got %v", mockServiceRequest.ProviderDetails[0].Price)
	}
}

func TestDeclineServiceRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	providerID := "provider-123"
	requestID := "request-456"

	mockServiceRequest := &model.ServiceRequest{
		ID:     requestID,
		Status: "Pending",
	}

	mockServiceRequestRepo.EXPECT().GetServiceRequestByID(requestID).Return(mockServiceRequest, nil)
	mockServiceRequestRepo.EXPECT().UpdateServiceRequest(mockServiceRequest).Return(nil)

	svc := service.NewServiceProviderService(nil, mockServiceRequestRepo, nil)

	err := svc.DeclineServiceRequest(providerID, requestID)
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

	mockServiceRepo.EXPECT().
		GetServiceByProviderID(providerID).
		Return(services, nil)

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

func TestGetReviews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider123"
	expectedReviews := []model.Review{
		{ID: "review1", Comments: "Great service!"},
		{ID: "review2", Comments: "Not bad!"},
	}

	mockServiceProviderRepo.EXPECT().
		GetReviewsByProviderID(providerID).
		Return(expectedReviews, nil)

	result, err := serviceProviderService.GetReviews(providerID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedReviews, result)
}
func TestGetAllServiceRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)
	mockServiceRequest := []model.ServiceRequest{
		{ID: "requestID",
			Status: "Pending"},
	}
	mockServiceRequestRepo.EXPECT().GetAllServiceRequests().Return(mockServiceRequest, nil)

	result, err := serviceProviderService.GetAllServiceRequests()
	assert.NoError(t, err)
	assert.Equal(t, mockServiceRequest, result)

}
func TestViewApprovedRequestsByHouseholder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	// Create a sample service provider ID and mock service requests
	providerID := "provider-123"

	// Define mock data for service requests
	mockServiceRequests := []model.ServiceRequest{
		{
			ID:            "request-1",
			ApproveStatus: true, // Approved request
			ProviderDetails: []model.ServiceProviderDetails{
				{
					ServiceProviderID: providerID,
					Approve:           true, // Approved by this provider
				},
			},
		},
		{
			ID:            "request-2",
			ApproveStatus: false, // Not approved
			ProviderDetails: []model.ServiceProviderDetails{
				{
					ServiceProviderID: providerID,
					Approve:           false, // Not approved by this provider
				},
			},
		},
		{
			ID:            "request-3",
			ApproveStatus: true, // Approved request
			ProviderDetails: []model.ServiceProviderDetails{
				{
					ServiceProviderID: "other-provider", // Different provider
					Approve:           true,
				},
			},
		},
	}

	// Set up the mock to return the sample service requests
	mockServiceRequestRepo.EXPECT().GetServiceRequestsByProviderID(providerID).Return(mockServiceRequests, nil)

	// Initialize the ServiceProviderService with the mock repository
	svc := service.NewServiceProviderService(nil, mockServiceRequestRepo, nil)

	// Call the function to test
	approvedRequests, err := svc.ViewApprovedRequestsByHouseholder(providerID)

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// We expect only 1 request to be approved by the given provider
	expectedApprovedCount := 1
	if len(approvedRequests) != expectedApprovedCount {
		t.Errorf("expected %d approved requests, got %d", expectedApprovedCount, len(approvedRequests))
	}

	// Check if the correct request is returned
	if approvedRequests[0].ID != "request-1" {
		t.Errorf("expected request ID 'request-1', got %v", approvedRequests[0].ID)
	}
}

func TestViewApprovedRequestsByHouseholder_NoApprovedRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	// Create a sample service provider ID and mock service requests
	providerID := "provider-123"

	// Define mock data for service requests with no approved requests
	mockServiceRequests := []model.ServiceRequest{
		{
			ID:            "request-1",
			ApproveStatus: false,
			ProviderDetails: []model.ServiceProviderDetails{
				{
					ServiceProviderID: providerID,
					Approve:           false,
				},
			},
		},
		{
			ID:            "request-2",
			ApproveStatus: false,
			ProviderDetails: []model.ServiceProviderDetails{
				{
					ServiceProviderID: providerID,
					Approve:           false,
				},
			},
		},
	}

	// Set up the mock to return the sample service requests
	mockServiceRequestRepo.EXPECT().GetServiceRequestsByProviderID(providerID).Return(mockServiceRequests, nil)

	// Initialize the ServiceProviderService with the mock repository
	svc := service.NewServiceProviderService(nil, mockServiceRequestRepo, nil)

	// Call the function to test
	_, err := svc.ViewApprovedRequestsByHouseholder(providerID)

	// Assertions
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	expectedError := "no approved requests found for this provider"
	if err.Error() != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestViewApprovedRequestsByHouseholder_ErrorRetrievingRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	// Create a sample service provider ID
	providerID := "provider-123"

	// Set up the mock to return an error
	mockServiceRequestRepo.EXPECT().GetServiceRequestsByProviderID(providerID).Return(nil, errors.New("database error"))

	// Initialize the ServiceProviderService with the mock repository
	svc := service.NewServiceProviderService(nil, mockServiceRequestRepo, nil)

	// Call the function to test
	_, err := svc.ViewApprovedRequestsByHouseholder(providerID)

	// Assertions
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	expectedError := "could not retrieve service requests: database error"
	if err.Error() != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestViewServices_ErrorFromRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider1"
	expectedError := errors.New("database error")

	mockServiceRepo.EXPECT().
		GetServiceByProviderID(providerID).
		Return(nil, expectedError)

	result, err := serviceProviderService.ViewServices(providerID)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

func TestViewServices_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "provider1"
	services := []model.Service{} // Empty result

	mockServiceRepo.EXPECT().
		GetServiceByProviderID(providerID).
		Return(services, nil)

	result, err := serviceProviderService.ViewServices(providerID)
	assert.NoError(t, err)
	assert.Equal(t, services, result)
}

func TestViewServices_InvalidProviderID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderRepo := mocks.NewMockServiceProviderRepository(ctrl)
	mockServiceRepo := mocks.NewMockServiceRepository(ctrl)
	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)

	serviceProviderService := service.NewServiceProviderService(mockServiceProviderRepo, mockServiceRequestRepo, mockServiceRepo)

	providerID := "" // Invalid provider ID
	services := []model.Service{}

	mockServiceRepo.EXPECT().
		GetServiceByProviderID(providerID).
		Return(services, nil)

	result, err := serviceProviderService.ViewServices(providerID)
	assert.NoError(t, err)
	assert.Equal(t, services, result)
}
func TestGetServiceRequestByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceRequestRepo := mocks.NewMockServiceRequestRepository(ctrl)
	serviceProviderService := service.NewServiceProviderService(nil, mockServiceRequestRepo, nil)

	requestID := "request123"
	expectedRequest := &model.ServiceRequest{
		ID:                 requestID,
		HouseholderID:      convertNullString(sql.NullString{String: "1001", Valid: true}),
		HouseholderName:    "John Doe",
		HouseholderAddress: convertNullString(sql.NullString{String: "123 Main St", Valid: true}),
		ServiceID:          "service123",
		RequestedTime:      time.Now(),
		ScheduledTime:      time.Now().Add(24 * time.Hour),
		Status:             "Pending",
		ApproveStatus:      false,
	}

	tests := []struct {
		name          string
		mockGetByID   func()
		expectedReq   *model.ServiceRequest
		expectedError error
	}{
		{
			name: "Successful retrieval",
			mockGetByID: func() {
				mockServiceRequestRepo.EXPECT().
					GetServiceRequestByID(requestID).
					Return(expectedRequest, nil)
			},
			expectedReq:   expectedRequest,
			expectedError: nil,
		},
		{
			name: "Error retrieving request",
			mockGetByID: func() {
				mockServiceRequestRepo.EXPECT().
					GetServiceRequestByID(requestID).
					Return(nil, errors.New("request not found"))
			},
			expectedReq:   nil,
			expectedError: errors.New("request not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockGetByID()

			result, err := serviceProviderService.GetServiceRequestByID(requestID)

			assert.Equal(t, tt.expectedReq, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
