//go:build !test
// +build !test

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"serviceNest/config"
	"serviceNest/logger"
	"syscall"
)

func main() {
	client, err := config.GetMySQLDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		client.Close()
	}()

	if client == nil {
		fmt.Errorf("error connecting to database")
	} else {
		log.Println("Connected to database MySql")
	}

	// Handle interrupt signals for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nDisconnecting from MySql...")
		client.Close()
		os.Exit(1)
	}()
	logger.Info("Start the application..", nil)
	start(client)
	//if err := runApp(); err != nil {
	//	log.Fatal(err)
	//}

}

//func runApp() error {
//
//	client, err := config.GetMySQLDB()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer func() {
//		client.Close()
//	}()
//
//	if client == nil {
//		return fmt.Errorf("error connecting to database")
//	} else {
//		log.Println("Connected to database MySql")
//	}
//
//	// Handle interrupt signals for graceful shutdown
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//	go func() {
//		<-c
//		fmt.Println("\nDisconnecting from MySql...")
//		client.Close()
//		os.Exit(1)
//	}()
//
//	for {
//		fmt.Println("-----------------------Welcome-----------------------")
//		color.Blue("For SignUp press 1\n")
//		color.Blue("For Login press 2\n")
//		color.Blue("For Exit press 3\n")
//		var choice int
//		color.Cyan("Enter your choice: ")
//		fmt.Scanln(&choice)
//		switch choice {
//		case 1:
//			if err := SignUpUser(client); err != nil {
//				color.Red("Error during signup: %s", err)
//			}
//		case 2:
//			if err := LoginUser(client); err != nil {
//				color.Red("Error during login: %s", err)
//			}
//		case 3:
//			return nil
//		default:
//			color.Red("Invalid choice")
//		}
//	}
//}
