package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"serviceNest/controllers"
	"serviceNest/model"
	"serviceNest/tests/mocks"
	"strings"
	"testing"
	"time"
)

func TestAddService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	requestBody := `{
        "name": "Cleaning",
        "description": "House cleaning service",
        "price": 100,
        "category": "Home"
    }`

	req, _ := http.NewRequest("POST", "/add-service", bytes.NewBufferString(requestBody))

	role := "ServiceProvider"
	ctx := context.WithValue(req.Context(), "userID", "provider123") // Set userID correctly
	ctx = context.WithValue(ctx, "role", role)                       // Also set role
	req = req.WithContext(ctx)

	// Set a fixed ID for testing
	serviceID := "fixed-service-id"

	// Mock the ID generation function to return the fixed ID
	original := controllers.GenerateUniqueID
	defer func() { controllers.GenerateUniqueID = original }()
	controllers.GenerateUniqueID = func() string {
		return serviceID
	}

	// Expect AddService to be called with specific arguments and return nil error
	newService := model.Service{
		ID:          serviceID, // Using the fixed ID here
		Name:        "Cleaning",
		Description: "House cleaning service",
		Price:       100,
		Category:    "Home",
		ProviderID:  "provider123",
	}

	mockService.EXPECT().AddService("provider123", newService).Return(nil)

	serviceProviderController.AddService(w, req)

	// Validate the response code
	assert.Equal(t, http.StatusOK, w.Code)

	// Decode the full response and check for service_id inside the data
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	// Assert the message is correct
	assert.Equal(t, "Service added successfully", resp["message"])

	// Assert the "data" field contains the service_id
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, serviceID, data["service_id"]) // Assert the returned ID is correct
}

func TestAddService_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	requestBody := `{
		"name": "",
		"description": "",
		"price": 100,
		"category": ""
	}`

	req, _ := http.NewRequest("POST", "/add-service", bytes.NewBufferString(requestBody))

	userID := "provider123"
	role := "ServiceProvider"
	ctx := context.WithValue(req.Context(), "userID", userID)
	ctx = context.WithValue(req.Context(), "role", role)
	req = req.WithContext(ctx)

	serviceProviderController.AddService(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Invalid request body", resp["message"])
}

func TestViewServices_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/view-services", nil)
	role := "ServiceProvider"
	ctx := context.WithValue(req.Context(), "userID", "provider123") // Set userID correctly
	ctx = context.WithValue(ctx, "role", role)                       // Also set role
	req = req.WithContext(ctx)

	services := []model.Service{
		{
			ID:          "service1",
			Name:        "Cleaning",
			Description: "House cleaning service",
			Price:       100,
			Category:    "Home",
			ProviderID:  "provider123",
		},
	}

	mockService.EXPECT().ViewServices("provider123").Return(services, nil)

	serviceProviderController.ViewServices(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "service fetch successfully", resp["message"])

	// Assert the services inside "data" as a list
	data := resp["data"].([]interface{}) // Fix: Expect a list, not a map
	assert.NotNil(t, data)
	assert.Len(t, data, 1) // Expecting 1 service

	// Further check that the service in the list has expected values
	serviceData := data[0].(map[string]interface{})
	assert.Equal(t, "service1", serviceData["id"])
	assert.Equal(t, "Cleaning", serviceData["name"])
	assert.Equal(t, "House cleaning service", serviceData["description"])
	assert.Equal(t, 100.0, serviceData["price"]) // JSON unmarshals numbers as float64
	assert.Equal(t, "Home", serviceData["category"])
}

func TestViewServices_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/view-services", nil)
	role := "ServiceProvider"
	ctx := context.WithValue(req.Context(), "userID", "provider123") // Set userID correctly
	ctx = context.WithValue(ctx, "role", role)                       // Also set role
	req = req.WithContext(ctx)

	mockService.EXPECT().ViewServices("provider123").Return(nil, errors.New("error fetching services"))

	serviceProviderController.ViewServices(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)

	// Check the error message
	assert.Equal(t, "Error fetching services", resp["message"])
	assert.Nil(t, resp["data"]) // Ensure data is nil on error
}

func TestAcceptServiceRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	requestBody := `{
		"request_id": "request123",
		"price": "150"
	}`

	req, _ := http.NewRequest("POST", "/accept-request", bytes.NewBufferString(requestBody))
	userID := "provider123"
	role := "ServiceProvider"
	ctx := context.WithValue(req.Context(), "userID", userID)
	ctx = context.WithValue(ctx, "role", role)
	req = req.WithContext(ctx)

	mockService.EXPECT().AcceptServiceRequest("provider123", "request123", "150").Return(nil)

	serviceProviderController.AcceptServiceRequest(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)

	// Check if the response message is as expected
	assert.Equal(t, "Request accept successfully", resp["message"])

	// Check if the data is nil or not
	//assert.NotNil(t, resp["data"]) // This assumes you will have a data field in the response
}

