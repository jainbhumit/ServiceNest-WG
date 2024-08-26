//func TestSaveService(t *testing.T) {
//	repo := repository.NewServiceRepository()
//
//	service := model.Service{
//		ID:          "test-id",
//		Name:        "Test Service",
//		Description: "Test Description",
//		Price:       100.0,
//		ProviderID:  "test-provider-id",
//		Category:    "Test Category",
//	}
//
//	err := repo.SaveService(service)
//	assert.NoError(t, err)
//}

// package repository_test
//
// import (
//
//	"github.com/stretchr/testify/assert"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
//	"serviceNest/model"
//	"serviceNest/repository"
//	"testing"
//
// )
//
//	func TestServiceRepository(t *testing.T) {
//		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
//		//defer mt.Close()
//
//		// Create the repository instance with the mock collection
//		repo := repository.NewServiceRepository(mt.Coll)
//
//		t.Run("GetAllServices_Success", func(t *testing.T) {
//			services := []model.Service{
//				{ID: "1", Name: "Service 1", Category: "Cleaning"},
//				{ID: "2", Name: "Service 2", Category: "Plumbing"},
//			}
//			mt.AddMockResponses(mtest.CreateCursorResponse(1, "serviceNestDB.services", mtest.FirstBatch, bson.D{
//				{Key: "id", Value: "1"},
//				{Key: "name", Value: "Service 1"},
//				{Key: "category", Value: "Cleaning"},
//			}, bson.D{
//				{Key: "id", Value: "2"},
//				{Key: "name", Value: "Service 2"},
//				{Key: "category", Value: "Plumbing"},
//			}))
//
//			result, err := repo.GetAllServices()
//			assert.NoError(t, err)
//			assert.Equal(t, 2, len(result))
//			assert.Equal(t, services[0].Name, result[0].Name)
//		})
//
//		t.Run("GetServiceByID_Success", func(t *testing.T) {
//			service := model.Service{ID: "1", Name: "Service 1", Category: "Cleaning"}
//			mt.AddMockResponses(mtest.CreateCursorResponse(1, "serviceNestDB.services", mtest.FirstBatch, bson.D{
//				{Key: "id", Value: "1"},
//				{Key: "name", Value: "Service 1"},
//				{Key: "category", Value: "Cleaning"},
//			}))
//
//			result, err := repo.GetServiceByID("1")
//			assert.NoError(t, err)
//			assert.NotNil(t, result)
//			assert.Equal(t, service.Name, result.Name)
//		})
//
//		t.Run("SaveService_Success", func(t *testing.T) {
//			service := model.Service{ID: "3", Name: "Service 3", Category: "Electrical"}
//			mt.AddMockResponses(mtest.CreateSuccessResponse())
//
//			err := repo.SaveService(service)
//			assert.NoError(t, err)
//		})
//
//		t.Run("RemoveService_Success", func(t *testing.T) {
//			mt.AddMockResponses(mtest.CreateSuccessResponse())
//
//			err := repo.RemoveService("1")
//			assert.NoError(t, err)
//		})
//	}
package repository_test

import (
	"log"
	"serviceNest/database"
	"serviceNest/model"
	"serviceNest/repository"
	"testing"
)

func TestSaveService(t *testing.T) {
	client := database.MockConnect()
	defer database.MockDisconnect()

	if client == nil {
		log.Fatal("Error connecting to database")
	}
	collection := database.GetCollection("TestServiceNestDB", "TestServices")
	repo := repository.NewServiceRepository(collection)
	var mockData model.Service
	mockData = model.Service{
		ID:          "service1",
		Name:        "Cleaning",
		Description: "House cleaning service",
		Price:       50.0,
		ProviderID:  "provider1",
		Category:    "Household",
	}
	err := repo.SaveService(mockData)
	if err != nil {
		t.Error(err)
	}
}

//func TestGetServiceByID(t *testing.T) {
//	collection := database.GetCollection("TestServiceNestDB", "TestServices")
//	repo := repository.NewServiceRepository(collection)
//
//	service, err := repo.GetServiceByID("service1")
//	assert.NoError(t, err, "Failed to get service by ID")
//	assert.Equal(t, "service1", service.ID, "Service ID does not match")
//}
//
//func TestGetAllServices(t *testing.T) {
//	collection := database.GetCollection("TestServiceNestDB", "TestServices")
//	repo := repository.NewServiceRepository(collection)
//
//	services, err := repo.GetAllServices()
//	assert.NoError(t, err, "Failed to get all services")
//	assert.Greater(t, len(services), 0, "No services found")
//}
//
//func TestRemoveService(t *testing.T) {
//	collection := database.GetCollection("TestServiceNestDB", "TestServices")
//	repo := repository.NewServiceRepository(collection)
//
//	err := repo.RemoveService("service1")
//	assert.NoError(t, err, "Failed to remove service")
//}
