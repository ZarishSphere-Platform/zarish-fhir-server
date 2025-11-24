package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zarishsphere-platform/zarish-fhir-server/internal/auth"
)

func RegisterRoutes(r *gin.Engine) {
	fhir := r.Group("/fhir/R4")
	fhir.Use(auth.AuthMiddleware())
	{
		fhir.POST("/:resourceType", CreateResource)
		fhir.GET("/:resourceType", SearchResource)
		fhir.GET("/:resourceType/:id", GetResource)
	}
}
