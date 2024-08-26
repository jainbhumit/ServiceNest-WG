package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"os/signal"
	"serviceNest/database"
	"syscall"
)

func main() {
	// Initialize MongoDB Connection
	client := database.Connect()
	defer database.Disconnect()

	if client == nil {
		log.Fatal("Error connecting to database")
	}

	// Handle interrupt signals for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nDisconnecting from MongoDB...")
		database.Disconnect()
		os.Exit(1)
	}()

	for {
		fmt.Println("-----------------------Welcome-----------------------")
		color.Blue("For SignUp press 1\n")
		color.Blue("For Login press 2\n")
		color.Blue("For Exit press 3\n")
		var choice int
		color.Cyan("Enter your choice: ")
		fmt.Scanln(&choice)
		switch choice {
		case 1:

			SignUpUser()
		case 2:

			LoginUser()
		case 3:
			return
		default:
			color.Red("Invalid choice")

		}
	}

}
