// //go:build !test
// // +build !test
package main

//
//import (
//	"database/sql"
//	"github.com/fatih/color"
//	"serviceNest/model"
//	"serviceNest/repository"
//)
//
//func SignUpUser(client *sql.DB) error {
//	userRepo := repository.NewUserRepository(client)
//
//	_, err := SignUp(userRepo)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func LoginUser(client *sql.DB) error {
//	userRepo := repository.NewUserRepository(client)
//
//	user, err := Login(userRepo)
//	if err != nil {
//		return err
//	}
//	dashBoard(user, client)
//	return nil
//}
//
//func dashBoard(user *model.User, client *sql.DB) {
//	color.Blue("Welcome to Service Nest")
//
//	if user.Role == "Householder" {
//		householderDashboard(user, client)
//	} else if user.Role == "ServiceProvider" {
//		serviceProviderDashboard(user, client)
//	} else {
//		admin := &model.Admin{
//			user,
//		}
//		adminDashboard(admin, client)
//	}
//
//}
