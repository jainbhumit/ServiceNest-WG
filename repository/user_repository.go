// package repository
//
// import (
//
//	"encoding/json"
//	"fmt"
//	"golang.org/x/crypto/bcrypt"
//	"io/ioutil"
//	"os"
//	"serviceNest/database"
//	"serviceNest/model"
//
// )
//
// //type IUserRepository interface {
// //	GetUserByID(userID string) (*model.User, error)
// //	GetUserByEmail(email string) (*model.User, error)
// //	UpdateUser(updatedUser *model.User) error
// //	SaveUser(user model.User) error
// //	loadUsers() ([]model.User, error)
// //	saveUsers(users []model.User) error
// //}
//
//	type UserRepository struct {
//		filePath string,
//		collection *mongo.Collection,
//	}
//
//	func NewUserRepository(filePath string) *UserRepository {
//		collection := database.GetCollection("serviceNestDB", "users")
//		return &UserRepository{filePath: filePath,
//			collection: collection}
//	}
//
// // SaveUser adds a new user to the repository
//
//	func (repo *UserRepository) SaveUser(user model.User) error {
//		users, err := repo.loadUsers()
//		if err != nil {
//			return err
//		}
//
//		users = append(users, user)
//
//		return repo.saveUsers(users)
//	}
//
// // GetUserByEmail retrieves a user by their email
//
//	func (repo *UserRepository) GetUserByEmail(email string) (*model.User, error) {
//		users, err := repo.loadUsers()
//		if err != nil {
//			return nil, err
//		}
//
//		for _, user := range users {
//			if user.Email == email {
//				return &user, nil
//			}
//		}
//
//		return nil, fmt.Errorf("user not found")
//	}
//
// // UpdateUser updates an existing user's details
//
//	func (repo *UserRepository) UpdateUser(updatedUser *model.User) error {
//		users, err := repo.loadUsers()
//		if err != nil {
//			return err
//		}
//
//		for i, user := range users {
//			if user.ID == updatedUser.ID {
//				// Update the user's details
//				if updatedUser.Email != "" && user.Email != updatedUser.Email {
//					// Ensure the new email doesn't already exist in the system
//					if _, err := repo.GetUserByEmail(updatedUser.Email); err == nil {
//						return fmt.Errorf("email already in use")
//					}
//					user.Email = updatedUser.Email
//				}
//				if updatedUser.Password != "" {
//					hashPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
//					if err != nil {
//						return err
//					}
//					user.Password = string(hashPassword)
//				}
//				if updatedUser.Address != "" {
//					user.Address = updatedUser.Address
//				}
//				if updatedUser.Contact != "" {
//					user.Contact = updatedUser.Contact
//				}
//				// Update the slice with the modified user
//				users[i] = user
//				break
//			}
//		}
//
//		// Save the updated user list back to the file
//		return repo.saveUsers(users)
//	}
//
// // loadUsers loads users from the file
//
//	func (repo *UserRepository) loadUsers() ([]model.User, error) {
//		var users []model.User
//
//		file, err := ioutil.ReadFile(repo.filePath)
//		if err != nil {
//			if os.IsNotExist(err) {
//				return users, nil
//			}
//			return nil, err
//		}
//
//		err = json.Unmarshal(file, &users)
//		if err != nil {
//			return nil, err
//		}
//
//		return users, nil
//	}
//
// // saveUsers saves the list of users to the file
//
//	func (repo *UserRepository) saveUsers(users []model.User) error {
//		data, err := json.MarshalIndent(users, "", "  ")
//		if err != nil {
//			return err
//		}
//
//		return ioutil.WriteFile(repo.filePath, data, 0644)
//	}
//
// // GetUserByID retrieves a user by their unique ID
//
//	func (repo *UserRepository) GetUserByID(userID string) (*model.User, error) {
//		users, err := repo.loadUsers()
//		if err != nil {
//			return nil, err
//		}
//
//		for _, user := range users {
//			if user.ID == userID {
//				return &user, nil
//			}
//		}
//
//		return nil, fmt.Errorf("user not found")
//	}
package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"serviceNest/database"
	"serviceNest/model"
	"time"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() *UserRepository {
	collection := database.GetCollection("serviceNestDB", "users")
	return &UserRepository{collection: collection}
}

// SaveUser adds a new user to the repository
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
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		update["password"] = string(hashPassword)
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
