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

type ProjectSecurityGateChecklistResponseStatusController struct {
	service services.ProjectSecurityGateChecklistResponseStatusServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateChecklistResponseStatusController(service services.ProjectSecurityGateChecklistResponseStatusServiceInterface) *ProjectSecurityGateChecklistResponseStatusController {
	return &ProjectSecurityGateChecklistResponseStatusController{
		service: service,
	}
}

func (psgcrsc *ProjectSecurityGateChecklistResponseStatusController) GetAllProjectSecurityGateChecklistResponseStatuses(c *gin.Context) {
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

	responseStatuses, err := psgcrsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(psgcrsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, responseStatuses)
}

func (psgcrsc *ProjectSecurityGateChecklistResponseStatusController) GetProjectSecurityGateChecklistResponseStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	responseStatus, err := psgcrsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, responseStatus)
}
