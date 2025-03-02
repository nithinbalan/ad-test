package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApplicabilityStatusController struct {
	service services.ApplicabilityStatusServiceInterface
	db      *gorm.DB
}

func NewApplicabilityStatusController(service services.ApplicabilityStatusServiceInterface) *ApplicabilityStatusController {
	return &ApplicabilityStatusController{
		service: service,
	}
}

func (asc *ApplicabilityStatusController) GetAllApplicabilityStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	applicabilityStatuses, err := asc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, applicabilityStatuses)
}

func (asc *ApplicabilityStatusController) GetApplicabilityStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, "Invalid id")
		return
	}

	applicabilityStatus, err := asc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, applicabilityStatus)
}
