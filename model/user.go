//go:build !test
// +build !test

package model

type User struct {
	ID        string  `json:"id" bson:"id"`
	Name      string  `json:"name" bson:"name"`
	Email     string  `json:"email" bson:"email"`
	Password  string  `json:"password" bson:"password"`
	Role      string  `json:"role" bson:"role"` // Householder or ServiceProvider
	Address   string  `json:"address" bson:"address"`
	Contact   string  `json:"contact" bson:"contact"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}
