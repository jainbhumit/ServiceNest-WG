package controllers_test

import (
	"context"
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

func TestGetAvailableServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockHouseholderService(ctrl)
	controller := controllers.NewHouseholderController(mockService)

	t.Run("Get all services", func(t *testing.T) {
		mockService.EXPECT().GetAvailableServices().Return([]model.Service{}, nil)

		req, err := http.NewRequest("GET", "/services", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.GetAvailableServices)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
	t.Run("Error while Get all services", func(t *testing.T) {
		mockService.EXPECT().GetAvailableServices().Return(nil, errors.New("Error getting all services"))

		req, err := http.NewRequest("GET", "/services", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.GetAvailableServices)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("Get services by category", func(t *testing.T) {
		mockService.EXPECT().GetServicesByCategory("Cleaning").Return([]model.Service{}, nil)

		req, err := http.NewRequest("GET", "/services?category=Cleaning", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.GetAvailableServices)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Error while fetching Get service by category", func(t *testing.T) {
		mockService.EXPECT().GetServicesByCategory("Cleaning").Return(nil, errors.New("Error fetching services"))

		req, err := http.NewRequest("GET", "/services?category=Cleaning", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.GetAvailableServices)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestRequestService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockHouseholderService(ctrl)
	controller := controllers.NewHouseholderController(mockService)

	t.Run("Valid service request", func(t *testing.T) {
		// Set the expectation for the mock method
		mockService.EXPECT().
			RequestService("provider123", "House Cleaning", gomock.Any()).
			Return("123", nil) // Return appropriate value or error as needed

		// Set up the request
		requestBody := `{"service_name": "House Cleaning", "scheduled_time": "2023-09-22 15:04"}`
		req, err := http.NewRequest("POST", "/service/request", strings.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		// Add userID and role to context
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.RequestService)
		handler.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("Invalid service request with malformed JSON", func(t *testing.T) {
		// Set up the request with malformed JSON
		requestBody := `{"service_name": , "scheduled_time": "2023-09-22 15:04"}`
		req, err := http.NewRequest("POST", "/service/request", strings.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		// Add userID and role to context
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.RequestService)
		handler.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid input")
	})

	t.Run("Invalid role", func(t *testing.T) {
		// Set up the request with invalid role
		requestBody := `{"service_name": "House Cleaning", "scheduled_time": "2023-09-22 15:04"}`
		req, err := http.NewRequest("POST", "/service/request", strings.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		// Add userID and an invalid role to context
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "ServiceProvider"))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.RequestService)
		handler.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid role")
	})

	t.Run("Missing user ID for Admin role", func(t *testing.T) {
		// Set up the request with Admin role but no user_id query parameter
		requestBody := `{"service_name": "House Cleaning", "scheduled_time": "2023-09-22 15:04"}`
		req, err := http.NewRequest("POST", "/service/request", strings.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		// Add Admin role but omit userID
		req = req.WithContext(context.WithValue(req.Context(), "role", "Admin"))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.RequestService)
		handler.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "user ID is required")
	})

	t.Run("Invalid time format", func(t *testing.T) {
		// Set up the request with an invalid time format
		requestBody := `{"service_name": "House Cleaning", "scheduled_time": "2023-09-2215:04"}`
		req, err := http.NewRequest("POST", "/service/request", strings.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		// Add userID and role to context
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.RequestService)
		handler.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid time format")
	})

	t.Run("Error requesting service", func(t *testing.T) {
		// Set the expectation for the mock method
		mockService.EXPECT().
			RequestService("provider123", "House Cleaning", gomock.Any()).
			Return("", errors.New("error requesting service")) // Return appropriate value or error as needed

		// Set up the request
		requestBody := `{"service_name": "House Cleaning", "scheduled_time": "2023-09-22 15:04"}`
		req, err := http.NewRequest("POST", "/service/request", strings.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		// Add userID and role to context
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.RequestService)
		handler.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "error requesting service")
	})
}

func TestCancelServiceRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockHouseholderService(ctrl)
	controller := controllers.NewHouseholderController(mockService)

	t.Run("Cancel service request", func(t *testing.T) {
		// Mock service expectation
		mockService.EXPECT().CancelServiceRequest("requestID", "provider123").Return(nil)

		// Set up the router and request with the path variable
		req, err := http.NewRequest("DELETE", "/service/cancel/requestID", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder")) // Add role to context

		rr := httptest.NewRecorder()

		// Set up router to handle path variable extraction
		router := mux.NewRouter()
		router.HandleFunc("/service/cancel/{request_id}", controller.CancelServiceRequest)
		router.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Error canceling service request", func(t *testing.T) {
		// Mock service expectation for an error scenario
		mockService.EXPECT().CancelServiceRequest("requestID", "provider123").Return(errors.New("not found"))

		// Set up the router and request with the path variable
		req, err := http.NewRequest("DELETE", "/service/cancel/requestID", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder")) // Add role to context

		rr := httptest.NewRecorder()

		// Set up router to handle path variable extraction
		router := mux.NewRouter()
		router.HandleFunc("/service/cancel/{request_id}", controller.CancelServiceRequest)
		router.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Invalid role canceling service request", func(t *testing.T) {

		// Set up the router and request with the path variable
		req, err := http.NewRequest("DELETE", "/service/cancel/requestID", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "ServiceProvider")) // Add role to context

		rr := httptest.NewRecorder()

		// Set up router to handle path variable extraction
		router := mux.NewRouter()
		router.HandleFunc("/service/cancel/{request_id}", controller.CancelServiceRequest)
		router.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestRescheduleServiceRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockHouseholderService(ctrl)
	controller := controllers.NewHouseholderController(mockService)

	t.Run("Successfully reschedule service request", func(t *testing.T) {
		// Mock service expectation
		mockService.EXPECT().RescheduleServiceRequest("requestID", gomock.Any(), "householder123").Return(nil)

		// Prepare valid request body
		reqBody := `{"id":"requestID", "scheduled_time":"2024-09-23 14:00"}`
		req, err := http.NewRequest("POST", "/service/reschedule", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		// Call the handler
		handler := http.HandlerFunc(controller.RescheduleServiceRequest)
		handler.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Error rescheduling service request due to invalid input", func(t *testing.T) {
		reqBody := `{"id":"", "scheduled_time":"2024-09-23 14:00"}` // Invalid input
		req, err := http.NewRequest("POST", "/service/reschedule", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}

		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.RescheduleServiceRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid request body")
	})
	t.Run("Error rescheduling service request due to invalid json input", func(t *testing.T) {
		reqBody := `{"id":"", "scheduledtime":"2024-09-23 14:00"}` // Invalid input
		req, err := http.NewRequest("POST", "/service/reschedule", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}

		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.RescheduleServiceRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

	})

	t.Run("Invalid time format", func(t *testing.T) {
		reqBody := `{"id":"requestID", "scheduled_time":"invalid-time"}`
		req, err := http.NewRequest("POST", "/service/reschedule", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}

		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.RescheduleServiceRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid time format")
	})

	t.Run("Error occuring due to db failure", func(t *testing.T) {
		mockService.EXPECT().RescheduleServiceRequest("requestID", gomock.Any(), "householder123").Return(errors.New("db error"))

		// Prepare valid request body
		reqBody := `{"id":"requestID", "scheduled_time":"2024-09-23 14:00"}`
		req, err := http.NewRequest("POST", "/service/reschedule", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		// Call the handler
		handler := http.HandlerFunc(controller.RescheduleServiceRequest)
		handler.ServeHTTP(rr, req)

		// Check the response code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Error rescheduling service request due to invalid role", func(t *testing.T) {
		reqBody := `{"id":"", "scheduled_time":"2024-09-23 14:00"}` // Invalid input
		req, err := http.NewRequest("POST", "/service/reschedule", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}

		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "ServiceProvider"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.RescheduleServiceRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

	})
}

func TestViewBookingHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockHouseholderService(ctrl)
	controller := controllers.NewHouseholderController(mockService)

	t.Run("Successfully fetch booking history", func(t *testing.T) {
		// Mock service expectation
		mockService.EXPECT().ViewStatus("householder123").Return([]model.ServiceRequest{
			{ID: "requestID", Status: "Pending", ServiceName: "Plumbing"},
		}, nil)

		req, err := http.NewRequest("GET", "/service/history", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.ViewBookingHistory)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("No service requests found", func(t *testing.T) {
		mockService.EXPECT().ViewStatus("householder123").Return([]model.ServiceRequest{}, nil)

		req, err := http.NewRequest("GET", "/service/history", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.ViewBookingHistory)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "No service request found")
	})

	t.Run("Error fetching booking history", func(t *testing.T) {
		// Mock service expectation
		mockService.EXPECT().ViewStatus("householder123").Return(nil, errors.New("error fetching status"))

		req, err := http.NewRequest("GET", "/service/history", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.ViewBookingHistory)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("Successfully fetch booking history with provider detail", func(t *testing.T) {
		// Mock service expectation
		providerDetail := []model.ServiceProviderDetails{
			{
				Name:    "provider1",
				Contact: "77994844",
				Address: "noida",
			},
		}

		mockService.EXPECT().ViewStatus("householder123").Return([]model.ServiceRequest{{
			ID: "requestID", Status: "Pending", ServiceName: "Plumbing", ApproveStatus: false, ProviderDetails: providerDetail,
		}}, nil)

		req, err := http.NewRequest("GET", "/service/history", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.ViewBookingHistory)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid role for fetching booking history", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/service/history", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "ServiceProvider"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.ViewBookingHistory)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestApproveRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockHouseholderService(ctrl)
	controller := controllers.NewHouseholderController(mockService)

	t.Run("Successfully approve service request", func(t *testing.T) {
		mockService.EXPECT().ApproveServiceRequest("requestID", "providerID", "householder123").Return(nil)

		reqBody := `{"request_id":"requestID", "provider_id":"providerID"}`
		req, err := http.NewRequest("POST", "/service/approve", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.ApproveRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
	t.Run("Invalid json for approve service request", func(t *testing.T) {

		reqBody := `{"requestid":"requestID", "provider_id":"providerID"}`
		req, err := http.NewRequest("POST", "/service/approve", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.ApproveRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("Error approving service request due to missing fields", func(t *testing.T) {
		reqBody := `{"request_id":"", "provider_id":""}`
		req, err := http.NewRequest("POST", "/service/approve", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}

		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.ApproveRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid request body")
	})

	t.Run("Error approving service request", func(t *testing.T) {
		mockService.EXPECT().ApproveServiceRequest("requestID", "providerID", "householder123").Return(errors.New("error approving service"))

		reqBody := `{"request_id":"requestID", "provider_id":"providerID"}`
		req, err := http.NewRequest("POST", "/service/approve", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "Householder"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.ApproveRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("Invalid role for approving service request", func(t *testing.T) {
		mockService.EXPECT().ApproveServiceRequest("requestID", "providerID", "householder123").Return(errors.New("error approving service"))

		reqBody := `{"request_id":"requestID", "provider_id":"providerID"}`
		req, err := http.NewRequest("POST", "/service/approve", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))
		req = req.WithContext(context.WithValue(req.Context(), "role", "ServiceProvider"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.ApproveRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
func TestLeaveReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockHouseholderService(ctrl)
	controller := controllers.NewHouseholderController(mockService)

	t.Run("Successfully add a review", func(t *testing.T) {
		mockService.EXPECT().AddReview("providerID", "householder123", "serviceID", "Great service!", 5.0).Return(nil)

		reqBody := `{
			"service_id": "serviceID", 
			"provider_id": "providerID", 
			"review_text": "Great service!", 
			"rating": 5.0
		}`
		req, err := http.NewRequest("POST", "/service/review", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.LeaveReview)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid review due to bad rating", func(t *testing.T) {
		reqBody := `{
			"service_id": "serviceID", 
			"provider_id": "providerID", 
			"review_text": "Great service!", 
			"rating": 6.0
		}`
		req, err := http.NewRequest("POST", "/service/review", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.LeaveReview)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Rating should be between 1 and 5")
	})
	t.Run("Invalid review due to wrong json", func(t *testing.T) {
		reqBody := `{
			"service_id": "serviceID", 
			"provider_id": "providerID", 
			"review_text: "Great service!", 
			"rating": 6.0
		}`
		req, err := http.NewRequest("POST", "/service/review", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.LeaveReview)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		//assert.Contains(t, rr.Body.String(), "Rating should be between 1 and 5")
	})
	t.Run("Invalid review due to invalid request body", func(t *testing.T) {
		reqBody := `{
			"service_id": "serviceID", 
			"provider_id": "providerID", 
			"review_text": "", 
			"rating": 6.0
		}`
		req, err := http.NewRequest("POST", "/service/review", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.LeaveReview)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		//assert.Contains(t, rr.Body.String(), "Rating should be between 1 and 5")
	})
	t.Run("Error adding a review", func(t *testing.T) {
		mockService.EXPECT().AddReview("providerID", "householder123", "serviceID", "Great service!", 5.0).Return(errors.New("error adding review"))

		reqBody := `{
			"service_id": "serviceID", 
			"provider_id": "providerID", 
			"review_text": "Great service!", 
			"rating": 5.0
		}`
		req, err := http.NewRequest("POST", "/service/review", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), "userID", "householder123"))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.LeaveReview)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestViewApprovedRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockHouseholderService(ctrl)
	controller := controllers.NewHouseholderController(mockService)

	t.Run("Successful retrieval of approved requests", func(t *testing.T) {

		mockService.EXPECT().ViewApprovedRequests("provider123").Return([]model.ServiceRequest{}, nil)

		req, err := http.NewRequest("GET", "/service/approved", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.ViewApprovedRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

	})

	t.Run("No approved requests found", func(t *testing.T) {
		mockService.EXPECT().ViewApprovedRequests("provider123").Return([]model.ServiceRequest{}, nil)

		req, err := http.NewRequest("GET", "/service/approved", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.ViewApprovedRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Error retrieving approved requests", func(t *testing.T) {
		mockService.EXPECT().ViewApprovedRequests("provider123").Return(nil, errors.New("db error"))

		req, err := http.NewRequest("GET", "/service/approved", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.ViewApprovedRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
func TestViewApprovedRequestAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockHouseholderService(ctrl)
	controller := controllers.NewHouseholderController(mockService)

	t.Run("Successful retrieval of approved requests with providers", func(t *testing.T) {
		// Mock data with multiple approved requests, one with approved providers
		approvedRequests := []model.ServiceRequest{
			{
				ID:            "req1",
				ServiceName:   "House Cleaning",
				ServiceID:     "service1",
				RequestedTime: time.Now(),
				ScheduledTime: time.Now().Add(2 * time.Hour),
				Status:        "Approved",
				ApproveStatus: true,
				ProviderDetails: []model.ServiceProviderDetails{
					{
						Name:    "Provider A",
						Contact: "123456789",
						Address: "123 Street",
						Price:   "100",
						Rating:  4.5,
						Approve: true,
					},
					{
						Name:    "Provider B",
						Contact: "987654321",
						Address: "456 Avenue",
						Price:   "120",
						Rating:  4.8,
						Approve: false,
					},
				},
			},
		}

		mockService.EXPECT().ViewApprovedRequests("provider123").Return(approvedRequests, nil)

		req, err := http.NewRequest("GET", "/service/approved", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.ViewApprovedRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Ensure the provider details are in the response
		assert.Contains(t, rr.Body.String(), "Provider A")
		assert.NotContains(t, rr.Body.String(), "Provider B") // As Provider B is not approved
	})

	t.Run("Successful retrieval of approved requests without providers", func(t *testing.T) {
		// Mock data with an approved request but no provider approvals
		approvedRequests := []model.ServiceRequest{
			{
				ID:            "req1",
				ServiceName:   "House Cleaning",
				ServiceID:     "service1",
				RequestedTime: time.Now(),
				ScheduledTime: time.Now().Add(2 * time.Hour),
				Status:        "Approved",
				ApproveStatus: true,
				ProviderDetails: []model.ServiceProviderDetails{
					{
						Name:    "Provider A",
						Contact: "123456789",
						Address: "123 Street",
						Price:   "100",
						Rating:  4.5,
						Approve: false,
					},
				},
			},
		}

		mockService.EXPECT().ViewApprovedRequests("provider123").Return(approvedRequests, nil)

		req, err := http.NewRequest("GET", "/service/approved", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.ViewApprovedRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Ensure no provider details are included in the response
		assert.NotContains(t, rr.Body.String(), "Provider A")
	})

	t.Run("No approved requests found", func(t *testing.T) {
		mockService.EXPECT().ViewApprovedRequests("provider123").Return([]model.ServiceRequest{}, nil)

		req, err := http.NewRequest("GET", "/service/approved", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.ViewApprovedRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Ensure the response message mentions no approved requests found
		assert.Contains(t, rr.Body.String(), "No approved service requests found")
	})

	t.Run("Error retrieving approved requests", func(t *testing.T) {
		mockService.EXPECT().ViewApprovedRequests("provider123").Return(nil, errors.New("db error"))

		req, err := http.NewRequest("GET", "/service/approved", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userID", "provider123"))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.ViewApprovedRequest)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "db error")
	})
}
