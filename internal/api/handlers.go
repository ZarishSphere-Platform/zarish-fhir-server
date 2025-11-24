package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/database"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/models"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/search"
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
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process JSON"})
		return
	}
	jsonContent := datatypes.JSON(jsonBytes)

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

	// Index in Elasticsearch
	go search.IndexResource(resourceType, id, body)

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

func SearchResource(c *gin.Context) {
	resourceType := c.Param("resourceType")
	queryParams := make(map[string]string)

	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}

	results, err := search.SearchResources(resourceType, queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceType": "Bundle",
		"type":         "searchset",
		"total":        len(results),
		"entry":        results,
	})
}
