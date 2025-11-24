package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/database"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/models"
	"gorm.io/datatypes"
)

func CreateResource(c *gin.Context) {
	resourceType := c.Param("resourceType")
	var body map[string]interface{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Basic validation: Ensure resourceType in body matches URL
	if rt, ok := body["resourceType"].(string); !ok || rt != resourceType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resourceType mismatch or missing"})
		return
	}

	// Generate ID if not present
	id, ok := body["id"].(string)
	if !ok || id == "" {
		id = uuid.New().String()
		body["id"] = id
	}

	// Convert body to JSON for storage
	jsonContent, err := datatypes.JSONValue(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process JSON"})
		return
	}

	resource := models.FHIRResource{
		ID:           id,
		ResourceType: resourceType,
		Content:      jsonContent,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if result := database.DB.Create(&resource); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, body)
}

func GetResource(c *gin.Context) {
	resourceType := c.Param("resourceType")
	id := c.Param("id")

	var resource models.FHIRResource
	if result := database.DB.Where("id = ? AND resource_type = ?", id, resourceType).First(&resource); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	c.JSON(http.StatusOK, resource.Content)
}
