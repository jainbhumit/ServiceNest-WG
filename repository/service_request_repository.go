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

type ServiceRequestRepository struct {
	collection *mongo.Collection
}

// NewServiceRequestRepository initializes a new ServiceRequestRepository with MongoDB
func NewServiceRequestRepository(collection *mongo.Collection) interfaces.ServiceRequestRepository {
	if collection == nil {
		collection = database.GetCollection("serviceNestDB", "serviceRequests")
	}
	return &ServiceRequestRepository{collection: collection}
}

// SaveServiceRequest saves a service_test request to the MongoDB collection
func (repo *ServiceRequestRepository) SaveServiceRequest(request model.ServiceRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.collection.InsertOne(ctx, request)
	return err
}

// GetServiceRequestByID retrieves a service_test request by its ID from MongoDB
func (repo *ServiceRequestRepository) GetServiceRequestByID(requestID string) (*model.ServiceRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var request model.ServiceRequest
	err := repo.collection.FindOne(ctx, bson.M{"ID": requestID}).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("service_test request not found")
		}
		return nil, err
	}

	return &request, nil
}

// GetServiceRequestsByHouseholderID retrieves all service_test requests made by a specific householder from MongoDB
func (repo *ServiceRequestRepository) GetServiceRequestsByHouseholderID(householderID string) ([]model.ServiceRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := repo.collection.Find(ctx, bson.M{"HouseholderID": householderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []model.ServiceRequest
	if err = cursor.All(ctx, &requests); err != nil {
		return nil, err
	}

	return requests, nil
}

// UpdateServiceRequest updates an existing service_test request in MongoDB
func (repo *ServiceRequestRepository) UpdateServiceRequest(updatedRequest model.ServiceRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.collection.UpdateOne(
		ctx,
		bson.M{"ID": updatedRequest.ID},
		bson.M{"$set": updatedRequest},
	)
	return err
}

// GetAllServiceRequests retrieves all service_test requests from MongoDB
func (repo *ServiceRequestRepository) GetAllServiceRequests() ([]model.ServiceRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []model.ServiceRequest
	if err = cursor.All(ctx, &requests); err != nil {
		return nil, err
	}

	return requests, nil
}

// SaveAllServiceRequests saves all service_test requests to the MongoDB collection (batch save)
func (repo *ServiceRequestRepository) SaveAllServiceRequests(serviceRequests []model.ServiceRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var operations []mongo.WriteModel
	for _, request := range serviceRequests {
		operation := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"ID": request.ID}).
			SetUpdate(bson.M{"$set": request}).
			SetUpsert(true)
		operations = append(operations, operation)
	}

	_, err := repo.collection.BulkWrite(ctx, operations)
	return err
}
