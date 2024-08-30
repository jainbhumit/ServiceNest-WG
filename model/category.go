//go:build !test
// +build !test

package model

type Category struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
