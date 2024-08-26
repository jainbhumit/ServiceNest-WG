package model

type Admin struct {
	User *User `json:"user" bson:"user"`
	//ID       string `json:"id"`
	//Name     string `json:"name"`
	//Email    string `json:"email"`
	//Password string `json:"password"`
	//Role     string `json:"role"` // should be "admin"
}
