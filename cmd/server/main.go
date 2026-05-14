package main

import (
	"dss/internal/app"
	"log"
)

func main() {
	container, err := app.NewContainer()
	if err != nil {
		log.Fatalf("failed to create container: %v", err)
	}
	defer container.Close()

	router := app.SetupRouter(container)

	log.Println("Starting server on :8088")
	if err := router.Run(":8088"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
