package controllers

import (
	"jethur/grc/kpi-register/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIReviewStatusController struct {
	service services.KPIReviewStatusServiceInterface
	db      *gorm.DB
}

func NewKPIReviewStatusController(service services.KPIReviewStatusServiceInterface) *KPIReviewStatusController {
	return &KPIReviewStatusController{
		service: service,
	}
}

func (krsc *KPIReviewStatusController) GetAllKPIReviewStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	kpiReviewStatuses, err := krsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(krsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, kpiReviewStatuses)
}

func (krsc *KPIReviewStatusController) GetKPIReviewStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiReviewStatus, err := krsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, kpiReviewStatus)
}
