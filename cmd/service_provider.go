package main

import (
	"fmt"
	"github.com/fatih/color"
	"serviceNest/model"
	"serviceNest/repository"
	"serviceNest/service"
	"serviceNest/util"
	"time"
)

func serviceProviderDashboard(user *model.User) {
	serviceRepo := repository.NewServiceRepository("services.json")
	requestRepo := repository.NewServiceRequestRepository("service_requests.json")
	providerRepo := repository.NewServiceProviderRepository("service_providers.json")

	providerService := service.NewServiceProviderService(providerRepo, requestRepo, serviceRepo)
	provider := &model.ServiceProvider{
		User:            *user,
		ServicesOffered: []model.Service{},
		Rating:          0.0,
		Reviews:         []*model.Review{},
		Availability:    true,
	}
	// Check if ServiceProvider already exists
	provider, err := providerRepo.GetProviderByID(user.ID)
	if err != nil {
		// If not found, create a new ServiceProvider
		provider = &model.ServiceProvider{
			User:            *user,
			ServicesOffered: []model.Service{},
			Rating:          0.0,
			Reviews:         []*model.Review{},
			Availability:    true,
		}

		// Save the new ServiceProvider to the repository
		err = providerRepo.SaveServiceProvider(*provider)
		if err != nil {
			color.Red("Error saving new service provider: %v", err)
			return
		}
	}
	for {
		color.Blue("1. View Profile")
		color.Blue("2. Add Service")
		color.Blue("3. View your Services")
		color.Blue("4. Update Service")
		color.Blue("5. Remove Service")
		color.Blue("6. View and Accept Service Request")
		color.Blue("7. Decline Service Request")
		color.Blue("8. Update Availability")
		color.Blue("9. Exit")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			viewProfile(user)
		case 2:
			addService(providerService, provider)
		case 3:
			viewProviderServices(providerService, provider)
		case 4:
			updateService(providerService, provider)
		case 5:
			removeService(providerService, provider)
		case 6:
			viewAndAcceptServiceRequest(providerService, provider)
		case 7:
			declineServiceRequest(providerService, provider)
		case 8:
			updateAvailability(providerService, provider)
		case 9:
			return
		default:
			color.Red("Invalid choice")
		}
	}
}

func addService(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
	var serviceName, description string
	var price float64

	fmt.Print("Enter service name: ")
	fmt.Scanln(&serviceName)
	fmt.Print("Enter description: ")
	fmt.Scanln(&description)
	fmt.Print("Enter price: ")
	fmt.Scanln(&price)

	service := model.Service{
		ID:          util.GenerateUniqueID(), // Implement this function for generating unique IDs
		Name:        serviceName,
		Description: description,
		Price:       price,
	}

	err := providerService.AddService(provider.ID, service)
	if err != nil {
		color.Red("Error adding service: %v", err)
		return
	}

	color.Green("Service added successfully!")
}

func updateService(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
	var serviceID, newName, newDescription string
	var newPrice float64

	fmt.Print("Enter service ID to update: ")
	fmt.Scanln(&serviceID)
	fmt.Print("Enter new service name: ")
	fmt.Scanln(&newName)
	fmt.Print("Enter new description: ")
	fmt.Scanln(&newDescription)
	fmt.Print("Enter new price: ")
	fmt.Scanln(&newPrice)

	updatedService := model.Service{
		ID:          serviceID,
		Name:        newName,
		Description: newDescription,
		Price:       newPrice,
	}

	err := providerService.UpdateService(provider.ID, serviceID, updatedService)
	if err != nil {
		color.Red("Error updating service: %v", err)
		return
	}

	color.Green("Service updated successfully!")
}

func removeService(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
	var serviceID string

	fmt.Print("Enter service ID to remove: ")
	fmt.Scanln(&serviceID)

	err := providerService.RemoveService(provider.ID, serviceID)
	if err != nil {
		color.Red("Error removing service: %v", err)
		return
	}

	color.Green("Service removed successfully!")
}

func acceptServiceRequest(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
	var requestID string

	fmt.Print("Enter service request ID to accept: ")
	fmt.Scanln(&requestID)

	err := providerService.AcceptServiceRequest(provider.ID, requestID)
	if err != nil {
		color.Red("Error accepting service request: %v", err)
		return
	}

	color.Green("Service request accepted successfully!")
}

func declineServiceRequest(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
	var requestID string

	fmt.Print("Enter service request ID to decline: ")
	fmt.Scanln(&requestID)

	err := providerService.DeclineServiceRequest(provider.ID, requestID)
	if err != nil {
		color.Red("Error declining service request: %v", err)
		return
	}

	color.Green("Service request declined successfully!")
}

func updateAvailability(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
	var available string

	fmt.Print("Are you available? (yes/no): ")
	fmt.Scanln(&available)

	isAvailable := available == "yes"
	err := providerService.UpdateAvailability(provider.ID, isAvailable)
	if err != nil {
		color.Red("Error updating availability: %v", err)
		return
	}

	color.Green("Availability updated successfully!")
}
func viewProviderServices(serviceProviderService *service.ServiceProviderService, provider *model.ServiceProvider) {
	services, err := serviceProviderService.ViewServices(provider.ID)
	if err != nil {
		color.Red("Error viewing services: %v", err)
		return
	}

	if len(services) == 0 {
		color.Cyan("You do not have any services offered.")
		return
	}

	color.Cyan("Services Offered:")
	for _, service := range services {
		color.Cyan("-ID-%s %s: %s (Price: %.2f)", service.ID, service.Name, service.Description, service.Price)
	}
}
func viewAndAcceptServiceRequest(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
	// Fetch all service requests
	serviceRequests, err := providerService.GetAllServiceRequests()
	if err != nil {
		color.Red("Error fetching service requests: %v", err)
		return
	}

	// Filter and display only pending requests
	var pendingRequests []model.ServiceRequest
	for _, request := range serviceRequests {
		if request.Status == "Pending" {
			pendingRequests = append(pendingRequests, request)
			color.Cyan("Request ID: %s, Service ID: %s", request.ID, request.ServiceID)
		}
	}

	if len(pendingRequests) == 0 {
		color.Yellow("No pending service requests available.")
		return
	}

	var requestID string

	fmt.Print("Enter the Service Request ID to view and accept: ")
	fmt.Scanln(&requestID)

	// Fetch the service request by ID
	serviceRequest, err := providerService.GetServiceRequestByID(requestID)
	if err != nil {
		color.Red("Error fetching service request: %v", err)
		return
	}

	// Display the details of the service request
	color.Cyan("Service Request Details:")
	color.Cyan("Request ID: %s", serviceRequest.ID)
	color.Cyan("Householder ID: %s", *serviceRequest.HouseholderID)
	color.Cyan("Service ID: %s", serviceRequest.ServiceID)
	color.Cyan("Requested Time: %s", serviceRequest.RequestedTime.Format(time.RFC1123))
	color.Cyan("Scheduled Time: %s", serviceRequest.ScheduledTime.Format(time.RFC1123))
	color.Cyan("Status: %s", serviceRequest.Status)

	// Ask if the service provider wants to accept the request
	var accept string
	fmt.Print("Do you want to accept this request? (yes/no): ")
	fmt.Scanln(&accept)

	if accept == "yes" {
		// Accept the service request
		err = providerService.AcceptServiceRequest(provider.ID, requestID)
		if err != nil {
			color.Red("Error accepting service request: %v", err)
			return
		}
		color.Green("Service request accepted successfully!")
	} else {
		color.Yellow("Service request not accepted.")
	}
}
