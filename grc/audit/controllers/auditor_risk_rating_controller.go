package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditorRiskRatingController struct {
	service services.AuditorRiskRatingServiceInterface
	db      *gorm.DB
}

func NewAuditorRiskRatingController(service services.AuditorRiskRatingServiceInterface) *AuditorRiskRatingController {
	return &AuditorRiskRatingController{
		service: service,
	}
}

func (arrc *AuditorRiskRatingController) GetAllAuditorRiskRatings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	auditorRiskRating, err := arrc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditorRiskRating)
}

func (arrc *AuditorRiskRatingController) GetAuditorRiskRatingDetails(c *gin.Context) {
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
