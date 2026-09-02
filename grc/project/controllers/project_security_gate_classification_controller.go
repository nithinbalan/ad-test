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

type ProjectSecurityGateClassificationController struct {
	service services.ProjectSecurityGateClassificationServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateClassificationController(service services.ProjectSecurityGateClassificationServiceInterface) *ProjectSecurityGateClassificationController {
	return &ProjectSecurityGateClassificationController{
		service: service,
	}
}

func (psgcc *ProjectSecurityGateClassificationController) GetAllProjectSecurityGateClassifications(c *gin.Context) {
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

	classifications, err := psgcc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(psgcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, classifications)
}

func (psgcc *ProjectSecurityGateClassificationController) GetProjectSecurityGateClassificationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	classification, err := psgcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, classification)
}
