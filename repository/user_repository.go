//go:build !test
// +build !test

package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"serviceNest/config"
	"serviceNest/database"
	"serviceNest/interfaces"
	"serviceNest/model"
	"time"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) interfaces.UserRepository {
	if collection == nil {
		collection = database.GetCollection(config.DB, config.USERCOLLECTION)
	}

	return &UserRepository{collection: collection}
}

// SaveUser adds a new user to the repository_test
func (repo *UserRepository) SaveUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.collection.InsertOne(ctx, user)
	return err
}

// GetUserByEmail retrieves a user by their email
func (repo *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user *model.User
	err := repo.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user's details
func (repo *UserRepository) UpdateUser(updatedUser *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": updatedUser.ID}
	update := bson.M{}

	if updatedUser.Email != "" {
		// Ensure the new email doesn't already exist in the system
		existingUser, err := repo.GetUserByEmail(updatedUser.Email)
		if err == nil && existingUser.ID != updatedUser.ID {
			return fmt.Errorf("email already in use")
		}
		update["email"] = updatedUser.Email
	}

	if updatedUser.Password != "" {
		//hashPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		//if err != nil {
		//	return err
		//}
		update["password"] = updatedUser.Password
	}

	if updatedUser.Address != "" {
		update["address"] = updatedUser.Address
	}

	if updatedUser.Contact != "" {
		update["contact"] = updatedUser.Contact
	}

	updateOp := bson.M{"$set": update}
	_, err := repo.collection.UpdateOne(ctx, filter, updateOp, options.Update().SetUpsert(true))
	return err
}

// GetUserByID retrieves a user by their unique ID
func (repo *UserRepository) GetUserByID(userID string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User
	err := repo.collection.FindOne(ctx, bson.M{"id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

//func (repo *UserRepository) EnsureConnection() error {
//	if err := repo.client.Ping(context.TODO(), nil); err != nil {
//		return fmt.Errorf("lost connection to MongoDB: %v", err)
//	}
//	return nil
//}
