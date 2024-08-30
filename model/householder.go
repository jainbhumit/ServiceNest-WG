//go:build !test
// +build !test

package model

type Householder struct {
	User
	ListOfServices []string `json:"list_of_services"`
}
