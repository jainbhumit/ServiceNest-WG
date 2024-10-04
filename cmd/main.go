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
	runApp(client)
	//if err := runApp(); err != nil {
	//	log.Fatal(err)
	//}

}