func TestAcceptServiceRequest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	requestBody := `{
        "request_id": "request123",
        "price": "150"
    }`

	req, _ := http.NewRequest("POST", "/accept-request", bytes.NewBufferString(requestBody))
	userID := "provider123"
	role := "ServiceProvider"
	ctx := context.WithValue(req.Context(), "userID", userID)
	ctx = context.WithValue(ctx, "role", role)
	req = req.WithContext(ctx)

	mockService.EXPECT().AcceptServiceRequest("provider123", "request123", "150").Return(errors.New("error accepting request"))

	serviceProviderController.AcceptServiceRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)

	// Safely check if "message" exists and is not nil
	message, ok := resp["message"].(string)
	if !ok || message == "" {
		t.Fatalf("Expected 'message' field to be non-nil string, got %v", resp["message"])
	}

	assert.Equal(t, "Error accepting service request", message)

	// Optionally check for other fields in the response if needed
	assert.Nil(t, resp["data"]) // Assuming data might be present, even in error
}
func TestAcceptServiceRequest_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	t.Run("Invalid Body", func(t *testing.T) {
		w := httptest.NewRecorder()
		requestBody := `{
        "request_id" "request123",
        "price": "150"
    }`

		req, _ := http.NewRequest("POST", "/accept-request", bytes.NewBufferString(requestBody))
		userID := "provider123"
		role := "ServiceProvider"
		ctx := context.WithValue(req.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "role", role)
		req = req.WithContext(ctx)

		//mockService.EXPECT().AcceptServiceRequest("provider123", "request123", "150").Return(errors.New("error accepting request"))

		serviceProviderController.AcceptServiceRequest(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)

		// Safely check if "message" exists and is not nil
		message, ok := resp["message"].(string)
		if !ok || message == "" {
			t.Fatalf("Expected 'message' field to be non-nil string, got %v", resp["message"])
		}

		assert.Equal(t, "Invalid body", message)

		// Optionally check for other fields in the response if needed
		assert.Nil(t, resp["data"]) // Assuming data might be present, even in error
	})

}
func TestViewApprovedRequests_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/view-approved-requests", nil)
	role := "ServiceProvider"
	ctx := context.WithValue(req.Context(), "userID", "provider123")
	ctx = context.WithValue(ctx, "role", role)
	req = req.WithContext(ctx)

	// Define the expected response
	approvedRequests := []model.ServiceRequest{
		{
			ID:            "3983",
			ServiceName:   "",
			ServiceID:     "1950",
			RequestedTime: time.Now(),
			ScheduledTime: time.Now().Add(24 * time.Hour),
			Status:        "Accepted",
		},
		{
			ID:            "3983",
			ServiceName:   "",
			ServiceID:     "1950",
			RequestedTime: time.Now(),
			ScheduledTime: time.Now().Add(24 * time.Hour),
			Status:        "Accepted",
		},
	}

	mockService.EXPECT().ViewApprovedRequestsByProvider("provider123").Return(approvedRequests, nil)

	serviceProviderController.ViewApprovedRequests(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)

	// Check the message field
	assert.Equal(t, "Approve requests fetched successfully", resp["message"])

	// Ensure the data field is a slice
	data, ok := resp["data"].([]interface{})
	if !ok || data == nil {
		t.Fatalf("Expected 'data' field to be a non-nil array, got %v", resp["data"])
	}

	// Loop over each request and validate its structure
	for _, reqData := range data {
		reqMap, ok := reqData.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected request to be a map, got %v", reqData)
		}

		// Validate fields
		assert.Equal(t, "3983", reqMap["request_id"])
		assert.Equal(t, "1950", reqMap["service_id"])
		assert.Equal(t, "Accepted", reqMap["status"])

		// Validate provider details

	}
}

func TestViewApprovedRequests_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/view-approved-requests", nil)
	role := "ServiceProvider"
	ctx := context.WithValue(req.Context(), "userID", "provider123")
	ctx = context.WithValue(ctx, "role", role)
	req = req.WithContext(ctx)

	// Mock the service to return an error
	mockService.EXPECT().ViewApprovedRequestsByProvider("provider123").Return(nil, errors.New("something went wrong"))

	serviceProviderController.ViewApprovedRequests(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Parse the response
	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)

	// Check the message field
	assert.Equal(t, "something went wrong", resp["message"])

	// Ensure data is nil or not present
	data, dataExists := resp["data"]
	assert.False(t, dataExists || data != nil)
}

func TestViewServiceRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/service/request", nil)

	serviceRequests := []model.ServiceRequest{
		{
			ID:                 "123",
			ServiceName:        "ServiceA",
			ServiceID:          "987",
			RequestedTime:      time.Now(),
			ScheduledTime:      time.Now().Add(24 * time.Hour),
			HouseholderAddress: func() *string { addr := "123 Main St"; return &addr }(),
			ApproveStatus:      false,
			Status:             "Pending",
		},
	}

	mockService.EXPECT().GetAllServiceRequests().Return(serviceRequests, nil)

	serviceProviderController.ViewServiceRequest(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Service request fetched successfully", resp["message"])
	data := resp["data"].([]interface{})
	assert.NotNil(t, data)
}

func TestViewServiceRequest_NoRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/service/request", nil)

	mockService.EXPECT().GetAllServiceRequests().Return([]model.ServiceRequest{}, nil)

	serviceProviderController.ViewServiceRequest(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "No pending service requests available", resp["message"])
}

func TestViewServiceRequest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/service/request", nil)

	mockService.EXPECT().GetAllServiceRequests().Return(nil, errors.New("error fetching request"))

	serviceProviderController.ViewServiceRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "error fetching request error fetching request", resp["message"])
}

func TestRemoveService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()

	// Create a request with the service ID in the URL
	req, _ := http.NewRequest("DELETE", "/service/remove/123", nil)

	// Set up the mux vars for the request to mimic a real request with a service ID
	req = mux.SetURLVars(req, map[string]string{"service_id": "123"})

	// Set the user ID in the context
	ctx := context.WithValue(req.Context(), "userID", "provider123")
	req = req.WithContext(ctx)

	// Mock expectation: should expect "provider123" as the provider ID and "123" as the service ID
	mockService.EXPECT().RemoveService("provider123", "123").Return(nil)

	serviceProviderController.RemoveService(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Service removed successfully", resp["message"])
}

func TestRemoveService_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()

	// Create a request with the service ID in the URL
	req, _ := http.NewRequest("DELETE", "/service/remove/123", nil)

	// Set up the mux vars for the request to mimic a real request with a service ID
	req = mux.SetURLVars(req, map[string]string{"service_id": "123"})

	// Set the user ID in the context
	ctx := context.WithValue(req.Context(), "userID", "provider123")
	req = req.WithContext(ctx)

	// Mock expectation: should expect "provider123" as the provider ID and "123" as the service ID
	mockService.EXPECT().RemoveService("provider123", "123").Return(errors.New("failed to remove service"))

	// Call the controller method
	serviceProviderController.RemoveService(w, req)

	// Assert that the response status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert that the response body contains the error message
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "failed to remove service", resp["message"])
}

func TestUpdateService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceProviderService(ctrl)
	serviceProviderController := controllers.NewServiceProviderController(mockService)

	w := httptest.NewRecorder()

	// Define the service update request body
	requestBody := `{
        "name": "New Service",
        "description": "Updated Description",
        "price": 150.50,
        "category": "CategoryA"
    }`

	// Create a request with the service ID in the URL
	req, _ := http.NewRequest("PUT", "/service/update/123", strings.NewReader(requestBody))

	// Set up the mux vars for the request to mimic a real request with a service ID
	req = mux.SetURLVars(req, map[string]string{"service_id": "123"})

	// Set the user ID in the context
	ctx := context.WithValue(req.Context(), "userID", "provider123")
	req = req.WithContext(ctx)

	// Define the expected updated service object
	updatedService := model.Service{
		ID:          "123",
		Name:        "New Service",
		Description: "Updated Description",
		Price:       150.50,
		Category:    "CategoryA",
		ProviderID:  "provider123",
	}

	// Mock expectation: should expect "provider123" as the provider ID and "123" as the service ID
	mockService.EXPECT().UpdateService("provider123", "123", updatedService).Return(nil)

	// Call the controller method
	serviceProviderController.UpdateService(w, req)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert that the response body contains a success message
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Service updated successfully", resp["message"])
}

func TestUpdateService_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceProviderController := controllers.NewServiceProviderController(nil)
	w := httptest.NewRecorder()

	reqBody := `{"name": "", "description": "Invalid Description", "price": 150.50, "category": ""}`
	req, _ := http.NewRequest("PUT", "/service/update/123", strings.NewReader(reqBody))

	serviceProviderController.UpdateService(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Invalid request body", resp["message"])
}

