package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditRiskRatingController struct {
	service services.AuditRiskRatingServiceInterface
	db      *gorm.DB
}

func NewAuditRiskRatingController(service services.AuditRiskRatingServiceInterface) *AuditRiskRatingController {
	return &AuditRiskRatingController{
		service: service,
	}
}

func (arrc *AuditRiskRatingController) GetAllAuditRiskRatings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	auditRiskRating, err := arrc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditRiskRating)
}

func (arrc *AuditRiskRatingController) GetAuditRiskRatingDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditStatus, err := arrc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditStatus)
}
