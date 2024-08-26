// package repository
//
// import (
//
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"os"
//	"serviceNest/model"
//
// )
//
//	type HouseholderRepository struct {
//		filePath string
//	}
//
// // NewHouseholderRepository initializes a new HouseholderRepository
//
//	func NewHouseholderRepository(filePath string) *HouseholderRepository {
//		return &HouseholderRepository{filePath: filePath}
//	}
//
// // SaveHouseholder saves a new householder to the file
//
//	func (repo *HouseholderRepository) SaveHouseholder(householder model.Householder) error {
//		householders, err := repo.loadHouseholders()
//		if err != nil {
//			return err
//		}
//
//		householders = append(householders, householder)
//
//		return repo.saveHouseholders(householders)
//	}
//
// // GetHouseholderByID retrieves a householder by their ID
//
//	func (repo *HouseholderRepository) GetHouseholderByID(id string) (*model.Householder, error) {
//		householders, err := repo.loadHouseholders()
//		if err != nil {
//			return nil, err
//		}
//
//		for _, householder := range householders {
//			if householder.ID == id {
//				return &householder, nil
//			}
//		}
//
//		return nil, fmt.Errorf("householder not found")
//	}
//
// // Private helper methods for loading and saving householders and service requests
//
//	func (repo *HouseholderRepository) loadHouseholders() ([]model.Householder, error) {
//		var householders []model.Householder
//
//		file, err := ioutil.ReadFile(repo.filePath)
//		if err != nil {
//			if os.IsNotExist(err) {
//				return householders, nil
//			}
//			return nil, err
//		}
//
//		err = json.Unmarshal(file, &householders)
//		if err != nil {
//			return nil, err
//		}
//
//		return householders, nil
//	}
//
//	func (repo *HouseholderRepository) saveHouseholders(householders []model.Householder) error {
//		data, err := json.MarshalIndent(householders, "", "  ")
//		if err != nil {
//			return err
//		}
//
//		return ioutil.WriteFile(repo.filePath, data, 0644)
//	}
package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"serviceNest/database"
	"serviceNest/model"
	"time"
)

type HouseholderRepository struct {
	collection *mongo.Collection
}

// NewHouseholderRepository initializes a new HouseholderRepository
func NewHouseholderRepository() *HouseholderRepository {
	collection := database.GetCollection("serviceNestDB", "users")
	return &HouseholderRepository{collection: collection}
}

// SaveHouseholder saves a new householder to the database
func (repo *HouseholderRepository) SaveHouseholder(householder model.Householder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.collection.InsertOne(ctx, householder)
	return err
}

// GetHouseholderByID retrieves a householder by their ID
func (repo *HouseholderRepository) GetHouseholderByID(id string) (*model.Householder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var householder model.Householder
	err := repo.collection.FindOne(ctx, bson.M{"id": id}).Decode(&householder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("householder not found")
		}
		return nil, err
	}
	return &householder, nil
}
