package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	fhir := r.Group("/fhir/R4")
	{
		fhir.POST("/:resourceType", CreateResource)
		fhir.GET("/:resourceType/:id", GetResource)
	}
}
