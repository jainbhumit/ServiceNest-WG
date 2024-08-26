package main

import (
	"fmt"
	"github.com/fatih/color"
	"serviceNest/model"
	"serviceNest/repository"
	"serviceNest/service"
	"serviceNest/util"
	"strings"
	"time"
)

// ViewProfile allows the user to view their profile details
func updateProfile(user *model.User) {
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)

	userID := user.ID
	var choice int
	color.Blue("For Update Email press 1")
	color.Blue("For Update Password press 2")
	color.Blue("For Update Contact press 3")
	color.Blue("For Update Address press 4")
	color.Blue("For going back press 5")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		var newEmail string
		for {
			fmt.Print("Enter Email: ")
			fmt.Scanln(&newEmail)
			newEmail = strings.TrimSpace(newEmail)
			err := userService.UpdateUser(userID, &newEmail, &user.Password, &user.Address, &user.Contact)
			if err != nil {
				color.Red("%v", err)
			} else {
				fmt.Println("User updated successfully!")
			}
			break
		}
	case 2:
		var newPassword string
		for {
			fmt.Print("Enter Password: ")
			fmt.Scanln(&newPassword)
			newPassword = strings.TrimSpace(newPassword)
			err := userService.UpdateUser(userID, &user.Email, &newPassword, &user.Address, &user.Contact)
			if err != nil {
				color.Red("%v", err)
			} else {
				fmt.Println("User updated successfully!")
			}
			break

		}
	case 3:
		var newPhone string
		for {
			fmt.Print("Enter Contact : ")
			fmt.Scanln(&newPhone)
			newPhone = strings.TrimSpace(newPhone)
			err := userService.UpdateUser(userID, &user.Email, &user.Password, &user.Address, &newPhone)
			if err != nil {
				color.Red("%v", err)
			} else {
				fmt.Println("User updated successfully!")
			}
			break
		}
	case 4:
		var newAddress string
		fmt.Print("Enter Address: ")
		fmt.Scanln(&newAddress)
		newAddress = strings.TrimSpace(newAddress)
		err := userService.UpdateUser(userID, &user.Email, &user.Password, &newAddress, &user.Contact)
		if err != nil {
			color.Red("%v", err)
		} else {
			fmt.Println("User updated successfully!")

		}

	case 5:
		return
	default:
		color.Red("Invalid choice")
	}

}
func viewProfile(user *model.User) {
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	currUser, err := userService.ViewProfileByID(user.ID)
	if err != nil {
		color.Red("%v", err)
	}
	color.Cyan("User Name: %s\n", currUser.Name)
	color.Cyan("User Email: %s\n", currUser.Email)
	color.Cyan("User Address: %s\n", currUser.Address)
	color.Cyan("User Contact: %s\n", currUser.Contact)
	color.Cyan("User Role: %s\n", currUser.Role)
	for {
		color.Blue("For updateProfile press 1")
		color.Blue("For going back press 2")
		var choice int
		fmt.Scanln(&choice)
		if choice == 1 {
			updateProfile(currUser)
		} else {
			break
		}
	}

}

// view services provide by provider
func viewServices(householderService *service.HouseholderService) {
	//color.Blue("Available Service are: ")
	//services, err := householderService.GetAvailableServices()
	//if err != nil {
	//	color.Red("%v", err)
	//}
	//for _, service := range services {
	//	fmt.Println("Name: ", service.Name)
	//	fmt.Println("Description: ", service.Description)
	//	fmt.Println("Price: ", service.Price)
	//	fmt.Println("-------------------------------------")
	//	//color.Green(service.Name)
	//	//color.Green(service.Description)
	//	//color.Green("%v", service.Price)
	//
	//}
	color.Blue("View Services by Category:")
	var category string
	fmt.Print("Enter the category: ")
	fmt.Scanln(&category)

	services, err := householderService.GetServicesByCategory(category)
	if err != nil {
		color.Red("Error fetching services: %v", err)
		return
	}

	if len(services) == 0 {
		color.Cyan("No services found in this category.")
		return
	}

	color.Cyan("Available Services in %s:", category)
	for _, service := range services {
		color.Cyan("Service ID: %s, Name: %s, Description: %s, Price: %.2f, Provider: %s",
			service.ID, service.Name, service.Description, service.Price, service.ProviderName)
		color.Cyan("Provider Contact: %s, Provider Address: %s, Provider Rating: %.2f",
			service.ProviderContact, service.ProviderAddress, service.ProviderRating)
	}
}

