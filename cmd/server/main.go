package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/router"
)

func main() {
	r := gin.Default()

	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	router.SetupRoutes(r, db)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
