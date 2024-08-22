package model

type User struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Role      string  `json:"role"` // Householder or ServiceProvider
	Address   string  `json:"address"`
	Contact   string  `json:"contact"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
