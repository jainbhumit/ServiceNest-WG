package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"serviceNest/config"
	"serviceNest/database"
	"serviceNest/interfaces"
	"serviceNest/model"
	"time"
)

type HouseholderRepository struct {
	Collection interfaces.MongoCollection
}

// NewHouseholderRepository initializes a new HouseholderRepository
func NewHouseholderRepository() interfaces.HouseholderRepository {
	collection := &MongoCollectionImpl{collection: database.GetCollection(config.DB, config.USERCOLLECTION)}
	return &HouseholderRepository{Collection: collection}
}

// SaveHouseholder saves a new householder to the database
func (repo *HouseholderRepository) SaveHouseholder(householder model.Householder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.Collection.InsertOne(ctx, householder)
	return err
}

// GetHouseholderByID retrieves a householder by their ID
func (repo *HouseholderRepository) GetHouseholderByID(id string) (*model.Householder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var householder model.Householder
	err := repo.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&householder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("householder not found")
		}
		return nil, err
	}
	return &householder, nil
}