// SearchService allows the householder to search for available service providers
func searchService(householderService *service.HouseholderService, householder *model.Householder) {
	util.DisplayCategory()
	//var serviceType string
	//fmt.Print("Enter the type of service you're looking for: ")
	//fmt.Scanln(&serviceType)
	//
	//providers, err := householderService.SearchService(householder, serviceType)
	//if err != nil {
	//	color.Red("Error searching for service: %v", err)
	//	return
	//}
	//
	//color.Cyan("Found %d providers for %s:", len(providers), serviceType)
	//for _, provider := range providers {
	//	color.Cyan("- %s (%s)", provider.Name, provider.Contact)
	//}
	var category string
	fmt.Print("Enter the category of services you're looking for: ")
	fmt.Scanln(&category)

	services, err := householderService.GetServicesByCategory(category)
	if err != nil {
		color.Red("Error fetching services: %v", err)
		return
	}

	if len(services) == 0 {
		color.Cyan("No services found in this category.")
		return
	}

	color.Cyan("Available Services:")
	for _, service := range services {
		color.Cyan("Service ID: %s, Name: %s, Description: %s, Price: %.2f, Provider: %s",
			service.ID, service.Name, service.Description, service.Price, service.ProviderName)
		color.Cyan("Provider Contact: %s, Provider Address: %s, Provider Rating: %.2f",
			service.ProviderContact, service.ProviderAddress, service.ProviderRating)
	}
}

// RequestService allows the householder to request a specific service
func requestService(householderService *service.HouseholderService, user *model.User) {
	util.DisplayCategory()
	var serviceType string
	fmt.Print("Enter the type of service you want to request: ")
	fmt.Scanln(&serviceType)

	requestID, err := householderService.RequestService(user.ID, serviceType)
	if err != nil {
		color.Red("Error requesting service: %v", err)
		return
	}

	color.Green("Service requested successfully! Your request ID is %s", requestID)
}

// ViewBookingHistory allows the householder to view their booking history
func viewBookingHistory(householderService *service.HouseholderService, user *model.User) {
	history, err := householderService.ViewBookingHistory(user.ID)
	if err != nil {
		color.Red("Error viewing booking history: %v", err)
		return
	}

	if len(history) == 0 {
		color.Cyan("No booking history found.")
		return
	}

	color.Cyan("Booking History:")
	for _, booking := range history {
		color.Cyan("- %s: %s (%s)", booking.ID, booking.ServiceID, booking.Status)
	}
}

// LeaveReview allows the householder to leave a review for a service provider
func leaveReview(householderService *service.HouseholderService, user *model.User) {
	var serviceID, reviewText string
	var rating float64

	fmt.Print("Enter the Service ID you want to review: ")
	fmt.Scanln(&serviceID)
	fmt.Print("Enter your review: ")
	fmt.Scanln(&reviewText)
	fmt.Print("Enter your rating (1-5): ")
	fmt.Scanln(&rating)

	err := householderService.AddReview(user.ID, serviceID, reviewText, rating)
	if err != nil {
		color.Red("Error submitting review: %v", err)
		return
	}

	color.Green("Review submitted successfully!")
}

func cancelServiceRequest(householderService *service.HouseholderService) {
	var requestID string
	fmt.Print("Enter the Service Request ID you want to cancel: ")
	fmt.Scanln(&requestID)

	err := householderService.CancelServiceRequest(requestID)
	if err != nil {
		color.Red("Error canceling service request: %v", err)
		return
	}

	color.Green("Service request %s has been successfully canceled.", requestID)
}

func rescheduleServiceRequest(householderService *service.HouseholderService) {
	var requestID string
	fmt.Print("Enter the Service Request ID you want to reschedule: ")
	fmt.Scanln(&requestID)

	var newTimeStr string
	fmt.Print("Enter the new scheduled time (YYYY-MM-DD HH:MM): ")
	fmt.Scanln(&newTimeStr)
	newTime, err := time.Parse("2006-01-02 15:04", newTimeStr)
	if err != nil {
		color.Red("Error parsing new time: %v", err)
		return
	}

	err = householderService.RescheduleServiceRequest(requestID, newTime)
	if err != nil {
		color.Red("Error rescheduling service request: %v", err)
		return
	}

	color.Green("Service request %s has been successfully rescheduled.", requestID)
}

