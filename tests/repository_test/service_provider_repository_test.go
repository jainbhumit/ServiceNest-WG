//func TestSaveService(t *testing.T) {
//	repo := repository_test.NewServiceRepository()
//
//	service_test := model.Service{
//		ID:          "test-id",
//		Name:        "Test Service",
//		Description: "Test Description",
//		Price:       100.0,
//		ProviderID:  "test-provider-id",
//		Category:    "Test Category",
//	}
//
//	err := repo.SaveService(service_test)
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
//	"serviceNest/repository_test"
//	"testing"
//
// )
//
//	func TestServiceRepository(t *testing.T) {
//		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
//		//defer mt.Close()
//
//		// Create the repository_test instance with the mock collection
//		repo := repository_test.NewServiceRepository(mt.Coll)
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
//			service_test := model.Service{ID: "1", Name: "Service 1", Category: "Cleaning"}
//			mt.AddMockResponses(mtest.CreateCursorResponse(1, "serviceNestDB.services", mtest.FirstBatch, bson.D{
//				{Key: "id", Value: "1"},
//				{Key: "name", Value: "Service 1"},
//				{Key: "category", Value: "Cleaning"},
//			}))
//
//			result, err := repo.GetServiceByID("1")
//			assert.NoError(t, err)
//			assert.NotNil(t, result)
//			assert.Equal(t, service_test.Name, result.Name)
//		})
//
//		t.Run("SaveService_Success", func(t *testing.T) {
//			service_test := model.Service{ID: "3", Name: "Service 3", Category: "Electrical"}
//			mt.AddMockResponses(mtest.CreateSuccessResponse())
//
//			err := repo.SaveService(service_test)
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
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"serviceNest/model"
	"serviceNest/tests/mocks"
	"testing"
)

//import (
//	"log"
//	"serviceNest/database"
//	"serviceNest/model"
//	"serviceNest/repository_test"
//	"testing"
//)
//
//func TestSaveService(t *testing.T) {
//	client := database.MockConnect()
//	defer database.MockDisconnect()
//
//	if client == nil {
//		log.Fatal("Error connecting to database")
//	}
//	collection := database.GetCollection("TestServiceNestDB", "TestServices")
//	repo := repository_test.NewServiceRepository(collection)
//	var mockData model.Service
//	mockData = model.Service{
//		ID:          "service1",
//		Name:        "Cleaning",
//		Description: "House cleaning service_test",
//		Price:       50.0,
//		ProviderID:  "provider1",
//		Category:    "Household",
//	}
//	err := repo.SaveService(mockData)
//	if err != nil {
//		t.Error(err)
//	}
//}

//	func TestGetServiceByID(t *testing.T) {
//		collection := database.GetCollection("TestServiceNestDB", "TestServices")
//		repo := repository_test.NewServiceRepository(collection)
//
//		service_test, err := repo.GetServiceByID("service1")
//		assert.NoError(t, err, "Failed to get service_test by ID")
//		assert.Equal(t, "service1", service_test.ID, "Service ID does not match")
//	}
//
//	func TestGetAllServices(t *testing.T) {
//		collection := database.GetCollection("TestServiceNestDB", "TestServices")
//		repo := repository_test.NewServiceRepository(collection)
//
//		services, err := repo.GetAllServices()
//		assert.NoError(t, err, "Failed to get all services")
//		assert.Greater(t, len(services), 0, "No services found")
//	}
//
//	func TestRemoveService(t *testing.T) {
//		collection := database.GetCollection("TestServiceNestDB", "TestServices")
//		repo := repository_test.NewServiceRepository(collection)
//
//		err := repo.RemoveService("service1")
//		assert.NoError(t, err, "Failed to remove service_test")
//	}
func TestServiceProviderRepository_AddReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	providerID := "provider123"
	householderID := "householder123"
	review := "Great service!"
	rating := 4.5

	// Expectations
	mockRepo.EXPECT().
		AddReview(providerID, householderID, review, rating).
		Return(nil)

	// Execute
	err := mockRepo.AddReview(providerID, householderID, review, rating)

	// Verify
	assert.NoError(t, err)
}

func TestServiceProviderRepository_GetProviderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	providerID := "provider123"
	expectedProvider := &model.ServiceProvider{
		User: model.User{
			ID:   providerID,
			Name: "Test Provider",
		},
	}

	// Expectations
	mockRepo.EXPECT().
		GetProviderByID(providerID).
		Return(expectedProvider, nil)

	// Execute
	result, err := mockRepo.GetProviderByID(providerID)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedProvider, result)
}

func TestServiceProviderRepository_GetProviderByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	providerID := "invalid-provider-id"
	expectedError := errors.New("provider not found")

	// Expectations
	mockRepo.EXPECT().
		GetProviderByID(providerID).
		Return(nil, expectedError)

	// Execute
	result, err := mockRepo.GetProviderByID(providerID)

	// Verify
	assert.Nil(t, result)
	assert.EqualError(t, err, expectedError.Error())
}

func TestServiceProviderRepository_GetProvidersByServiceType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	serviceType := "Plumbing"
	expectedProviders := []model.ServiceProvider{
		{
			User: model.User{
				ID:   "provider123",
				Name: "Provider 1",
			},
		},
		{
			User: model.User{
				ID:   "provider456",
				Name: "Provider 2",
			},
		},
	}

	// Expectations
	mockRepo.EXPECT().
		GetProvidersByServiceType(serviceType).
		Return(expectedProviders, nil)

	// Execute
	result, err := mockRepo.GetProvidersByServiceType(serviceType)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedProviders, result)
}

func TestServiceProviderRepository_SaveServiceProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	provider := model.ServiceProvider{
		User: model.User{
			ID:   "provider123",
			Name: "Test Provider",
		},
	}

	// Expectations
	mockRepo.EXPECT().
		SaveServiceProvider(provider).
		Return(nil)

	// Execute
	err := mockRepo.SaveServiceProvider(provider)

	// Verify
	assert.NoError(t, err)
}

func TestServiceProviderRepository_UpdateServiceProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceProviderRepository(ctrl)

	// Define test data
	provider := &model.ServiceProvider{
		User: model.User{
			ID:   "provider123",
			Name: "Updated Provider",
		},
	}

	// Expectations
	mockRepo.EXPECT().
		UpdateServiceProvider(provider).
		Return(nil)

	// Execute
	err := mockRepo.UpdateServiceProvider(provider)

	// Verify
	assert.NoError(t, err)
}
