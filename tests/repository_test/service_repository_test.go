package repository_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"serviceNest/model"
	"serviceNest/repository"
	"serviceNest/tests/mocks"
	"testing"
)

func TestServiceRepository_GetAllServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	repo := repository.NewServiceRepository(mockCollection)

	// Define expected data
	expectedServices := []*model.Service{
		{ID: "service1", Name: "Service 1"},
		{ID: "service2", Name: "Service 2"},
	}

	// Simulate MongoDB cursor behavior
	mockCursor := mocks.NewMockCursor(ctrl)
	mockCollection.EXPECT().
		Find(gomock.Any(), bson.M{}).
		Return(mockCursor, nil)

	mockCursor.EXPECT().All(gomock.Any(), gomock.Any()).SetArg(1, expectedServices).Return(nil)
	mockCursor.EXPECT().Close(gomock.Any()).Return(nil)

	// Execute
	result, err := repo.GetAllServices()

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedServices, result)
}

func TestServiceRepository_GetServiceByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockCollection(ctrl)
	repo := repository.NewServiceRepository(mockCollection)

	// Define expected data
	serviceID := "service1"
	expectedService := &model.Service{ID: serviceID, Name: "Service 1"}

	// Simulate MongoDB FindOne behavior
	mockSingleResult := mocks.NewMockSingleResult(ctrl)
	mockCollection.EXPECT().
		FindOne(gomock.Any(), bson.M{"id": serviceID}).
		Return(mockSingleResult)

	mockSingleResult.EXPECT().Decode(gomock.Any()).SetArg(0, *expectedService).Return(nil)

	// Execute
	result, err := repo.GetServiceByID(serviceID)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedService, result)
}

func TestServiceRepository_SaveService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockCollection(ctrl)
	repo := repository.NewServiceRepository(mockCollection)

	// Define test data
	newService := model.Service{ID: "service1", Name: "Service 1"}

	// Simulate MongoDB InsertOne behavior
	mockCollection.EXPECT().
		InsertOne(gomock.Any(), newService).
		Return(&mongo.InsertOneResult{}, nil)

	// Execute
	err := repo.SaveService(newService)

	// Verify
	assert.NoError(t, err)
}

func TestServiceRepository_SaveAllServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockCollection(ctrl)
	repo := repository.NewServiceRepository(mockCollection)

	// Define test data
	services := []model.Service{
		{ID: "service1", Name: "Service 1"},
		{ID: "service2", Name: "Service 2"},
	}

	// Simulate MongoDB BulkWrite behavior
	mockCollection.EXPECT().
		BulkWrite(gomock.Any(), gomock.Any()).
		Return(&mongo.BulkWriteResult{}, nil)

	// Execute
	err := repo.SaveAllServices(services)

	// Verify
	assert.NoError(t, err)
}

func TestServiceRepository_RemoveService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockCollection(ctrl)
	repo := repository.NewServiceRepository(mockCollection)

	// Define test data
	serviceID := "service1"

	// Simulate MongoDB DeleteOne behavior
	mockCollection.EXPECT().
		DeleteOne(gomock.Any(), bson.M{"id": serviceID}).
		Return(&mongo.DeleteResult{}, nil)

	// Execute
	err := repo.RemoveService(serviceID)

	// Verify
	assert.NoError(t, err)
}
