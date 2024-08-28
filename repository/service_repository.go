package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"serviceNest/database"
	"serviceNest/interfaces"
	"serviceNest/model"
	"time"
)

type ServiceRepository struct {
	Collection *mongo.Collection
}

// NewServiceRepository creates a new instance of ServiceRepository
func NewServiceRepository(collection *mongo.Collection) interfaces.ServiceRepository {
	if collection == nil {
		// Default to the real MongoDB collection if none is provided
		collection = database.GetCollection("serviceNestDB", "services")
	}
	return &ServiceRepository{Collection: collection}
}

// GetAllServices fetches all available services from the MongoDB database
func (repo *ServiceRepository) GetAllServices() ([]*model.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := repo.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var services []*model.Service
	if err = cursor.All(ctx, &services); err != nil {
		return nil, err
	}

	return services, nil
}

// GetServiceByID retrieves a service_test by its ID
func (repo *ServiceRepository) GetServiceByID(serviceID string) (*model.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var service model.Service
	err := repo.Collection.FindOne(ctx, bson.M{"id": serviceID}).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("service_test not found")
		}
		return nil, err
	}

	return &service, nil
}

// SaveService adds a new service_test to the MongoDB database
func (repo *ServiceRepository) SaveService(service model.Service) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.Collection.InsertOne(ctx, service)

	return err
}

// SaveAllServices saves the entire list of services to the MongoDB database
func (repo *ServiceRepository) SaveAllServices(services []model.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var operations []mongo.WriteModel
	for _, service := range services {
		operation := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"id": service.ID}).
			SetUpdate(bson.M{"$set": service}).
			SetUpsert(true)
		operations = append(operations, operation)
	}

	_, err := repo.Collection.BulkWrite(ctx, operations)
	return err
}

// RemoveService removes a service_test from the MongoDB database
func (repo *ServiceRepository) RemoveService(serviceID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.Collection.DeleteOne(ctx, bson.M{"id": serviceID})
	return err
}
