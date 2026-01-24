package main

import (
	"log"
	"os"

	"github.com/richardlinio/waterballsa-backend-golang/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Println("Failed to initialize application:", err)
		os.Exit(1)
	}

	err = application.Run()
	if err != nil {
		log.Println("Application terminated with error:", err)
		os.Exit(1)
	}
}
