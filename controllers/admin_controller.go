package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"serviceNest/interfaces"
	"serviceNest/logger"
	"serviceNest/response"
	"serviceNest/util"
)

type AdminController struct {
	adminService interfaces.AdminService
}

// NewAdminController initializes a new AdminController with the given service
func NewAdminController(adminService interfaces.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

// ManageServices handles the services management functionality
func (a *AdminController) ViewAllService(w http.ResponseWriter, r *http.Request) {
	// Fetch all services (without limit and offset in service/repository)
	limit, offset := util.GetPaginationParams(r)
	services, err := a.adminService.GetAllService(limit, offset)
	if err != nil {
		logger.Error("error fetching services", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "error fetching services", 1003)
		return
	}

	logger.Info("All services fetched successfully", nil)
	response.SuccessResponse(w, services, "All available services", http.StatusOK)
}

// DeleteService allows the admin to delete a service
func (a *AdminController) DeleteService(w http.ResponseWriter, r *http.Request) {
	serviceID := mux.Vars(r)["serviceID"]

	err := a.adminService.DeleteService(serviceID)
	if err != nil {
		logger.Error("error deleting service", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "Error deleting service", 1006)
		return
	}

	response.SuccessResponse(w, nil, "Service deleted successfully", http.StatusOK)
}

// ViewReports allows the admin to view reports
func (a *AdminController) ViewReports(w http.ResponseWriter, r *http.Request) {
	// Fetch all reports (without limit and offset in service/repository)
	limit, offset := util.GetPaginationParams(r)
	reports, err := a.adminService.ViewReports(limit, offset)
	if err != nil {
		logger.Error("error fetching reports", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "Error generating reports", 1006)
		return
	}

	paginatedReports := util.ApplyPagination(reports, limit, offset)

	response.SuccessResponse(w, paginatedReports, "Reports fetched successfully", http.StatusOK)
}

// DeactivateUserAccount allows the admin to deactivate a user account
func (a *AdminController) DeactivateUserAccount(w http.ResponseWriter, r *http.Request) {
	providerID := mux.Vars(r)["providerID"]

	err := a.adminService.DeactivateAccount(providerID)
	if err != nil {
		logger.Error("error deactivating account", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "Error deactivating account", 1006)
		return
	}

	response.SuccessResponse(w, nil, "Account deactivated successfully", http.StatusOK)
}

func (a *AdminController) AddService(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string `json:"category_name" validate:"required"`
		Description string `json:"description" validate:"required"`
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

	err = a.adminService.AddService(request.Name, request.Description)
	if err != nil {
		logger.Error("error Adding service", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), 1006)
		return
	}

	response.SuccessResponse(w, nil, "Service added successfully", http.StatusOK)
}

func (a *AdminController) ViewUserDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userEmail := vars["userEmail"]

	user, err := a.adminService.GetUserByEmail(userEmail)

	if err != nil {
		logger.Error("error fetching services", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), 1003)
		return
	}

	type responseBody struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
	}
	userDetail := responseBody{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Address: user.Address,
	}

	logger.Info("All services fetched successfully", nil)
	response.SuccessResponse(w, userDetail, "User fetch successfully", http.StatusOK)
}