//	func TestServiceProviderController_ViewReviews(t *testing.T) {
//		ctrl := gomock.NewController(t)
//		defer ctrl.Finish()
//
//		mockServiceProviderService := mocks.NewMockServiceProviderService(ctrl)
//		controller := controllers.NewServiceProviderController(mockServiceProviderService)
//
//		// Create a test HTTP request
//		req, err := http.NewRequest("GET", "/reviews", nil)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		// Set a context value for userID
//		ctx := context.WithValue(req.Context(), "userID", "providerID")
//		req = req.WithContext(ctx)
//
//		// Create a response recorder to capture the response
//		rr := httptest.NewRecorder()
//
//		// Test case: Successful fetch of reviews
//		reviews := []model.Review{
//			{ID: "1", ProviderID: "providerID", Comments: "Great service!", Rating: 5},
//			{ID: "2", ProviderID: "providerID", Comments: "Satisfactory experience.", Rating: 3},
//		}
//
//		mockServiceProviderService.EXPECT().
//			GetReviews("providerID").
//			Return(reviews, nil).
//			Times(1)
//
//		// Call the method
//		controller.ViewReviews(rr, req)
//
//		// Assert the response status and body
//		assert.Equal(t, http.StatusOK, rr.Code)
//		assert.JSONEq(t, `{"data":[{"ID":"1","ProviderID":"providerID","Comment":"Great service!","Rating":5},{"ID":"2","ProviderID":"providerID","Comment":"Satisfactory experience.","Rating":3}],"message":"Reviews fetched successfully"}`, rr.Body.String())
//
//		// Test case: Error fetching reviews
//		rr = httptest.NewRecorder() // Reset the recorder
//
//		mockServiceProviderService.EXPECT().
//			GetReviews("providerID").
//			Return(nil, errors.New("fetch error")).
//			Times(1)
//
//		controller.ViewReviews(rr, req)
//
//		// Assert the response status and body for the error case
//		assert.Equal(t, http.StatusInternalServerError, rr.Code)
//		assert.JSONEq(t, `{"error":"Failed to fetch reviews"}`, rr.Body.String())
//	}
func TestServiceProviderController_ViewReviews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderService := mocks.NewMockServiceProviderService(ctrl)

	controller := controllers.NewServiceProviderController(mockServiceProviderService)

	// Sample reviews to return
	reviews := []model.Review{
		{
			ID:            "1",
			ServiceID:     "existingServiceID",
			HouseholderID: "householderID",
			ProviderID:    "providerID",
			Rating:        5,
			Comments:      "Great service!",
		},
		{
			ID:            "2",
			ServiceID:     "existingServiceID",
			HouseholderID: "householderID",
			ProviderID:    "providerID",
			Rating:        3,
			Comments:      "Satisfactory experience.",
		},
	}

	// Mock the service call
	mockServiceProviderService.EXPECT().
		GetReviews("providerID").
		Return(reviews, nil).
		Times(1)

	// Create a request with context
	req := httptest.NewRequest(http.MethodGet, "/reviews", nil)
	ctx := context.WithValue(req.Context(), "userID", "providerID")
	req = req.WithContext(ctx)

	// Record the response
	rr := httptest.NewRecorder()

	// Call the ViewReviews method
	controller.ViewReviews(rr, req)

	// Check the response code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expectedJSON := `{"data":[{"id":"1","service_id":"existingServiceID","householder_id":"householderID","provider_id":"providerID","rating":5,"comments":"Great service!","review_date":"0001-01-01T00:00:00Z"},{"id":"2","service_id":"existingServiceID","householder_id":"householderID","provider_id":"providerID","rating":3,"comments":"Satisfactory experience.","review_date":"0001-01-01T00:00:00Z"}],"message":"Reviews fetched successfully","status":"Success"}`
	assert.JSONEq(t, expectedJSON, rr.Body.String())
}

func TestViewReviewsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProviderService := mocks.NewMockServiceProviderService(ctrl)

	controller := controllers.NewServiceProviderController(mockServiceProviderService)

	// Mock the service call
	mockServiceProviderService.EXPECT().
		GetReviews("providerID").
		Return(nil, errors.New("failed to get reviews")).
		Times(1)

	// Create a request with context
	req := httptest.NewRequest(http.MethodGet, "/reviews", nil)
	ctx := context.WithValue(req.Context(), "userID", "providerID")
	req = req.WithContext(ctx)

	// Record the response
	rr := httptest.NewRecorder()

	// Call the ViewReviews method
	controller.ViewReviews(rr, req)

	// Check the response code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Check the response body
	expectedJSON := `{"error_code":1003,"message":"Failed to fetch reviews","status":"Fail"}`
	assert.JSONEq(t, expectedJSON, rr.Body.String())
}
