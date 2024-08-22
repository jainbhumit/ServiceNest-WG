package model

type Admin struct {
	User
	Permissions []string `json:"permissions"`
}

func NewAdmin(id, name, email, password string) *Admin {
	return &Admin{
		User: User{
			ID:       id,
			Name:     name,
			Email:    email,
			Password: password,
			Role:     "Admin",
		},
		Permissions: []string{"ManageUsers", "ManageServices", "ManageRequests", "ManagePayments"},
	}
}
