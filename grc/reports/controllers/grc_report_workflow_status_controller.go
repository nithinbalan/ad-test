package controllers

import (
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GRCReportWorkflowStatusController struct {
	service services.GRCReportWorkflowStatusServiceInterface
	db      *gorm.DB
}

func NewGRCReportWorkflowStatusController(service services.GRCReportWorkflowStatusServiceInterface) *GRCReportWorkflowStatusController {
	return &GRCReportWorkflowStatusController{
		service: service,
	}
}

func (dsc *GRCReportWorkflowStatusController) GetAllGRCReportWorkflowStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	documentStatuses, err := dsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentStatuses)
}

func (dsc *GRCReportWorkflowStatusController) GetGRCReportWorkflowStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentStatus, err := dsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, documentStatus)
}
