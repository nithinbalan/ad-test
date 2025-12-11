package controllers

import (
	"net/http"

	"jethur/grc/settings/services"

	"github.com/gin-gonic/gin"
)

type HealthCheckController struct {
	service services.HealthCheckService
}

func NewHealthCheckController(service services.HealthCheckService) *HealthCheckController {
	return &HealthCheckController{service: service}
}

func (hc *HealthCheckController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func (hc *HealthCheckController) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

func (hc *HealthCheckController) Readiness(c *gin.Context) {

	if err := hc.service.CheckSystemPreference(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unavailable",
			"detail": "system preference check failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
