package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"serviceNest/interfaces"
	"serviceNest/logger"
	"serviceNest/model"
	"serviceNest/response"
	"time"
)

type HouseholderController struct {
	householderService interfaces.HouseholderService
}

func NewHouseholderController(householderService interfaces.HouseholderService) *HouseholderController {
	return &HouseholderController{
		householderService: householderService,
	}
}

func (h *HouseholderController) GetAvailableServices(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		services, err := h.householderService.GetAvailableServices()
		if err != nil {
			logger.Error("error fetching all service", nil)
			response.ErrorResponse(w, http.StatusInternalServerError, "internal server error", 1006)
			return
		}
		response.SuccessResponse(w, services, "Available services", http.StatusOK)
	} else {
		services, err := h.householderService.GetServicesByCategory(category)
		if err != nil {
			logger.Error("error fetching services", nil)
			response.ErrorResponse(w, http.StatusInternalServerError, "internal server error", 1006)
			return
		}
		response.SuccessResponse(w, services, "Available services", http.StatusOK)

	}
}

func (h *HouseholderController) RequestService(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ServiceName   string `json:"service_name" validate:"required"`
		ScheduledTime string `json:"scheduled_time" validate:"required"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Invalid input", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	// Validate the request body
	err := validate.Struct(request)
	if err != nil {
		logger.Error("Invalid request body", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}
	role := r.Context().Value("role").(string)
	var householderID string
	if role == "Admin" {
		householderID = r.URL.Query().Get("user_id")
		if householderID == "" {
			logger.Error("No query param", nil)
			response.ErrorResponse(w, http.StatusBadRequest, "user ID is required", 2001)
			return
		}
	} else if role == "Householder" {
		householderID = r.Context().Value("userID").(string)
	} else {
		logger.Error("Invalid role", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid role", 1007)
		return
	}
	// Example of retrieving the householder ID from the context

	scheduleTime, err := time.Parse("2006-01-02 15:04", request.ScheduledTime)
	if err != nil {
		logger.Error("Invalid request body", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid time format", 1001)
		return
	}

	// Pass the request's ServiceName and ScheduledTime to the service layer
	requestID, err := h.householderService.RequestService(householderID, request.ServiceName, &scheduleTime)
	if err != nil {
		logger.Error("error requesting service", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "error requesting service", 1006)
		return
	}

	logger.Info("Service request successfully", nil)
	var respone struct {
		ID string `json:"request_id"`
	}
	respone.ID = requestID
	response.SuccessResponse(w, respone, "Service request successfully", http.StatusCreated)
}

func (h *HouseholderController) CancelServiceRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID, ok := vars["request_id"]
	if !ok {
		logger.Error("Missing request Id in params", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Missing request Id in params", 2002)
		return
	}
	role := r.Context().Value("role").(string)
	var householderID string
	if role == "Admin" {
		householderID = r.URL.Query().Get("user_id")
		if householderID == "" {
			logger.Error("No query param", nil)
			response.ErrorResponse(w, http.StatusBadRequest, "user ID is required", 2001)
			return
		}
	} else if role == "Householder" {
		householderID = r.Context().Value("userID").(string)
	} else {
		logger.Error("Invalid role", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid role", 1007)
		return
	}
	err := h.householderService.CancelServiceRequest(requestID, householderID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error cancelling request %v", err), nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error", 1006)
		return
	}
	logger.Info("Request cancelled successfully", nil)
	response.SuccessResponse(w, nil, "Request cancelled successfully", http.StatusOK)
	//color.Green("Service request %s has been successfully canceled.", requestID)
}

func (h *HouseholderController) RescheduleServiceRequest(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID            string `json:"id" validate:"required"`
		ScheduledTime string `json:"scheduled_time" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Invalid input", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}
	err := validate.Struct(request)
	if err != nil {
		logger.Error("Invalid request body", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}

	newTime, err := time.Parse("2006-01-02 15:04", request.ScheduledTime)
	if err != nil {
		logger.Error("Invalid request body", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid time format", 1001)
		return
	}
	role := r.Context().Value("role").(string)
	var householderID string
	if role == "Admin" {
		householderID = r.URL.Query().Get("user_id")
		if householderID == "" {
			logger.Error("No query param", nil)
			response.ErrorResponse(w, http.StatusBadRequest, "user ID is required", 2001)
			return
		}
	} else if role == "Householder" {
		householderID = r.Context().Value("userID").(string)
	} else {
		logger.Error("Invalid role", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid role", 1007)
		return
	}
	err = h.householderService.RescheduleServiceRequest(request.ID, newTime, householderID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error rescheduling service %v", err), nil)
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), 1008)
		//color.Red("Error rescheduling service_test request: %v", err)
		return
	}
	logger.Info("Successfully rescheduled service request", nil)
	response.SuccessResponse(w, nil, "service request has been successfully rescheduled", http.StatusOK)
	//color.Green("Service request %s has been successfully rescheduled.", requestID)
}
func (h *HouseholderController) ViewBookingHistory(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role").(string)
	var householderID string
	if role == "Admin" {
		householderID = r.URL.Query().Get("user_id")
		if householderID == "" {
			logger.Error("No query param", nil)
			response.ErrorResponse(w, http.StatusBadRequest, "user ID is required", 2001)
			return
		}
	} else if role == "Householder" {
		householderID = r.Context().Value("userID").(string)
	} else {
		logger.Error("Invalid role", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid role", 1007)
		return
	}

	// Fetch service requests for the householder
	serviceRequests, err := h.householderService.ViewStatus(householderID)
	if err != nil {
		logger.Error("Failed to fetch service requests", map[string]interface{}{
			"householderID": householderID,
			"error":         err.Error(),
		})
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch service requests", 1003)
		return
	}

	// Check if there are any service requests
	if len(serviceRequests) == 0 {
		logger.Info("No service requests found for householder", map[string]interface{}{
			"householderID": householderID,
		})
		response.SuccessResponse(w, nil, "No service request found", http.StatusOK)
		return
	}

	type responseStruct struct {
		ID              string                         `json:"request_id"`
		ServiceName     string                         `json:"service_name,omitempty"`
		ServiceID       string                         `json:"service_id" `
		RequestedTime   time.Time                      `json:"requested_time"`
		ScheduledTime   time.Time                      `json:"scheduled_time"`
		Status          string                         `json:"status"` // Pending, Accepted, Approved, Cancelled
		ApproveStatus   bool                           `json:"approve_status" bson:"approveStatus"`
		ProviderDetails []model.ServiceProviderDetails `json:"provider_details,omitempty" bson:"providerDetails,omitempty"`
	}
	responseBody := make([]responseStruct, 0)
	for _, request := range serviceRequests {
		currRequest := &responseStruct{
			ID:            request.ID,
			ServiceName:   request.ServiceName,
			ServiceID:     request.ServiceID,
			RequestedTime: request.RequestedTime,
			ScheduledTime: request.ScheduledTime,
			Status:        request.Status,
		}
		if request.Status == "Accepted" && request.ProviderDetails != nil && !request.ApproveStatus {
			for _, provider := range request.ProviderDetails {
				currRequest.ProviderDetails = append(currRequest.ProviderDetails, model.ServiceProviderDetails{
					ServiceProviderID: provider.ServiceProviderID,
					Name:              provider.Name,
					Contact:           provider.Contact,
					Address:           provider.Address,
					Price:             provider.Price,
					Rating:            provider.Rating,
				})

			}
		}
		responseBody = append(responseBody, *currRequest)
	}
	// Return the service requests in the response
	logger.Info("Service requests fetched successfully", map[string]interface{}{
		"householderID": householderID,
	})
	response.SuccessResponse(w, responseBody, "Service Request fetched successfully", http.StatusOK)
}

func (h *HouseholderController) ViewApprovedRequest(w http.ResponseWriter, r *http.Request) {
	householderID := r.Context().Value("userID").(string)

	approvedRequests, err := h.householderService.ViewApprovedRequests(householderID)
	if err != nil {
		logger.Error(err.Error(), nil)
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), 1008)
		return
	}

	if len(approvedRequests) == 0 {
		logger.Info("No approve service requests", nil)
		response.SuccessResponse(w, nil, "No approved service requests found", http.StatusOK)
		//fmt.Println("No approved service requests found.")
		return
	}
	type responseStruct struct {
		ID              string                         `json:"request_id"`
		ServiceName     string                         `json:"service_name,omitempty"`
		ServiceID       string                         `json:"service_id" `
		RequestedTime   time.Time                      `json:"requested_time"`
		ScheduledTime   time.Time                      `json:"scheduled_time"`
		Status          string                         `json:"status"`
		ProviderDetails []model.ServiceProviderDetails `json:"provider_details,omitempty" bson:"providerDetails,omitempty"`
	}
	responseBody := make([]responseStruct, 0)
	for _, request := range approvedRequests {
		if request.ApproveStatus {
			currRequest := &responseStruct{
				ID:            request.ID,
				ServiceName:   request.ServiceName,
				ServiceID:     request.ServiceID,
				RequestedTime: request.RequestedTime,
				ScheduledTime: request.ScheduledTime,
				Status:        request.Status,
			}
			for _, provider := range request.ProviderDetails {
				if provider.Approve {
					currRequest.ProviderDetails = append(currRequest.ProviderDetails, model.ServiceProviderDetails{
						Name:    provider.Name,
						Contact: provider.Contact,
						Address: provider.Address,
						Price:   provider.Price,
						Rating:  provider.Rating,
					})
				}

			}
			responseBody = append(responseBody, *currRequest)

		}

	}
	logger.Info("approved request fetched", nil)
	response.SuccessResponse(w, responseBody, "Approved requests fetched", http.StatusOK)

}

