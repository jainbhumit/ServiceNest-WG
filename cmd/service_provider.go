// //go:build !test
// // +build !test
package main

//
//import (
//	"bufio"
//	"database/sql"
//	"fmt"
//	"github.com/fatih/color"
//	"os"
//	"serviceNest/interfaces"
//	"serviceNest/model"
//	"serviceNest/repository"
//	"serviceNest/service"
//	"serviceNest/util"
//	"time"
//)
//
//func serviceProviderDashboard(user *model.User, client *sql.DB) {
//	serviceRepo := repository.NewServiceRepository(client)
//	requestRepo := repository.NewServiceRequestRepository(client)
//	providerRepo := repository.NewServiceProviderRepository(client)
//
//	providerService := service.NewServiceProviderService(providerRepo, requestRepo, serviceRepo)
//	//provider := &model.ServiceProvider{
//	//	User:            *user,
//	//	ServicesOffered: []model.Service{},
//	//	Rating:          0.0,
//	//	Reviews:         []*model.Review{},
//	//	Availability:    true,
//	//	IsActive:        true,
//	//}
//	// Check if ServiceProvider already exists
//	if user == nil {
//		color.Red("User is nil. Cannot proceed with dashboard.")
//		return
//	}
//
//	provider, err := providerRepo.GetProviderByID(user.ID)
//	if err != nil {
//
//		// If not found, create a new ServiceProvider
//		provider = &model.ServiceProvider{
//			User:            *user,
//			ServicesOffered: []model.Service{},
//			Rating:          0.0,
//			Reviews:         []*model.Review{},
//			Availability:    true,
//			IsActive:        true,
//		}
//
//		// Save the new ServiceProvider to the repository_test
//		err = providerRepo.SaveServiceProvider(*provider)
//		if err != nil {
//			color.Red("Error saving new service provider: %v", err)
//			return
//		}
//	}
//	if !provider.IsActive {
//		color.Red("Service provider is deactivated by admin")
//		return
//	}
//	for {
//		color.Blue("1. View Profile")
//		color.Blue("2. Add Service")
//		color.Blue("3. View your Services")
//		color.Blue("4. Update Service")
//		color.Blue("5. Remove Service")
//		color.Blue("6. View and Accept Service Request")
//		color.Blue("7. Decline Service Request")
//		color.Blue("8. Update Availability")
//		color.Blue("9. View Approved Services")
//		color.Blue("10. View Reviews")
//
//		color.Blue("11. Exit")
//
//		var choice int
//		fmt.Scanln(&choice)
//
//		switch choice {
//		case 1:
//			viewProfile(user)
//		case 2:
//			addService(providerService, provider)
//		case 3:
//			viewProviderServices(providerService, provider)
//		case 4:
//			updateService(providerService, provider)
//		case 5:
//			removeService(providerService, provider)
//		case 6:
//			viewAndAcceptServiceRequest(providerService, provider)
//		case 7:
//			declineServiceRequest(providerService, provider)
//		case 8:
//			updateAvailability(providerService, provider)
//		case 9:
//			viewApprovedRequestsForProvider(providerService, provider.User.ID)
//		case 10:
//			viewReview(providerService, provider.User.ID)
//		case 11:
//			return
//		default:
//			color.Red("Invalid choice")
//		}
//	}
//}
//
//func addService(providerService interfaces.ServiceProviderService, provider *model.ServiceProvider) {
//	util.DisplayCategory()
//	var serviceName, category string
//	var price float64
//
//	fmt.Print("Enter service name: ")
//	fmt.Scanln(&serviceName)
//	fmt.Print("Enter description: ")
//	reader := bufio.NewReader(os.Stdin)
//	description, err := reader.ReadString('\n')
//	if err != nil {
//		color.Red("Error reading desciption")
//		return
//	}
//	fmt.Print("Enter category: ")
//	fmt.Scanln(&category)
//	fmt.Print("Enter price: ")
//	fmt.Scanln(&price)
//
//	service := model.Service{
//		ID:              util.GenerateUniqueID(),
//		Name:            serviceName,
//		Description:     description,
//		Price:           price,
//		Category:        category,
//		ProviderID:      provider.ID,
//		ProviderName:    provider.Name,
//		ProviderContact: provider.Contact,
//		ProviderAddress: provider.Address,
//		ProviderRating:  provider.Rating,
//	}
//
//	err = providerService.AddService(provider.ID, service)
//	if err != nil {
//		color.Red("Error adding service_test: %v", err)
//		return
//	}
//
//	color.Green("Service added successfully!")
//}
//
//func updateService(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
//	var serviceID, newName string
//	var newPrice float64
//
//	fmt.Print("Enter service ID to update: ")
//	fmt.Scanln(&serviceID)
//	fmt.Print("Enter new service name: ")
//	fmt.Scanln(&newName)
//	fmt.Print("Enter new description: ")
//	reader := bufio.NewReader(os.Stdin)
//	newDescription, err := reader.ReadString('\n')
//	if err != nil {
//		color.Red("Error reading desciption")
//		return
//	}
//	fmt.Print("Enter new price: ")
//	fmt.Scanln(&newPrice)
//
//	updatedService := model.Service{
//		ID:          serviceID,
//		Name:        newName,
//		Description: newDescription,
//		Price:       newPrice,
//	}
//
//	err = providerService.UpdateService(provider.ID, serviceID, updatedService)
//	if err != nil {
//		color.Red("Error updating service: %v", err)
//		return
//	}
//
//	color.Green("Service updated successfully!")
//}
//
//func removeService(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
//	var serviceID string
//
//	fmt.Print("Enter service ID to remove: ")
//	fmt.Scanln(&serviceID)
//
//	err := providerService.RemoveService(provider.ID, serviceID)
//	if err != nil {
//		color.Red("Error removing service: %v", err)
//		return
//	}
//
//	color.Green("Service removed successfully!")
//}
//
//func declineServiceRequest(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
//	var requestID string
//
//	fmt.Print("Enter service request ID to decline: ")
//	fmt.Scanln(&requestID)
//
//	err := providerService.DeclineServiceRequest(provider.ID, requestID)
//	if err != nil {
//		color.Red("Error declining service request: %v", err)
//		return
//	}
//
//	color.Green("Service request declined successfully!")
//}
//
//func updateAvailability(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
//	var available string
//
//	fmt.Print("Are you available? (yes/no): ")
//	fmt.Scanln(&available)
//
//	isAvailable := available == "yes"
//	err := providerService.UpdateAvailability(provider.ID, isAvailable)
//	if err != nil {
//		color.Red("Error updating availability: %v", err)
//		return
//	}
//
//	color.Green("Availability updated successfully!")
//}
//func viewProviderServices(serviceProviderService *service.ServiceProviderService, provider *model.ServiceProvider) {
//	services, err := serviceProviderService.ViewServices(provider.ID)
//	if err != nil {
//		color.Red("Error viewing services: %v", err)
//		return
//	}
//
//	if len(services) == 0 {
//		color.Cyan("You do not have any services offered.")
//		return
//	}
//
//	color.Cyan("Services Offered:")
//	for _, service := range services {
//		color.Cyan("-ID-%s %s: %s (Price: %.2f)", service.ID, service.Name, service.Description, service.Price)
//	}
//}
//func viewAndAcceptServiceRequest(providerService *service.ServiceProviderService, provider *model.ServiceProvider) {
//
//	// Fetch all service requests
//	serviceRequests, err := providerService.GetAllServiceRequests()
//	if err != nil {
//		color.Red("Error fetching service requests: %v", err)
//		return
//	}
//
//	// Filter and display only pending requests
//	var pendingRequests []model.ServiceRequest
//	for _, request := range serviceRequests {
//		if request.ApproveStatus == false && request.Status != "Cancelled" {
//			pendingRequests = append(pendingRequests, request)
//			color.Cyan("Request ID: %s, Service ID: %s", request.ID, request.ServiceID)
//		}
//	}
//
//	if len(pendingRequests) == 0 {
//		color.Yellow("No pending service requests available.")
//		return
//	}
//
//	var requestID string
//
//	fmt.Print("Enter the Service Request ID to view and accept: ")
//	fmt.Scanln(&requestID)
//
//	// Fetch the service request by ID
//	serviceRequest, err := providerService.GetServiceRequestByID(requestID)
//	if err != nil {
//		color.Red("Error fetching service request: %v", err)
//		return
//	}
//
//	// Display the details of the service_test request
//	color.Cyan("Service Request Details:")
//	color.Cyan("Request ID: %s", serviceRequest.ID)
//	color.Cyan("Householder ID: %v", serviceRequest.HouseholderID)
//	color.Cyan("Service Name: %s", serviceRequest.ServiceName)
//	color.Cyan("Service ID: %s", serviceRequest.ServiceID)
//	color.Cyan("Requested Time: %s", serviceRequest.RequestedTime.Format(time.RFC1123))
//	color.Cyan("Scheduled Time: %s", serviceRequest.ScheduledTime.Format(time.RFC1123))
//	//color.Cyan("Status: %s", serviceRequest.Status)
//
//	// Ask if the service provider wants to accept the request
//	var accept string
//	fmt.Print("Do you want to accept this request? (yes/no): ")
//	fmt.Scanln(&accept)
//
//	var estimatedPrice string
//	color.Cyan("Enter the price of request")
//	fmt.Scanln(&estimatedPrice)
//
//	if accept == "yes" {
//		// Accept the service_test request
//		err = providerService.AcceptServiceRequest(provider.ID, requestID, estimatedPrice)
//		if err != nil {
//			color.Red("Error accepting service request: %v", err)
//			return
//		}
//		color.Green("Service request accepted successfully!")
//	} else {
//		color.Yellow("Service request not accepted.")
//	}
//}
//func viewApprovedRequestsForProvider(serviceProviderService *service.ServiceProviderService, providerID string) {
//	// Call the service method to get approved requests
//	approvedRequests, err := serviceProviderService.ViewApprovedRequestsByProvider(providerID)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	// Display the approved requests
//	if len(approvedRequests) == 0 {
//		fmt.Println("No approved service requests found for this provider.")
//		return
//	}
//
//	fmt.Println("Approved Service Requests:")
//	for _, req := range approvedRequests {
//		fmt.Printf("\nRequest ID: %s\n", req.ID)
//		fmt.Printf("Service ID: %s\n", req.ServiceID)
//		fmt.Printf("Requested Time: %s\n", req.RequestedTime.Format("2006-01-02 15:04:05"))
//		fmt.Printf("Scheduled Time: %s\n", req.ScheduledTime.Format("2006-01-02 15:04:05"))
//		fmt.Printf("Status: %s\n", req.Status)
//		fmt.Printf("Householder Name: %s\n", req.HouseholderName)
//		fmt.Printf("Householder Address: %s\n", *req.HouseholderAddress)
//
//		// Display provider-specific details
//		for _, providerDetail := range req.ProviderDetails {
//			if providerDetail.ServiceProviderID == providerID {
//				fmt.Println("Provider Details:")
//				fmt.Printf("\tName: %s\n", providerDetail.Name)
//				fmt.Printf("\tContact: %s\n", providerDetail.Contact)
//				fmt.Printf("\tAddress: %s\n", providerDetail.Address)
//				fmt.Printf("\tPrice: %s\n", providerDetail.Price)
//				fmt.Printf("\tRating: %.2f\n", providerDetail.Rating)
//				fmt.Println("\tReviews:")
//				for _, review := range providerDetail.Reviews {
//					fmt.Printf("\t\tReview ID: %s\n", review.ID)
//					fmt.Printf("\t\tRating: %.2f\n", review.Rating)
//					fmt.Printf("\t\tComments: %s\n", review.Comments)
//					fmt.Printf("\t\tReview Date: %s\n", review.ReviewDate.Format("2006-01-02"))
//				}
//			}
//		}
//	}
//}
//
//func viewReview(serviceProviderService *service.ServiceProviderService, providerID string) {
//	reviews, err := serviceProviderService.GetReviews(providerID)
//	if err != nil {
//		color.Red("Error fetching reviews: %v", err)
//		return
//	}
//	for _, review := range reviews {
//		color.Cyan("ServiceID: %v ", review.ServiceID)
//		color.Cyan("Rating: %v ", review.Rating)
//		color.Cyan("Comments: %v", review.Comments)
//		color.Cyan("Date: %v", review.ReviewDate)
//		color.Cyan("----------------------------------------")
//
//	}
//}
