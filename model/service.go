package model

type Service struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ProviderID  string  `json:"provider_id"`
}
