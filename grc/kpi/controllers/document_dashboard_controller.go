package controllers

import (
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
	response, err := controller.service.GetDocumentSummary()
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *DocumentDashboardController) GetDocumentStatusByType(c *gin.Context) {
	response, err := controller.service.GetDocumentStatusByType()
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *DocumentDashboardController) GetDocumentOnTimeReviewTrend(c *gin.Context) {
	response, err := controller.service.GetDocumentOnTimeReviewTrend()
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
