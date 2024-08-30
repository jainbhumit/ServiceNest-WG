//go:build !test
// +build !test

package model

type ServiceArea struct {
	ID        string  `json:"id" bson:"id"`
	Name      string  `json:"name" bson:"name"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
	Radius    float64 `json:"radius" bson:"radius"` // Service area radius in kilometers
}
