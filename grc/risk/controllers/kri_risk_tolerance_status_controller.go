package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRIRiskToleranceStatusController struct {
	service services.KRIRiskToleranceStatusServiceInterface
	db      *gorm.DB
}

func NewKRIRiskToleranceStatusController(service services.KRIRiskToleranceStatusServiceInterface) *KRIRiskToleranceStatusController {
	return &KRIRiskToleranceStatusController{
		service: service,
	}
}

func (krtsc *KRIRiskToleranceStatusController) GetAllKRIRiskToleranceStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	statuses, err := krtsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(krtsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, statuses)
}

func (krtsc *KRIRiskToleranceStatusController) GetKRIRiskToleranceStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	status, err := krtsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, status)
}
