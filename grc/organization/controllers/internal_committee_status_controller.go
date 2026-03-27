package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InternalCommitteeStatusController struct {
	service services.InternalCommitteeStatusServiceInterface
	db      *gorm.DB
}

func NewInternalCommitteeStatusController(service services.InternalCommitteeStatusServiceInterface) *InternalCommitteeStatusController {
	return &InternalCommitteeStatusController{
		service: service,
	}
}

func (icsc *InternalCommitteeStatusController) GetAllInternalCommitteeStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	typeFilter := c.Query("type")

	internalCommitteeStatuses, err := icsc.service.GetAll(sort, sortBy, keyword, status, typeFilter)

	if err != nil {
		utility.SendInternalErrResp(icsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, internalCommitteeStatuses)
}

func (icsc *InternalCommitteeStatusController) GetInternalCommitteeStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	internalCommitteeStatus, err := icsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, internalCommitteeStatus)
}
