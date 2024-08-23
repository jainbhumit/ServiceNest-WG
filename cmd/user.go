package main

import (
	"fmt"
	"github.com/fatih/color"
	"serviceNest/model"
	"serviceNest/repository"
	"serviceNest/service"
)

func SignUpUser() {
	userRepo := repository.NewUserRepository("users.json")

	_, err := service.SignUp(userRepo)
	if err != nil {
		fmt.Println("Error during signup:", err)
	}
}

func LoginUser() {
	userRepo := repository.NewUserRepository("users.json")

	user, err := service.Login(userRepo)
	if err != nil {
		fmt.Println("Error during login:", err)
		return
	}
	dashBoard(user)
}

func dashBoard(user *model.User) {
	color.Blue("Welcome to Service Nest")

	if user.Role == "Householder" {
		householderDashboard(user)
	} else if user.Role == "ServiceProvider" {
		serviceProviderDashboard(user)
	} else {
		admin := model.Admin{
			user,
		}
		adminDashboard(&admin)
	}

}
