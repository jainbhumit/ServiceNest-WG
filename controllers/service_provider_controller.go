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
	"serviceNest/util"
	"time"
)

var GenerateUniqueID = util.GenerateUniqueID

type ServiceProviderController struct {
	serviceProviderService interfaces.ServiceProviderService
}

func NewServiceProviderController(serviceProviderService interfaces.ServiceProviderService) *ServiceProviderController {
	return &ServiceProviderController{serviceProviderService: serviceProviderService}
}

func (s *ServiceProviderController) AddService(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string  `json:"name" validate:"required"`
		Description string  `json:"description" validate:"required"`
		Price       float64 `json:"price" validate:"required"`
		Category    string  `json:"category" validate:"required"`
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
	providerID := r.Context().Value("userID").(string)

	newService := &model.Service{
		ID:          GenerateUniqueID(),
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Category:    request.Category,
		ProviderID:  providerID,
	}

	err = s.serviceProviderService.AddService(providerID, *newService)
	if err != nil {
		logger.Error("Error adding service", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), 1003)

		//http.Error(w, "Error adding service", http.StatusInternalServerError)
		return
	}
	var serviceID struct {
		ID string `json:"service_id"`
	}
	serviceID.ID = newService.ID
	response.SuccessResponse(w, serviceID, "Service added successfully", http.StatusOK)

}

func (s *ServiceProviderController) ViewServices(w http.ResponseWriter, r *http.Request) {
	providerID := r.Context().Value("userID").(string)

	services, err := s.serviceProviderService.ViewServices(providerID)
	if err != nil {
		logger.Error(err.Error(), nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "Error fetching services", 1006)
		//http.Error(w, "Error fetching services", http.StatusInternalServerError)
		return
	}
	if len(services) == 0 {
		response.SuccessResponse(w, nil, "Don't have a service offered", http.StatusOK)
		return
	}

	logger.Info("Service fetched successfully", nil)
	response.SuccessResponse(w, services, "service fetch successfully", http.StatusOK)

}
func (h *ServiceProviderController) UpdateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["service_id"]
	var request struct {
		Name        string  `json:"name" validate:"required"`
		Description string  `json:"description" validate:"required"`
		Price       float64 `json:"price" validate:"required"`
		Category    string  `json:"category" validate:"required"`
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
	providerID := r.Context().Value("userID").(string)

	updatedService := &model.Service{
		ID:          serviceID,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Category:    request.Category,
		ProviderID:  providerID,
	}

	err = h.serviceProviderService.UpdateService(providerID, serviceID, *updatedService)
	if err != nil {
		logger.Error("Error updating service", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), 1003)
		//http.Error(w, "Error updating service", http.StatusInternalServerError)
		return
	}
	logger.Info("update successfully", nil)
	response.SuccessResponse(w, nil, "Service updated successfully", http.StatusOK)
}

func (s *ServiceProviderController) RemoveService(w http.ResponseWriter, r *http.Request) {
	providerID := r.Context().Value("userID").(string)

	serviceID := mux.Vars(r)["service_id"]

	err := s.serviceProviderService.RemoveService(providerID, serviceID)
	if err != nil {
		logger.Error("Error removing service", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), 1008)
		//http.Error(w, "Error removing service", http.StatusInternalServerError)
		return
	}
	response.SuccessResponse(w, nil, "Service removed successfully", http.StatusOK)
	//json.NewEncoder(w).Encode("Service removed successfully")
}

