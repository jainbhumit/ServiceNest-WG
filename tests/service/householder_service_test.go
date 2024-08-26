package service

import (
	"log"
	"serviceNest/database"
	"serviceNest/repository"
	"serviceNest/service"
	"serviceNest/tests"
	"testing"
)

func connectDatabase() {
	client := database.MockConnect()

	if client == nil {
		log.Fatal("Error connecting to database")
	}
}
func TestGetServicesByCategory(t *testing.T) {
	connectDatabase()
	collection := database.MockGetCollection("TestServiceNestDB", "TestServices")
	serviceRepo := repository.NewServiceRepository(collection)
	householderService := service.NewHouseholderService(nil, nil, serviceRepo, nil)
	err1 := serviceRepo.SaveService(tests.MockServices[0])
	if err1 != nil {
		t.Error(err1)
	}
	services, err := householderService.GetServicesByCategory("Household")
	if err != nil {
		t.Error(err)
	}
	if len(services) == 0 {
		t.Error("Householder service not created")
	}
	if services[0].Category != "Household" {
		t.Error("Household service created but not fetch right category")
	}
	err = serviceRepo.RemoveService(services[0].ID)
}

func TestSearchService(t *testing.T) {
	//client := database.MockConnect()
	//defer database.MockDisconnect()
	//
	//if client == nil {
	//	log.Fatal("Error connecting to database")
	//}
	connectDatabase()

	collection := database.MockGetCollection("TestServiceNestDB", "TestServiceProvider")
	serviceProviderRepo := repository.NewServiceProviderRepository(collection)
	householderService := service.NewHouseholderService(nil, serviceProviderRepo, nil, nil)
	err1 := serviceProviderRepo.SaveServiceProvider(tests.MockServiceProvider[0])
	if err1 != nil {
		t.Error(err1)
	}
	services, err := householderService.SearchService(nil, "Household")
	if err != nil {
		t.Error(err)
	}
	if len(services) == 0 {
		t.Error("services is empty service not created")
	}
	if services[0].ServicesOffered[0].Name != "Household" {
		t.Error("Provider is present but the service fetch is not from require one")
	}
	//err = serviceProviderRepo.RemoveService(services[0].ID)
}

func TestRequestService(t *testing.T) {
	//client := database.MockConnect()
	//defer database.MockDisconnect()
	//
	//if client == nil {
	//	log.Fatal("Error connecting to database")
	//}
	connectDatabase()

	collection := database.MockGetCollection("TestServiceNestDB", "TestServiceRequests")
	requestRepo := repository.NewServiceRequestRepository(collection)
	householderService := service.NewHouseholderService(nil, nil, nil, requestRepo)
	request, err := householderService.RequestService("householder1", "household")
	if err != nil {
		t.Error(err)
	}
	if request == "" {
		t.Error("Request not generate")
	}

}

func TestViewBookingHistory(t *testing.T) {
	//client := database.MockConnect()
	//defer database.MockDisconnect()
	//
	//if client == nil {
	//	log.Fatal("Error connecting to database")
	//}
	connectDatabase()

	collection := database.MockGetCollection("TestServiceNestDB", "TestServiceRequests")
	requestRepo := repository.NewServiceRequestRepository(collection)
	householderService := service.NewHouseholderService(nil, nil, nil, requestRepo)
	householder := tests.MockHouseholder
	_, err1 := householderService.RequestService(householder[0].ID, "Household")
	if err1 != nil {
		t.Error(err1)
	}
	history, err := householderService.ViewBookingHistory("householder1")
	if err != nil {
		t.Error(err)
	}
	if len(history) == 0 {
		t.Error("history is empty")
	}

}
func TestGetAvailableServices(t *testing.T) {
	//client := database.MockConnect()
	//defer database.MockDisconnect()
	//
	//if client == nil {
	//	log.Fatal("Error connecting to database")
	//}
	connectDatabase()

	collection := database.MockGetCollection("TestServiceNestDB", "TestServices")
	serviceRepo := repository.NewServiceRepository(collection)
	householderService := service.NewHouseholderService(nil, nil, serviceRepo, nil)
	err1 := serviceRepo.SaveService(tests.MockServices[0])
	if err1 != nil {
		t.Error(err1)
	}
	services, err := householderService.GetAvailableServices()
	if err != nil {
		t.Error(err)
	}
	if len(services) == 0 {
		t.Error("services is empty service not created")
	}

}

func TestViewServiceRequestStatus(t *testing.T) {
	connectDatabase()
	collection := database.MockGetCollection("TestServiceNestDB", "TestServiceRequests")
	requestRepo := repository.NewServiceRequestRepository(collection)
	householderService := service.NewHouseholderService(nil, nil, nil, requestRepo)
	householder := tests.MockHouseholder
	request, err1 := householderService.RequestService(householder[0].ID, "Household")
	if err1 != nil {
		t.Error(err1)
	}
	status, err := householderService.ViewServiceRequestStatus(request)
	if err != nil {
		t.Error(err)
	}
	if status == "" {
		t.Error("status is empty")
	}

}
func TestCancelServiceRequest(t *testing.T) {
	//client := database.MockConnect()
	//defer database.MockDisconnect()
	//
	//if client == nil {
	//	log.Fatal("Error connecting to database")
	//}
	connectDatabase()
	defer database.MockDisconnect()

	collection := database.MockGetCollection("TestServiceNestDB", "TestServiceRequests")
	requestRepo := repository.NewServiceRequestRepository(collection)
	householderService := service.NewHouseholderService(nil, nil, nil, requestRepo)
	request, err := householderService.RequestService("householder1", "household")
	if err != nil {
		t.Error(err)
	}
	err = householderService.CancelServiceRequest(request)
	if err != nil {
		t.Error(err)
	}
}

func TestViewStatus(t *testing.T) {
	connectDatabase()
	collection := database.MockGetCollection("TestServiceNestDB", "TestServiceRequests")
	requestRepo := repository.NewServiceRequestRepository(collection)
	householderService := service.NewHouseholderService(nil, nil, nil, requestRepo)
	householder := tests.MockHouseholder
	_, err1 := householderService.RequestService(householder[0].ID, "Household")
	if err1 != nil {
		t.Error(err1)
	}
	requests, err := householderService.ViewStatus(nil, &householder[0])
	if err != nil {
		t.Error(err)
	}
	if len(requests) == 0 {
		t.Error("requests is empty")
	}

}

func TestCancelAcceptedRequest(t *testing.T) {
	connectDatabase()
	defer database.MockDisconnect()
	collection := database.MockGetCollection("TestServiceNestDB", "TestServiceRequests")
	requestRepo := repository.NewServiceRequestRepository(collection)
	householderService := service.NewHouseholderService(nil, nil, nil, requestRepo)
	householder := tests.MockHouseholder
	requestId, err1 := householderService.RequestService(householder[0].ID, "Household")

	if err1 != nil {
		t.Error(err1)
	}
	err := householderService.CancelAcceptedRequest(requestId, householder[0].ID)
	//if err != nil {
	//	if errors.Is(err, errors.New("only accepted service requests can be canceled")) != true {
	//		t.Error(err)
	//	}
	//	t.Error(err)
	//}
	if err == nil {
		t.Error("requestId is empty")
	}
}
