package controllers

import (
	"jethur/grc/violation-report/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportStatusController struct {
	service services.ViolationReportStatusServiceInterface
	db      *gorm.DB
}

func NewViolationReportStatusController(service services.ViolationReportStatusServiceInterface) *ViolationReportStatusController {
	return &ViolationReportStatusController{
		service: service,
	}
}

func (vrsc *ViolationReportStatusController) GetAllViolationReportStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	violationReportStatuses, err := vrsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(vrsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, violationReportStatuses)
}

func (vrsc *ViolationReportStatusController) GetViolationReportStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	violationReportStatus, err := vrsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, violationReportStatus)
}
