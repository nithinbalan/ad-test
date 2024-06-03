package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceReviewStatusController struct {
	service services.ComplianceReviewStatusServiceInterface
	db      *gorm.DB
}

func NewComplianceReviewStatusController(service services.ComplianceReviewStatusServiceInterface) *ComplianceReviewStatusController {
	return &ComplianceReviewStatusController{
		service: service,
	}
}

func (crsc *ComplianceReviewStatusController) GetAllComplianceReviewStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	complianceReviewStatuses, err := crsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(crsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, complianceReviewStatuses)
}

func (crsc *ComplianceReviewStatusController) GetComplianceReviewStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceReviewStatus, err := crsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, complianceReviewStatus)
}
