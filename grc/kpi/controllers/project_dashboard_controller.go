package controllers

import (
	projectDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectDashboardController struct {
	service services.ProjectDashboardServiceInterface
	db      *gorm.DB
}

func NewProjectDashboardController(service services.ProjectDashboardServiceInterface, db *gorm.DB) *ProjectDashboardController {
	return &ProjectDashboardController{service: service, db: db}
}

func (pdc *ProjectDashboardController) buildParams(c *gin.Context) *projectDTOs.ProjectDashboardQueryParamsDTO {
	year, _ := strconv.Atoi(c.Query("year"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	return &projectDTOs.ProjectDashboardQueryParamsDTO{
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Year:     year,
		Limit:    limit,
	}
}

func (pdc *ProjectDashboardController) GetProjectSummary(c *gin.Context) {
	result, err := pdc.service.GetProjectSummary(pdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(pdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (pdc *ProjectDashboardController) GetProjectsByPhase(c *gin.Context) {
	result, err := pdc.service.GetProjectsByPhase(pdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(pdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (pdc *ProjectDashboardController) GetProjectsByClassification(c *gin.Context) {
	result, err := pdc.service.GetProjectsByClassification(pdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(pdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (pdc *ProjectDashboardController) GetProjectsByStatus(c *gin.Context) {
	result, err := pdc.service.GetProjectsByStatus(pdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(pdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (pdc *ProjectDashboardController) GetProjectCompletion(c *gin.Context) {
	result, err := pdc.service.GetProjectCompletion(pdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(pdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (pdc *ProjectDashboardController) GetProjectsDueSoonAndOverdue(c *gin.Context) {
	result, err := pdc.service.GetProjectsDueSoonAndOverdue(pdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(pdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
