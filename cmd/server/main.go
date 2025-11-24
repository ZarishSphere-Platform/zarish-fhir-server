package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/api"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/database"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/search"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to Database
	database.Connect()

	// Init Search
	search.Init()

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
