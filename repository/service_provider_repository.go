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
)

type ServiceProviderRepository struct {
	Collection *mongo.Collection
}

// NewServiceProviderRepository initializes a new ServiceProviderRepository
func NewServiceProviderRepository(collection *mongo.Collection) interfaces.ServiceProviderRepository {
	if collection == nil {
		// Default to the real MongoDB collection if none is provided
		collection = database.GetCollection(config.DB, config.SERVICEPROVIDERCOLLECTION)
	}
	return &ServiceProviderRepository{Collection: collection}

}

//// SaveServiceProvider saves a new service_test provider to the file
//func (repo *ServiceProviderRepository) SaveServiceProvider(provider model.ServiceProvider) error {
//	providers, err := repo.loadServiceProviders()
//	if err != nil {
//		return err
//	}
//
//	providers = append(providers, provider)
//
//	return repo.saveServiceProviders(providers)
//}

func (repo *ServiceProviderRepository) SaveServiceProvider(provider model.ServiceProvider) error {
	_, err := repo.Collection.InsertOne(context.TODO(), provider)
	return err
}

// GetProviderByID retrieves a service_test provider by their ID
//
//	func (repo *ServiceProviderRepository) GetProviderByID(id string) (*model.ServiceProvider, error) {
//		providers, err := repo.loadServiceProviders()
//		if err != nil {
//			return nil, err
//		}
//
//		for _, provider := range providers {
//			if provider.ID == id {
//				return &provider, nil
//			}
//		}
//
//		return nil, fmt.Errorf("service_test provider not found")
//	}
func (repo *ServiceProviderRepository) GetProviderByID(providerID string) (*model.ServiceProvider, error) {
	var provider *model.ServiceProvider
	err := repo.Collection.FindOne(context.TODO(), bson.M{"user.id": providerID}).Decode(&provider)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

//// GetProvidersByServiceType retrieves all service_test providers that offer a specific service_test type
//func (repo *ServiceProviderRepository) GetProvidersByServiceType(serviceType string) ([]model.ServiceProvider, error) {
//	providers, err := repo.loadServiceProviders()
//	if err != nil {
//		return nil, err
//	}
//
//	var matchingProviders []model.ServiceProvider
//	for _, provider := range providers {
//		for _, service_test := range provider.ServicesOffered {
//			if service_test.Name == serviceType {
//				matchingProviders = append(matchingProviders, provider)
//				break
//			}
//		}
//	}
//
//	return matchingProviders, nil
//}

func (repo *ServiceProviderRepository) GetProvidersByServiceType(serviceType string) ([]model.ServiceProvider, error) {
	filter := bson.M{"services_offered.name": serviceType}
	cursor, err := repo.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var providers []model.ServiceProvider
	for cursor.Next(context.TODO()) {
		var provider model.ServiceProvider
		if err := cursor.Decode(&provider); err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return providers, nil
}

// GetProviderByServiceID retrieves all service_test providers that offer a specific service_test type
//
//	func (repo *ServiceProviderRepository) GetProviderByServiceID(serviceID string) (*model.ServiceProvider, error) {
//		providers, err := repo.loadServiceProviders()
//		if err != nil {
//			return nil, err
//		}
//
//		var matchingProviders model.ServiceProvider
//		for _, provider := range providers {
//			for _, service_test := range provider.ServicesOffered {
//				if service_test.ID == serviceID {
//					matchingProviders = provider
//					break
//				}
//			}
//		}
//
//		return &matchingProviders, nil
//	}
func (repo *ServiceProviderRepository) GetProviderByServiceID(serviceID string) (*model.ServiceProvider, error) {
	filter := bson.M{"services_offered.id": serviceID}
	var provider model.ServiceProvider
	err := repo.Collection.FindOne(context.TODO(), filter).Decode(&provider)
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// UpdateServiceProvider updates an existing service_test provider
//func (repo *ServiceProviderRepository) UpdateServiceProvider(updatedProvider *model.ServiceProvider) error {
//	providers, err := repo.loadServiceProviders()
//	if err != nil {
//		return err
//	}
//
//	// Find and update the service_test provider
//	for i, provider := range providers {
//		if provider.User.ID == updatedProvider.User.ID {
//			providers[i] = *updatedProvider
//			break
//		}
//	}
//
//	// Save the updated list of service_test providers back to the file
//	return repo.saveServiceProviders(providers)
//}

func (repo *ServiceProviderRepository) UpdateServiceProvider(provider *model.ServiceProvider) error {
	filter := bson.M{"user.id": provider.ID}
	update := bson.M{"$set": provider}
	_, err := repo.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// AddReview adds a review and rating for a service_test provider
//
//	func (repo *ServiceProviderRepository) AddReview(providerID, householderID, review string, rating float64) error {
//		providers, err := repo.loadServiceProviders()
//		if err != nil {
//			return err
//		}
//
//		for i, provider := range providers {
//			if provider.ID == providerID {
//				newReview := model.Review{
//					HouseholderID: householderID,
//					Comments:      review,
//					Rating:        rating,
//				}
//				providers[i].Reviews = append(providers[i].Reviews, &newReview)
//
//				return repo.saveServiceProviders(providers)
//			}
//		}
//
//		return fmt.Errorf("service_test provider not found")
//	}
func (repo *ServiceProviderRepository) AddReview(providerID, householderID, review string, rating float64) error {
	filter := bson.M{"user.id": providerID}
	update := bson.M{
		"$push": bson.M{
			"reviews": model.Review{
				HouseholderID: householderID,
				Comments:      review,
				Rating:        rating,
			},
		},
	}
	_, err := repo.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("could not add review: %v", err)
	}
	return nil
}

// Private helper methods for loading and saving service_test providers
//func (repo *ServiceProviderRepository) loadServiceProviders() ([]model.ServiceProvider, error) {
//	var providers []model.ServiceProvider
//
//	file, err := ioutil.ReadFile(repo.filePath)
//	if err != nil {
//		if os.IsNotExist(err) {
//			return providers, nil
//		}
//		return nil, err
//	}
//
//	err = json.Unmarshal(file, &providers)
//	if err != nil {
//		return nil, err
//	}
//
//	return providers, nil
//}
//
//func (repo *ServiceProviderRepository) saveServiceProviders(providers []model.ServiceProvider) error {
//	data, err := json.MarshalIndent(providers, "", "  ")
//	if err != nil {
//		return err
//	}
//
//	return ioutil.WriteFile(repo.filePath, data, 0644)
//}
