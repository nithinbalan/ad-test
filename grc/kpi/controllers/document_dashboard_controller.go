package controllers

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentDashboardController struct {
	service services.DocumentDashboardServiceInterface
	db      *gorm.DB
}

func NewDocumentDashboardController(service services.DocumentDashboardServiceInterface, db *gorm.DB) *DocumentDashboardController {
	return &DocumentDashboardController{service: service, db: db}
}

func (controller *DocumentDashboardController) GetDocumentSummary(c *gin.Context) {
	allowedParams := []string{
		"document_category_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())
	queryParams := &kpiDTOs.DocumentDashboardQueryParamsDTO{Filters: filters}

	response, err := controller.service.GetDocumentSummary(queryParams)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *DocumentDashboardController) GetDocumentStatusByType(c *gin.Context) {
	allowedParams := []string{
		"document_category_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())
	queryParams := &kpiDTOs.DocumentDashboardQueryParamsDTO{Filters: filters}

	response, err := controller.service.GetDocumentStatusByType(queryParams)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *DocumentDashboardController) GetDocumentOnTimeReviewTrend(c *gin.Context) {
	allowedParams := []string{
		"document_category_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())
	queryParams := &kpiDTOs.DocumentDashboardQueryParamsDTO{Filters: filters}

	response, err := controller.service.GetDocumentOnTimeReviewTrend(queryParams)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
