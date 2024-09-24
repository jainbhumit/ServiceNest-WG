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
//	"strings"
//	"time"
//)
//
//// ViewProfile allows the user to view their profile details
//func updateProfile(user *model.User) {
//	userRepo := repository.NewUserRepository(nil)
//	userService := service.NewUserService(userRepo)
//
//	userID := user.ID
//	var choice int
//	color.Blue("For Update Email press 1")
//	color.Blue("For Update Password press 2")
//	color.Blue("For Update Contact press 3")
//	color.Blue("For Update Address press 4")
//	color.Blue("For going back press 5")
//	fmt.Scanln(&choice)
//
//	switch choice {
//	case 1:
//		var newEmail string
//		for {
//			fmt.Print("Enter Email: ")
//			fmt.Scanln(&newEmail)
//			newEmail = strings.TrimSpace(newEmail)
//			err := userService.UpdateUser(userID, &newEmail, &user.Password, &user.Address, &user.Contact)
//			if err != nil {
//				color.Red("%v", err)
//			} else {
//				fmt.Println("User updated successfully!")
//			}
//			break
//		}
//	case 2:
//		var newPassword string
//		for {
//			fmt.Print("Enter Password: ")
//			fmt.Scanln(&newPassword)
//			newPassword = strings.TrimSpace(newPassword)
//			err := userService.UpdateUser(userID, &user.Email, &newPassword, &user.Address, &user.Contact)
//			if err != nil {
//				color.Red("%v", err)
//			} else {
//				fmt.Println("User updated successfully!")
//			}
//			break
//
//		}
//	case 3:
//		var newPhone string
//		for {
//			fmt.Print("Enter Contact : ")
//			fmt.Scanln(&newPhone)
//			newPhone = strings.TrimSpace(newPhone)
//			err := userService.UpdateUser(userID, &user.Email, &user.Password, &user.Address, &newPhone)
//			if err != nil {
//				color.Red("%v", err)
//			} else {
//				fmt.Println("User updated successfully!")
//			}
//			break
//		}
//	case 4:
//		var newAddress string
//		fmt.Print("Enter Address: ")
//		fmt.Scanln(&newAddress)
//		newAddress = strings.TrimSpace(newAddress)
//		err := userService.UpdateUser(userID, &user.Email, &user.Password, &newAddress, &user.Contact)
//		if err != nil {
//			color.Red("%v", err)
//		} else {
//			fmt.Println("User updated successfully!")
//
//		}
//
//	case 5:
//		return
//	default:
//		color.Red("Invalid choice")
//	}
//
//}
//func viewProfile(user *model.User) {
//	userRepo := repository.NewUserRepository(nil)
//	userService := service.NewUserService(userRepo)
//	currUser, err := userService.ViewProfileByID(user.ID)
//	if err != nil {
//		color.Red("%v", err)
//	}
//	color.Cyan("User Name: %s\n", currUser.Name)
//	color.Cyan("User Email: %s\n", currUser.Email)
//	color.Cyan("User Address: %s\n", currUser.Address)
//	color.Cyan("User Contact: %s\n", currUser.Contact)
//	color.Cyan("User Role: %s\n", currUser.Role)
//	for {
//		color.Blue("For updateProfile press 1")
//		color.Blue("For previous menu press 2")
//		var choice int
//		fmt.Scanln(&choice)
//		if choice == 1 {
//			updateProfile(currUser)
//		} else {
//			break
//		}
//	}
//
//}
//
//// view services provide by provider
//func viewServices(householderService interfaces.HouseholderService) {
//	//color.Blue("Available Service are: ")
//	//services, err := householderService.GetAvailableServices()
//	//if err != nil {
//	//	color.Red("%v", err)
//	//}
//	//for _, service := range services {
//	//	fmt.Println("Name: ", service.Name)
//	//	fmt.Println("Description: ", service.Description)
//	//	fmt.Println("Price: ", service.Price)
//	//	fmt.Println("-------------------------------------")
//	//	//color.Green(service_test.Name)
//	//	//color.Green(service_test.Description)
//	//	//color.Green("%v", service_test.Price)
//	//
//	//}
//	util.DisplayCategory()
//	color.Blue("View Services by Category:")
//	var category string
//	fmt.Print("Enter the category: ")
//	fmt.Scanln(&category)
//
//	services, err := householderService.GetServicesByCategory(category)
//	if err != nil {
//		color.Red("Error fetching services: %v", err)
//		return
//	}
//
//	if len(services) == 0 {
//		color.Cyan("No services found in this category.")
//		return
//	}
//
//	color.Cyan("Available Services in %s:", category)
//	for _, service := range services {
//		color.Cyan("Service ID: %s, Name: %s, Description: %s, Price: %.2f, Provider: %s",
//			service.ID, service.Name, service.Description, service.Price, service.ProviderName)
//		color.Cyan("Provider Contact: %s, Provider Address: %s, Provider Rating: %.2f",
//			service.ProviderContact, service.ProviderAddress, service.ProviderRating)
//	}
//}
//
//// SearchService allows the householder to search for available service_test providers
//func searchService(householderService interfaces.HouseholderService, householder *model.Householder) {
//	util.DisplayCategory()
//
//	var category string
//	fmt.Print("Enter the category of services you're looking for: ")
//	fmt.Scanln(&category)
//
//	services, err := householderService.GetServicesByCategory(category)
//	if err != nil {
//		color.Red("Error fetching services: %v", err)
//		return
//	}
//
//	if len(services) == 0 {
//		color.Cyan("No services found in this category.")
//		return
//	}
//
//	color.Cyan("Available Services:")
//	for _, service := range services {
//		color.Cyan("Service ID: %s, Name: %s, Description: %s, Price: %.2f, Provider: %s",
//			service.ID, service.Name, service.Description, service.Price, service.ProviderName)
//		color.Cyan("Provider Contact: %s, Provider Address: %s, Provider Rating: %.2f",
//			service.ProviderContact, service.ProviderAddress, service.ProviderRating)
//	}
//}
//
//// RequestService allows the householder to request a specific service_test
//func requestService(householderService interfaces.HouseholderService, user *model.Householder) {
//	util.DisplayCategory()
//	var serviceType string
//	fmt.Print("Enter the type of service you want to request: ")
//	fmt.Scanln(&serviceType)
//
//	reader := bufio.NewReader(os.Stdin)
//	fmt.Print("Enter the new scheduled time (YYYY-MM-DD HH:MM): ")
//	newTimeStr, _ := reader.ReadString('\n')
//	newTimeStr = strings.TrimSpace(newTimeStr)
//	newTime, err := time.Parse("2006-01-02 15:04", newTimeStr)
//	if err != nil {
//		color.Red("Error parsing new time: %v", err)
//		return
//	}
//	requestID, err := householderService.RequestService(user.ID, serviceType, &newTime)
//	if err != nil {
//		color.Red("Error requesting service: %v", err)
//		return
//	}
//
//	color.Green("Service requested successfully! Your request ID is %s", requestID)
//}
//
//// ViewBookingHistory allows the householder to view their booking history
//func viewBookingHistory(householderService interfaces.HouseholderService, user *model.User) {
//	history, err := householderService.ViewBookingHistory(user.ID)
//	if err != nil {
//		color.Red("Error viewing booking history: %v", err)
//		return
//	}
//
//	if len(history) == 0 {
//		color.Cyan("No booking history found.")
//		return
//	}
//
//	color.Cyan("Booking History:")
//	for _, booking := range history {
//		color.Cyan("- %s: %s (%s)", booking.ID, booking.ServiceID, booking.Status)
//	}
//}
//
//// LeaveReview allows the householder to leave a review for a service_test provider
//func leaveReview(householderService interfaces.HouseholderService, user *model.User) {
//	var serviceID, providerID string
//
//	var rating float64
//
//	fmt.Print("Enter the Service ID you want to review: ")
//	fmt.Scanln(&serviceID)
//
//	fmt.Print("Enter the Provider ID you want to review: ")
//	fmt.Scanln(&providerID)
//
//	fmt.Print("Enter your review: ")
//	reader := bufio.NewReader(os.Stdin)
//	reviewText, err := reader.ReadString('\n')
//	if err != nil {
//		color.Red("Error submitting review: %v", err)
//		return
//	}
//	fmt.Print("Enter your rating (1-5): ")
//	fmt.Scanln(&rating)
//
//	// Call the AddReview method with the providerID now included
//	err = householderService.AddReview(providerID, user.ID, serviceID, reviewText, rating)
//	if err != nil {
//		color.Red("Error submitting review: %v", err)
//		return
//	}
//
//	color.Green("Review submitted successfully!")
//}
//
//func cancelServiceRequest(householderService interfaces.HouseholderService, householderID string) {
//	var requestID string
//	fmt.Print("Enter the Service Request ID you want to cancel: ")
//	fmt.Scanln(&requestID)
//
//	err := householderService.CancelServiceRequest(requestID, householderID)
//	if err != nil {
//		color.Red("Error canceling service request: %v", err)
//		return
//	}
//
//	color.Green("Service request %s has been successfully canceled.", requestID)
//}
//
//func rescheduleServiceRequest(householderService interfaces.HouseholderService, householderID string) {
//	var requestID string
//	fmt.Print("Enter the Service Request ID you want to reschedule: ")
//	fmt.Scanln(&requestID)
//
//	reader := bufio.NewReader(os.Stdin)
//	fmt.Print("Enter the new scheduled time (YYYY-MM-DD HH:MM): ")
//	newTimeStr, _ := reader.ReadString('\n')
//	newTimeStr = strings.TrimSpace(newTimeStr)
//	newTime, err := time.Parse("2006-01-02 15:04", newTimeStr)
//	if err != nil {
//		color.Red("Error parsing new time: %v", err)
//		return
//	}
//
//	err = householderService.RescheduleServiceRequest(requestID, newTime, householderID)
//	if err != nil {
//		color.Red("Error rescheduling service_test request: %v", err)
//		return
//	}
//
//	color.Green("Service request %s has been successfully rescheduled.", requestID)
//}
//
////	func viewServiceRequestStatus(householderService *service.HouseholderService) {
////		var requestID string
////		fmt.Print("Enter the Service Request ID: ")
////		fmt.Scanln(&requestID)
////
////		status, err := householderService.ViewServiceRequestStatus(requestID)
////		if err != nil {
////			color.Red("Error viewing service_test request status: %v", err)
////			return
////		}
////
////		color.Cyan("Status of service_test request %s: %s", requestID, status)
////	}
//func viewStatus(householderService interfaces.HouseholderService, householder *model.Householder) {
//	// Fetch all service requests for the householder
//	requests, err := householderService.ViewStatus(householder.ID)
//	if err != nil {
//		color.Red("Error viewing status: %v", err)
//	}
//
//	if len(requests) == 0 {
//		color.Cyan("You have no service requests.")
//		return
//	}
//
//	for _, request := range requests {
//		color.Cyan("Request ID: %s, Service ID: %s, Status: %s", request.ID, request.ServiceID, request.Status)
//		if request.Status == "Accepted" && request.ProviderDetails != nil && !request.ApproveStatus {
//			for _, provider := range request.ProviderDetails {
//				color.Green("ServiceProvider Details:")
//				color.Green("ID: %v", provider.ServiceProviderID)
//				color.Green("Name: %s", provider.Name)
//				color.Green("Contact: %s", provider.Contact)
//				color.Green("Address: %s", provider.Address)
//				color.Green("Price: %s", provider.Price)
//				color.Green("Rating: %v", provider.Rating)
//				//fmt.Println("Reviews for ", provider.Name)
//				if len(provider.Reviews) == 0 {
//					color.Cyan("Provider has no Rivews.")
//					continue
//				}
//				var reviewPresent bool
//				for _, review := range provider.Reviews {
//					if review.ServiceID == request.ServiceID {
//						reviewPresent = true
//						color.Green("Rate[1-5]: %v", review.Rating)
//						color.Green("Comment: %v", review.Comments)
//						color.Green("Date: %v", review.ReviewDate)
//					}
//				}
//				if !reviewPresent {
//					color.Cyan("Provider has no Reviews for service %v.", request.ServiceID)
//
//				}
//				fmt.Println("-------------------------------------")
//			}
//		}
//
//		fmt.Println()
//
//	}
//	for {
//		var choice string
//		fmt.Println("1.Cancel any Accepted request")
//		fmt.Println("2.Approve Requests")
//		fmt.Println("3.Previous Menu")
//		fmt.Scanln(&choice)
//		switch choice {
//		case "1":
//			cancelAcceptedServiceRequest(householderService, householder.ID)
//		case "2":
//			ApproveRequest(householderService, householder.ID)
//		case "3":
//			return
//		default:
//			color.Red("Invalid choice")
//		}
//	}
//
//}
//func cancelAcceptedServiceRequest(householderService interfaces.HouseholderService, householderID string) {
//	var requestID string
//	fmt.Print("Enter the Service Request ID you want to cancel: ")
//	fmt.Scanln(&requestID)
//
//	err := householderService.CancelAcceptedRequest(requestID, householderID)
//	if err != nil {
//		color.Red("Error canceling service request: %v", err)
//		return
//	}
//
//	color.Green("Service request %s has been successfully canceled.", requestID)
//}
//
//func ApproveRequest(householderService interfaces.HouseholderService, householderID string) {
//	reader := bufio.NewReader(os.Stdin)
//
//	// Prompt householder for the service request ID
//	fmt.Print("Enter Service Request ID to approve: ")
//	requestID, _ := reader.ReadString('\n')
//	requestID = strings.TrimSpace(requestID)
//
//	// Prompt householder for the provider ID
//	fmt.Print("Enter Service Provider ID you wish to approve: ")
//	providerID, _ := reader.ReadString('\n')
//	providerID = strings.TrimSpace(providerID)
//
//	// Call the approval function
//	if err := householderService.ApproveServiceRequest(requestID, providerID, householderID); err != nil {
//		color.Red("Error approving service request: %v", err)
//		return
//	}
//
//	color.Green("Service request approved successfully!")
//}
//func viewApprovedRequests(householderService interfaces.HouseholderService, householderID string) {
//	// Call the service method to get approved requests
//	approvedRequests, err := householderService.ViewApprovedRequests(householderID)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	// Display the approved requests
//	if len(approvedRequests) == 0 {
//		fmt.Println("No approved service requests found.")
//		return
//	}
//
//	fmt.Println("Approved Service Requests:")
//	for _, req := range approvedRequests {
//		if req.ApproveStatus {
//			fmt.Printf("\nRequest ID: %s\n", req.ID)
//			fmt.Printf("Service ID: %s\n", req.ServiceID)
//			fmt.Printf("Requested Time: %s\n", req.RequestedTime.Format("2006-01-02 15:04:05"))
//			fmt.Printf("Scheduled Time: %s\n", req.ScheduledTime.Format("2006-01-02 15:04:05"))
//			fmt.Printf("Status: %s\n", req.Status)
//			fmt.Println("Provider Details:")
//			for _, provider := range req.ProviderDetails {
//
//				if provider.Approve {
//					fmt.Printf("\tProvider ID: %s\n", provider.ServiceProviderID)
//					fmt.Printf("\tName: %s\n", provider.Name)
//					fmt.Printf("\tContact: %s\n", provider.Contact)
//					fmt.Printf("\tAddress: %s\n", provider.Address)
//					fmt.Printf("\tPrice: %s\n", provider.Price)
//					fmt.Printf("\tRating: %.2f\n", provider.Rating)
//					//fmt.Println("\tReviews:")
//					for _, review := range provider.Reviews {
//						fmt.Printf("\t\tReview ID: %s\n", review.ID)
//						fmt.Printf("\t\tRating: %.2f\n", review.Rating)
//						fmt.Printf("\t\tComments: %s\n", review.Comments)
//						fmt.Printf("\t\tReview Date: %s\n", review.ReviewDate.Format("2006-01-02"))
//					}
//				}
//			}
//		}
//
//	}
//}
//
//// HouseholderDashboard is the main dashboard for householder actions
//func householderDashboard(user *model.User, client *sql.DB) {
//
//	// Initialize repositories and services
//	householderRepo := repository.NewHouseholderRepository(client)
//	serviceRequestRepo := repository.NewServiceRequestRepository(client)
//	serviceProviderRepo := repository.NewServiceProviderRepository(client)
//	serviceRepo := repository.NewServiceRepository(client)
//	householderService := service.NewHouseholderService(householderRepo, serviceProviderRepo, serviceRepo, serviceRequestRepo)
//
//	// Convert the User to a Householder
//	householder := &model.Householder{
//		User: *user,
//	}
//	for {
//		color.Blue("-------------Householder Dashboard--------------")
//		color.Blue("1. View Profile")
//		color.Blue("2  View Services")
//		color.Blue("3. Search Service")
//		color.Blue("4. Request Service")
//		color.Blue("5. View Booking History")
//		color.Blue("6. Leave Review")
//		color.Blue("7. Cancel Service Request")
//		color.Blue("8. Reschedule Service Request")
//		color.Blue("9. View Service Request Status")
//		color.Blue("10. View Approved Request")
//		color.Blue("11.Exit")
//
//		var choice int
//		fmt.Scanln(&choice)
//
//		switch choice {
//		case 1:
//			viewProfile(user)
//		case 2:
//			viewServices(householderService)
//		case 3:
//			searchService(householderService, householder)
//		case 4:
//			requestService(householderService, householder)
//		case 5:
//			viewBookingHistory(householderService, user)
//		case 6:
//			leaveReview(householderService, user)
//		case 7:
//			cancelServiceRequest(householderService, householder.ID)
//		case 8:
//			rescheduleServiceRequest(householderService, householder.ID)
//		case 9:
//			viewStatus(householderService, householder)
//		case 10:
//			viewApprovedRequests(householderService, user.ID)
//		case 11:
//			return
//		default:
//			color.Red("Invalid choice")
//		}
//	}
//}
