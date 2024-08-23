package main

import (
	"fmt"
	"github.com/fatih/color"
	"serviceNest/model"
	"serviceNest/repository"
	"serviceNest/service"
)

// AdminDashboard is the main dashboard for admin actions
func adminDashboard(admin *model.Admin) {
	adminRepo := repository.NewAdminRepository("admins.json")
	serviceRepo := repository.NewServiceRepository("services.json")
	userRepo := repository.NewUserRepository("users.json")
	serviceRequestRepo := repository.NewServiceRequestRepository("service_requests.json")
	providerRepo := repository.NewServiceProviderRepository("service_providers.json")

	adminService := service.NewAdminService(adminRepo, serviceRepo, serviceRequestRepo, userRepo, providerRepo)

	for {
		color.Blue("Admin Dashboard")
		color.Blue("1. Manage Services")
		color.Blue("2. View Reports")
		color.Blue("3. Deactivate User Account")
		//color.Blue("4. Add Service Area")
		//color.Blue("5. Remove Service Area")
		//color.Blue("6. Update Service Area")
		color.Blue("4. Exit")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			manageServices(adminService)
		case 2:
			viewReports(adminService)
		case 3:
			deactivateUserAccount(adminService)
		case 4:
			return

		default:
			color.Red("Invalid choice")
		}
	}
}

// ManageServices handles the services management functionality
func manageServices(adminService *service.AdminService) {
	color.Blue("Manage Services")
	color.Blue("1. View All Services")
	//color.Blue("2. Update Service")
	color.Blue("2. Delete Service")
	color.Blue("3. Back to Dashboard")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		viewAllServices(adminService)
	case 2:
		deleteService(adminService)
	case 4:
		return
	default:
		color.Red("Invalid choice")
	}
}

// ViewAllServices displays all the services available
func viewAllServices(adminService *service.AdminService) {
	services, err := adminService.GetAllService()
	if err != nil {
		color.Red("Error retrieving services: %v", err)
		return
	}

	for _, svc := range services {
		color.Cyan("Service ID: %s, Name: %s, Description: %s, Price: %.2f", svc.ID, svc.Name, svc.Description, svc.Price)
	}
}

// DeleteService allows the admin to delete a service
func deleteService(adminService *service.AdminService) {
	var serviceID string
	fmt.Print("Enter Service ID to delete: ")
	fmt.Scanln(&serviceID)

	err := adminService.DeleteService(serviceID)
	if err != nil {
		color.Red("Error deleting service: %v", err)
	} else {
		color.Green("Service deleted successfully")
	}
}

// ViewReports allows the admin to view various reports
func viewReports(adminService *service.AdminService) {
	color.Blue("View Reports")
	reports, err := adminService.ViewReports()
	if err != nil {
		color.Red("Error generating reports: %v", err)
		return
	}

	for _, report := range reports {
		color.Cyan("Report ID: %v, Details: %v", report.ID, report.ProviderDetails)
	}
}

// DeactivateUserAccount allows the admin to deactivate a user account
func deactivateUserAccount(adminService *service.AdminService) {
	var userID string
	fmt.Print("Enter User ID to deactivate: ")
	fmt.Scanln(&userID)

	err := adminService.DeactivateAccount(userID)
	if err != nil {
		color.Red("Error deactivating account: %v", err)
	} else {
		color.Green("Account deactivated successfully")
	}
}
