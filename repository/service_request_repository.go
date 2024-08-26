//package repository
//
//import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"os"
//	"serviceNest/model"
//)
//
//type ServiceRequestRepository struct {
//	filePath string
//}
//
//// NewServiceRequestRepository initializes a new ServiceRequestRepository
//func NewServiceRequestRepository(filePath string) *ServiceRequestRepository {
//	return &ServiceRequestRepository{filePath: filePath}
//}
//
//// SaveServiceRequest saves a service request to the file
//func (repo *ServiceRequestRepository) SaveServiceRequest(request model.ServiceRequest) error {
//	requests, err := repo.loadServiceRequests()
//	if err != nil {
//		return err
//	}
//
//	requests = append(requests, request)
//
//	return repo.saveServiceRequests(requests)
//}
//
//// GetServiceRequestByID retrieves a service request by its ID
//func (repo *ServiceRequestRepository) GetServiceRequestByID(requestID string) (*model.ServiceRequest, error) {
//	requests, err := repo.loadServiceRequests()
//	if err != nil {
//		return nil, err
//	}
//
//	for _, request := range requests {
//		if request.ID == requestID {
//			return &request, nil
//		}
//	}
//
//	return nil, fmt.Errorf("service request with ID %s not found", requestID)
//}
//
//// GetServiceRequestsByHouseholderID retrieves all service requests made by a specific householder
//func (repo *ServiceRequestRepository) GetServiceRequestsByHouseholderID(householderID string) ([]model.ServiceRequest, error) {
//	requests, err := repo.loadServiceRequests()
//	if err != nil {
//		return nil, err
//	}
//
//	var householderRequests []model.ServiceRequest
//	for _, request := range requests {
//		if *request.HouseholderID == householderID {
//			householderRequests = append(householderRequests, request)
//		}
//	}
//
//	return householderRequests, nil
//}
//
//// UpdateServiceRequest updates an existing service request
//func (repo *ServiceRequestRepository) UpdateServiceRequest(updatedRequest model.ServiceRequest) error {
//	requests, err := repo.loadServiceRequests()
//	if err != nil {
//		return err
//	}
//
//	for i, request := range requests {
//		if request.ID == updatedRequest.ID {
//			requests[i] = updatedRequest
//			break
//		}
//	}
//
//	return repo.saveServiceRequests(requests)
//}
//
//// GetAllServiceRequests retrieves all service requests from the file
//func (r *ServiceRequestRepository) GetAllServiceRequests() ([]model.ServiceRequest, error) {
//	var serviceRequests []model.ServiceRequest
//
//	file, err := os.Open(r.filePath)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	bytes, err := ioutil.ReadAll(file)
//	if err != nil {
//		return nil, err
//	}
//
//	err = json.Unmarshal(bytes, &serviceRequests)
//	if err != nil {
//		return nil, err
//	}
//
//	return serviceRequests, nil
//}
//
//// SaveAllServiceRequests saves all service requests to the file
//func (r *ServiceRequestRepository) SaveAllServiceRequests(serviceRequests []model.ServiceRequest) error {
//	// Convert the service requests slice to JSON
//	data, err := json.MarshalIndent(serviceRequests, "", "  ")
//	if err != nil {
//		return err
//	}
//
//	// Write the JSON data to the file
//	err = ioutil.WriteFile(r.filePath, data, os.ModePerm)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// Private helper methods for loading and saving service requests
//func (repo *ServiceRequestRepository) loadServiceRequests() ([]model.ServiceRequest, error) {
//	var serviceRequests []model.ServiceRequest
//
//	// Check if the file exists
//	if _, err := os.Stat(repo.filePath); os.IsNotExist(err) {
//		// File does not exist, return an empty slice
//		return serviceRequests, nil
//	}
//
//	// File exists, proceed to read it
//	file, err := ioutil.ReadFile(repo.filePath)
//	if err != nil {
//		return nil, fmt.Errorf("could not read file: %v", err)
//	}
//
//	err = json.Unmarshal(file, &serviceRequests)
//	if err != nil {
//		return nil, fmt.Errorf("could not unmarshal file: %v", err)
//	}
//
//	return serviceRequests, nil
//}
//
//func (repo *ServiceRequestRepository) saveServiceRequests(requests []model.ServiceRequest) error {
//	data, err := json.MarshalIndent(requests, "", "  ")
//	if err != nil {
//		return err
//	}
//
//	return ioutil.WriteFile(repo.filePath, data, 0644)
//}

package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"serviceNest/database"
	"serviceNest/model"
	"time"
)

type ServiceRequestRepository struct {
	collection *mongo.Collection
}

// NewServiceRequestRepository initializes a new ServiceRequestRepository with MongoDB
func NewServiceRequestRepository(collection *mongo.Collection) *ServiceRequestRepository {
	if collection == nil {
		collection = database.GetCollection("serviceNestDB", "serviceRequests")
	}
	return &ServiceRequestRepository{collection: collection}
}

// SaveServiceRequest saves a service request to the MongoDB collection
func (repo *ServiceRequestRepository) SaveServiceRequest(request model.ServiceRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.collection.InsertOne(ctx, request)
	return err
}

// GetServiceRequestByID retrieves a service request by its ID from MongoDB
func (repo *ServiceRequestRepository) GetServiceRequestByID(requestID string) (*model.ServiceRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var request model.ServiceRequest
	err := repo.collection.FindOne(ctx, bson.M{"ID": requestID}).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("service request not found")
		}
		return nil, err
	}

	return &request, nil
}

// GetServiceRequestsByHouseholderID retrieves all service requests made by a specific householder from MongoDB
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

// UpdateServiceRequest updates an existing service request in MongoDB
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

// GetAllServiceRequests retrieves all service requests from MongoDB
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

// SaveAllServiceRequests saves all service requests to the MongoDB collection (batch save)
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
