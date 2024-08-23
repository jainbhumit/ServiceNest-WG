package model

import "time"

type ServiceRequest struct {
	ID                string                  `json:"id"`
	HouseholderID     *string                 `json:"householder_id"`
	ServiceID         string                  `json:"service_id"`
	ServiceProviderID string                  `json:"service_provider_id"`
	RequestedTime     time.Time               `json:"requested_time"`
	ScheduledTime     time.Time               `json:"scheduled_time"`
	Status            string                  `json:"status"` // Pending, Accepted, Completed, Cancelled
	ProviderDetails   *ServiceProviderDetails `json:"provider_details,omitempty"`
}
type ServiceProviderDetails struct {
	Name    string    `json:"name"`
	Contact string    `json:"contact"`
	Address string    `json:"address"`
	Rating  float64   `json:"rating"`
	Reviews []*Review `json:"reviews"`
}