//	func viewServiceRequestStatus(householderService *service.HouseholderService) {
//		var requestID string
//		fmt.Print("Enter the Service Request ID: ")
//		fmt.Scanln(&requestID)
//
//		status, err := householderService.ViewServiceRequestStatus(requestID)
//		if err != nil {
//			color.Red("Error viewing service request status: %v", err)
//			return
//		}
//
//		color.Cyan("Status of service request %s: %s", requestID, status)
//	}
func viewStatus(householderService *service.HouseholderService, householder *model.Householder) {
	// Fetch all service requests for the householder
	requests, err := householderService.ViewStatus(householderService, householder)
	if err != nil {
		color.Red("Error viewing status: %v", err)
	}

	if len(requests) == 0 {
		color.Cyan("You have no service requests.")
		return
	}

	for _, request := range requests {
		color.Cyan("Request ID: %s, Service ID: %s, Status: %s", request.ID, request.ServiceID, request.Status)
		if request.Status == "Accepted" && request.ProviderDetails != nil {
			color.Green("ServiceProvider Details:")
			color.Green("Name: %s", request.ProviderDetails.Name)
			color.Green("Contact: %s", request.ProviderDetails.Contact)
			color.Green("Address: %s", request.ProviderDetails.Address)
			color.Green("Rating: %v", request.ProviderDetails.Rating)
			color.Green("Review for : %v", request.ServiceID)

			if len(request.ProviderDetails.Reviews) == 0 {
				color.Cyan("Provider has no Rivews.")
				continue
			}
			var reviewPresent bool
			for _, review := range request.ProviderDetails.Reviews {
				if review.ServiceID == request.ServiceID {
					reviewPresent = true
					color.Green("Rate[1-5]: %v", review.Rating)
					color.Green("Comment: %v", review.Comments)
					color.Green("Date: %v", review.ReviewDate)
				}
			}
			if !reviewPresent {
				color.Cyan("Provider has no Rivews for servive %v.", request.ServiceID)

			}
		}
		fmt.Println()

	}
	for {
		var choice string
		fmt.Println("1.Cancel any Accepted request")
		fmt.Println("2.Going back")

		fmt.Scanln(&choice)
		switch choice {
		case "1":
			cancelAcceptedServiceRequest(householderService, householder.ID)
		case "2":
			return
		default:
			color.Red("Invalid choice")
		}
	}

}
func cancelAcceptedServiceRequest(householderService *service.HouseholderService, householderID string) {
	var requestID string
	fmt.Print("Enter the Service Request ID you want to cancel: ")
	fmt.Scanln(&requestID)

	err := householderService.CancelAcceptedRequest(requestID, householderID)
	if err != nil {
		color.Red("Error canceling service request: %v", err)
		return
	}

	color.Green("Service request %s has been successfully canceled.", requestID)
}

// HouseholderDashboard is the main dashboard for householder actions
func householderDashboard(user *model.User) {
	// Initialize repositories and services
	householderRepo := repository.NewHouseholderRepository()
	serviceRequestRepo := repository.NewServiceRequestRepository()
	serviceProviderRepo := repository.NewServiceProviderRepository("service_providers.json")
	serviceRepo := repository.NewServiceRepository(nil)
	householderService := service.NewHouseholderService(householderRepo, serviceProviderRepo, serviceRepo, serviceRequestRepo)

	// Convert the User to a Householder
	householder := &model.Householder{
		User: *user,
	}
	for {
		color.Blue("-------------Householder Dashboard--------------")
		color.Blue("1. View Profile")
		color.Blue("2  View Services")
		color.Blue("3. Search Service")
		color.Blue("4. Request Service")
		color.Blue("5. View Booking History")
		color.Blue("6. Leave Review")
		color.Blue("7. Cancel Service Request")
		color.Blue("8. Reschedule Service Request")
		color.Blue("9. View Service Request Status")
		color.Blue("10.Exit")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			viewProfile(user)
		case 2:
			viewServices(householderService)
		case 3:
			searchService(householderService, householder)
		case 4:
			requestService(householderService, user)
		case 5:
			viewBookingHistory(householderService, user)
		case 6:
			leaveReview(householderService, user)
		case 7:
			cancelServiceRequest(householderService)
		case 8:
			rescheduleServiceRequest(householderService)
		case 9:
			viewStatus(householderService, householder)
		case 10:
			return
		default:
			color.Red("Invalid choice")
		}
	}
}
