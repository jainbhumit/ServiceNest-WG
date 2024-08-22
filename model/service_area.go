package model

type ServiceArea struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"` // Service area radius in kilometers
}

func NewServiceArea(id, name string, lat, lon, radius float64) *ServiceArea {
	return &ServiceArea{
		ID:        id,
		Name:      name,
		Latitude:  lat,
		Longitude: lon,
		Radius:    radius,
	}
}
