package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"serviceNest/config"
	"serviceNest/repository"
	"serviceNest/routers"
	"serviceNest/service"
)

//var userService *service.UserService
//var householderService *service.HouseholderService
//var providerService *service.ServiceProviderService
//var adminService *service.AdminService

func start(client *sql.DB) {
	// initialize all repository
	userRepo := repository.NewUserRepository(client)
	householderRepo := repository.NewHouseholderRepository(client)
	requestRepo := repository.NewServiceRequestRepository(client)
	providerRepo := repository.NewServiceProviderRepository(client)
	serviceRepo := repository.NewServiceRepository(client)

	// initialize all services
	userService := service.NewUserService(userRepo)
	householderService := service.NewHouseholderService(householderRepo, providerRepo, serviceRepo, requestRepo)
	providerService := service.NewServiceProviderService(providerRepo, requestRepo, serviceRepo)
	adminService := service.NewAdminService(serviceRepo, requestRepo, userRepo, providerRepo)

	router := routers.SetupRouter(userService, householderService, providerService, adminService)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Health Good")

	})
	log.Println("Sever Starting on Port 8080...")
	err := http.ListenAndServe(config.PORT, router)
	if err != nil {
		log.Fatal(err)
	}
}
