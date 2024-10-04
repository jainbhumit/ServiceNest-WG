package controllers

import (
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
	services, err := a.adminService.GetAllService()
	if err != nil {
		logger.Error("error fetching services", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "error fetching services", 1003)
		return
	}

	// Extract limit and offset from query parameters
	limit, offset := util.GetPaginationParams(r)

	// Apply pagination on the result (slice the result set)
	paginatedServices := util.ApplyPagination(services, limit, offset)

	logger.Info("All services fetched successfully", nil)
	response.SuccessResponse(w, paginatedServices, "All available services", http.StatusOK)
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
	reports, err := a.adminService.ViewReports()
	if err != nil {
		logger.Error("error fetching reports", nil)
		response.ErrorResponse(w, http.StatusInternalServerError, "Error generating reports", 1006)
		return
	}

	// Extract limit and offset from query parameters
	limit, offset := util.GetPaginationParams(r)

	// Apply pagination on the result (slice the result set)
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
