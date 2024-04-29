package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceStatusController struct {
	service services.ComplianceStatusServiceInterface
	db      *gorm.DB
}

func NewComplianceStatusController(service services.ComplianceStatusServiceInterface) *ComplianceStatusController {
	return &ComplianceStatusController{
		service: service,
	}
}

func (csc *ComplianceStatusController) GetAllComplianceStatuses(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	complianceStatuses, err := csc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(csc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, complianceStatuses)
}

func (csc *ComplianceStatusController) GetComplianceStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceStatus, err := csc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, complianceStatus)
}
