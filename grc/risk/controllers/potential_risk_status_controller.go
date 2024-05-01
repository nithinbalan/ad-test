package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PotentialRiskStatusController struct {
	service services.PotentialRiskStatusServiceInterface
	db      *gorm.DB
}

func NewPotentialRiskStatusController(service services.PotentialRiskStatusServiceInterface) *PotentialRiskStatusController {
	return &PotentialRiskStatusController{
		service: service,
	}
}

func (prsc *PotentialRiskStatusController) GetAllPotentialRiskStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	potentialRiskStatuses, err := prsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(prsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, potentialRiskStatuses)
}

func (prsc *PotentialRiskStatusController) GetPotentialRiskStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	potentialRiskStatus, err := prsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, potentialRiskStatus)
}
