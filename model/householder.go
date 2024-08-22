package model

type Householder struct {
	User
	ListOfServices []string `json:"list_of_services"`
}
