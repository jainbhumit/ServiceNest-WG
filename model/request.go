package model

import "time"

type ServiceRequest struct {
	ID                string                  `json:"id" bson:"ID"`
	HouseholderID     *string                 `json:"householder_id" bson:"HouseholderID"`
	ServiceID         string                  `json:"service_id" bson:"serviceID"`
	ServiceProviderID string                  `json:"service_provider_id" bson:"service_providerID"`
	RequestedTime     time.Time               `json:"requested_time" bson:"requestedTime"`
	ScheduledTime     time.Time               `json:"scheduled_time" bson:"scheduledTime"`
	Status            string                  `json:"status" bson:"status"` // Pending, Accepted, Completed, Cancelled
	ProviderDetails   *ServiceProviderDetails `json:"provider_details,omitempty" bson:"providerDetails,omitempty"`
}
type ServiceProviderDetails struct {
	Name    string    `json:"name" bson:"name"`
	Contact string    `json:"contact" bson:"contact"`
	Address string    `json:"address" bson:"address"`
	Rating  float64   `json:"rating" bson:"rating"`
	Reviews []*Review `json:"reviews" bson:"reviews"`
}
