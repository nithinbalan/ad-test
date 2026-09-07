package controllers

import (
	"fmt"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateReadinessStatusController struct {
	service services.ProjectSecurityGateReadinessStatusServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateReadinessStatusController(service services.ProjectSecurityGateReadinessStatusServiceInterface) *ProjectSecurityGateReadinessStatusController {
	return &ProjectSecurityGateReadinessStatusController{
		service: service,
	}
}

func (psgrsc *ProjectSecurityGateReadinessStatusController) GetAllProjectSecurityGateReadinessStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	validSortFields := map[string]bool{
		"id":         true,
		"title_en":   true,
		"title_ar":   true,
		"type":       true,
		"color_code": true,
		"status":     true,
	}

	if sortBy != "" && !validSortFields[sortBy] {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid 'sort_by' parameter"))
		return
	}

	readinessStatuses, err := psgrsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(psgrsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, readinessStatuses)
}

func (psgrsc *ProjectSecurityGateReadinessStatusController) GetProjectSecurityGateReadinessStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	readinessStatus, err := psgrsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, readinessStatus)
}