func (s *ServiceProviderController) ViewServiceRequest(w http.ResponseWriter, r *http.Request) {

	serviceRequests, err := s.serviceProviderService.GetAllServiceRequests()
	if err != nil {
		logger.Error(err.Error(), nil)
		response.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error fetching request %v", err), 1006)
		return
	}

	type responseStruct struct {
		ID                 string    `json:"request_id"`
		ServiceName        string    `json:"service_name,omitempty"`
		ServiceID          string    `json:"service_id" `
		RequestedTime      time.Time `json:"requested_time"`
		ScheduledTime      time.Time `json:"scheduled_time"`
		HouseholderAddress string    `json:"address"`
	}
	responseBody := make([]responseStruct, 0)
	// Filter and display only pending requests

	for _, request := range serviceRequests {
		if request.ApproveStatus == false && request.Status != "Cancelled" {
			responseBody = append(responseBody, responseStruct{
				ID:                 request.ID,
				ServiceName:        request.ServiceName,
				ServiceID:          request.ServiceID,
				RequestedTime:      request.RequestedTime,
				ScheduledTime:      request.ScheduledTime,
				HouseholderAddress: *request.HouseholderAddress,
			})
		}
	}

	if len(responseBody) == 0 {
		logger.Info("No service requests found", nil)
		response.SuccessResponse(w, nil, "No pending service requests available", http.StatusOK)
		//color.Yellow("No pending service requests available.")
		return
	}
	logger.Info("Service request fetched successfully", nil)
	response.SuccessResponse(w, responseBody, "Service request fetched successfully", http.StatusOK)

}

func (s *ServiceProviderController) AcceptServiceRequest(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID             string `json:"request_id" validate:"required"`
		EstimatedPrice string `json:"price" validate:"required"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid body", 1001)
		//http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = validate.Struct(request)
	if err != nil {
		logger.Error("Invalid request body", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}

	providerID := r.Context().Value("userID").(string)
	err = s.serviceProviderService.AcceptServiceRequest(providerID, request.ID, request.EstimatedPrice)
	if err != nil {
		logger.Error(err.Error(), nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "Error accepting service request", 1008)
		//http.Error(w, "Error accepting service request", http.StatusInternalServerError)
		return
	}

	logger.Info("Request accept successfully", nil)
	response.SuccessResponse(w, nil, "Request accept successfully", http.StatusOK)

}

func (s *ServiceProviderController) ViewApprovedRequests(w http.ResponseWriter, r *http.Request) {
	providerID := r.Context().Value("userID").(string)

	approvedRequests, err := s.serviceProviderService.ViewApprovedRequestsByProvider(providerID)

	if err != nil {
		logger.Error(err.Error(), nil)
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), 1008)

		//http.Error(w, "Error fetching approved requests", http.StatusInternalServerError)
		return
	}
	type responseStruct struct {
		ID          string `json:"request_id"`
		ServiceName string `json:"service_name,omitempty"`

		HouseholderId      string                         `json:"householder_id"`
		HouseholderName    string                         `json:"householder_name"`
		HouseholderAddress string                         `json:"householder_address"`
		ApproveStatus      bool                           `json:"approve_status"`
		ServiceID          string                         `json:"service_id" `
		RequestedTime      time.Time                      `json:"requested_time"`
		ScheduledTime      time.Time                      `json:"scheduled_time"`
		Status             string                         `json:"status"`
		ProviderDetails    []model.ServiceProviderDetails `json:"provider_details,omitempty" bson:"providerDetails,omitempty"`
	}
	responseBody := make([]responseStruct, 0)
	for _, request := range approvedRequests {
		if request.ApproveStatus {
			currRequest := &responseStruct{
				ID:                 request.ID,
				ServiceName:        request.ServiceName,
				ServiceID:          request.ServiceID,
				RequestedTime:      request.RequestedTime,
				ScheduledTime:      request.ScheduledTime,
				Status:             request.Status,
				HouseholderId:      *request.HouseholderID,
				HouseholderAddress: *request.HouseholderAddress,
				HouseholderName:    request.HouseholderName,
			}
			responseBody = append(responseBody, *currRequest)

		}

	}
	response.SuccessResponse(w, responseBody, "Approve requests fetched successfully", http.StatusOK)

	//json.NewEncoder(w).Encode(approvedRequests)
}

func (s *ServiceProviderController) ViewReviews(w http.ResponseWriter, r *http.Request) {
	providerID := r.Context().Value("userID").(string)

	// Call the service to get the reviews
	reviews, err := s.serviceProviderService.GetReviews(providerID)
	if err != nil {
		logger.Error(err.Error(), nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch reviews", 1003)

		//http.Error(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	// Send reviews as JSON response
	response.SuccessResponse(w, reviews, "Reviews fetched successfully", http.StatusOK)
}
