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

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
