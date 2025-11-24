package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/api"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/database"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to Database
	database.Connect()

	r := gin.Default()

	// Register Routes
	api.RegisterRoutes(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"service": "zarish-fhir-server",
		})
	})

	log.Printf("Starting Zarish Sphere FHIR Server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
