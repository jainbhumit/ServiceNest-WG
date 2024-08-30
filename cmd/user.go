//go:build !test
// +build !test

package main

import (
	"github.com/fatih/color"
	"serviceNest/model"
	"serviceNest/repository"
	"serviceNest/service"
)

func SignUpUser() error {
	userRepo := repository.NewUserRepository(nil)

	_, err := service.SignUp(userRepo)
	if err != nil {
		return err
	}
	return nil
}

func LoginUser() error {
	userRepo := repository.NewUserRepository(nil)

	user, err := service.Login(userRepo)
	if err != nil {
		return err
	}
	dashBoard(user)
	return nil
}

func dashBoard(user *model.User) {
	color.Blue("Welcome to Service Nest")

	if user.Role == "Householder" {
		householderDashboard(user)
	} else if user.Role == "ServiceProvider" {
		serviceProviderDashboard(user)
	} else {
		admin := &model.Admin{
			user,
		}
		adminDashboard(admin)
	}

}
