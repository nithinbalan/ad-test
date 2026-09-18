package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateDashboardController struct {
	service services.ProjectSecurityGateDashboardServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateDashboardController(service services.ProjectSecurityGateDashboardServiceInterface, db *gorm.DB) *ProjectSecurityGateDashboardController {
	return &ProjectSecurityGateDashboardController{service: service, db: db}
}

func (psgdc *ProjectSecurityGateDashboardController) buildParams(c *gin.Context) *dtos.ProjectSecurityGateDashboardQueryParamsDTO {
	year, _ := strconv.Atoi(c.Query("year"))
	return &dtos.ProjectSecurityGateDashboardQueryParamsDTO{
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Year:     year,
	}
}

func (psgdc *ProjectSecurityGateDashboardController) GetSummary(c *gin.Context) {
	result, err := psgdc.service.GetSummary(psgdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(psgdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (psgdc *ProjectSecurityGateDashboardController) GetExceptionSummary(c *gin.Context) {
	result, err := psgdc.service.GetExceptionSummary(psgdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(psgdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (psgdc *ProjectSecurityGateDashboardController) GetChecklistsByResponseStatus(c *gin.Context) {
	result, err := psgdc.service.GetChecklistsByResponseStatus(psgdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(psgdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (psgdc *ProjectSecurityGateDashboardController) GetByReadinessStatus(c *gin.Context) {
	result, err := psgdc.service.GetByReadinessStatus(psgdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(psgdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (psgdc *ProjectSecurityGateDashboardController) GetByClassification(c *gin.Context) {
	result, err := psgdc.service.GetByClassification(psgdc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(psgdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
