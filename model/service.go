//go:build !test
// +build !test

package model

type Service struct {
	ID              string  `json:"id" bson:"id"`
	Name            string  `json:"name" bson:"name"`
	Description     string  `json:"description" bson:"description"`
	Price           float64 `json:"price" bson:"price"`
	ProviderID      string  `json:"provider_id" bson:"provider_id"`
	Category        string  `json:"category" bson:"category"`
	ProviderName    string  `json:"provider_name" bson:"provider_name"`
	ProviderContact string  `json:"provider_contact" bson:"provider_contact"`
	ProviderAddress string  `json:"provider_address" bson:"provider_address"`
	ProviderRating  float64 `json:"provider_rating" bson:"provider_rating"`
}
