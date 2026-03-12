package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommitteeTypeController struct {
	service services.CommitteeTypeServiceInterface
	db      *gorm.DB
}

func NewCommitteeTypeController(service services.CommitteeTypeServiceInterface) *CommitteeTypeController {
	return &CommitteeTypeController{
		service: service,
	}
}

func (ctc *CommitteeTypeController) GetAllCommitteeTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	typeFilter := c.Query("type")

	committeeTypes, err := ctc.service.GetAll(sort, sortBy, keyword, status, typeFilter)

	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, committeeTypes)
}

func (ctc *CommitteeTypeController) GetCommitteeTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	committeeType, err := ctc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, committeeType)
}
