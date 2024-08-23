package model

type ServiceProvider struct {
	User
	ServicesOffered []Service `json:"services_offered"`
	Rating          float64   `json:"rating"`
	Reviews         []*Review `json:"reviews"`
	Availability    bool      `json:"availability"`
	IsActive        bool      `json:"is_active"`
}
