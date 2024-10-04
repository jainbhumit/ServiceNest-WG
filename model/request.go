package model

import "time"

type ServiceRequest struct {
	ID                 string                   `json:"id" bson:"ID"`
	HouseholderID      *string                  `json:"householder_id" bson:"HouseholderID"`
	HouseholderName    string                   `json:"householder_name" bson:"HouseholderName"`
	HouseholderAddress *string                  `json:"householder_address" bson:"HouseholderAddress"`
	ServiceName        string                   `json:"service_name" bson:"ServiceName"`
	ServiceID          string                   `json:"service_id" bson:"serviceID"`
	RequestedTime      time.Time                `json:"requested_time" bson:"requestedTime"`
	ScheduledTime      time.Time                `json:"scheduled_time" bson:"scheduledTime"`
	Status             string                   `json:"status" bson:"status"` // Pending, Accepted, Completed, Cancelled
	ApproveStatus      bool                     `json:"approve_status" bson:"approveStatus"`
	ProviderDetails    []ServiceProviderDetails `json:"provider_details,omitempty" bson:"providerDetails,omitempty"`
}
type ServiceProviderDetails struct {
	ServiceProviderID string   `json:"service_provider_id" bson:"serviceProviderID"`
	Name              string   `json:"name" bson:"name"`
	Contact           string   `json:"contact" bson:"contact"`
	Address           string   `json:"address" bson:"address"`
	Price             string   `json:"price" bson:"price"`
	Rating            float64  `json:"rating" bson:"rating"`
	Reviews           []Review `json:"reviews,omitempty" bson:"reviews"`
	Approve           bool     `json:"approve" bson:"approve"`
}