func (h *HouseholderController) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RequestID  string `json:"request_id" validate:"required"`
		ProviderID string `json:"provider_id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Invalid input", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	err := validate.Struct(request)
	if err != nil {
		logger.Error("Invalid request body", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}
	role := r.Context().Value("role").(string)
	var householderID string
	if role == "Admin" {
		householderID = r.URL.Query().Get("user_id")
		if householderID == "" {
			logger.Error("No query param", nil)
			response.ErrorResponse(w, http.StatusBadRequest, "user ID is required", 2001)
			return
		}
	} else if role == "Householder" {
		householderID = r.Context().Value("userID").(string)
	} else {
		logger.Error("Invalid role", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid role", 1007)
		return
	}
	// Call the approval function
	if err = h.householderService.ApproveServiceRequest(request.RequestID, request.ProviderID, householderID); err != nil {
		logger.Error(err.Error(), nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", 1006)
		return
	}

	logger.Info("Request approve successfully", nil)
	response.SuccessResponse(w, nil, "Request approve successfully", http.StatusOK)

}

func (h *HouseholderController) LeaveReview(w http.ResponseWriter, r *http.Request) {
	var reviewRequest struct {
		ServiceID  string  `json:"service_id" validate:"required"`
		ProviderID string  `json:"provider_id" validate:"required"`
		ReviewText string  `json:"review_text" validate:"required"`
		Rating     float64 `json:"rating" validate:"required"`
	}

	// Decode request body into reviewRequest struct
	if err := json.NewDecoder(r.Body).Decode(&reviewRequest); err != nil {
		logger.Error(err.Error(), nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Error decoding review request", 1001)
		return
	}
	err := validate.Struct(reviewRequest)
	if err != nil {
		logger.Error("Invalid request body", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}
	// Get user from context (assuming AuthMiddleware has set user info in the context)
	userID := r.Context().Value("userID").(string)

	// Validate rating input (between 1 and 5)
	if reviewRequest.Rating < 1 || reviewRequest.Rating > 5 {
		response.ErrorResponse(w, http.StatusBadRequest, "Rating should be between 1 and 5", 1001)
		return
	}

	// Call the householder service to add the review
	err = h.householderService.AddReview(reviewRequest.ProviderID, userID, reviewRequest.ServiceID, reviewRequest.ReviewText, reviewRequest.Rating)
	if err != nil {
		logger.Error(fmt.Sprintf("Error adding review %v", err), nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to submit review", 1006)
		return
	}

	// Successfully added the review
	logger.Info("Review added successfully", nil)
	response.SuccessResponse(w, nil, "Review added successfully", http.StatusOK)
}
