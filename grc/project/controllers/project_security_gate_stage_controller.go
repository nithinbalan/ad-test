package controllers

import (
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateStageController struct {
	service services.ProjectSecurityGateStageServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateStageController(service services.ProjectSecurityGateStageServiceInterface) *ProjectSecurityGateStageController {
	return &ProjectSecurityGateStageController{
		service: service,
	}
}

func (psgsc *ProjectSecurityGateStageController) GetAllProjectSecurityGateStages(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	stages, err := psgsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(psgsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, stages)
}

func (psgsc *ProjectSecurityGateStageController) GetProjectSecurityGateStageDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	stage, err := psgsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, stage)
}
