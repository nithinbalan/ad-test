package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommitteeTermController struct {
	service services.CommitteeTermServiceInterface
	db      *gorm.DB
}

func NewCommitteeTermController(service services.CommitteeTermServiceInterface) *CommitteeTermController {
	return &CommitteeTermController{
		service: service,
	}
}

func (ctc *CommitteeTermController) GetAllCommitteeTerms(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	typeFilter := c.Query("type")

	committeeTerms, err := ctc.service.GetAll(sort, sortBy, keyword, status, typeFilter)

	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, committeeTerms)
}

func (ctc *CommitteeTermController) GetCommitteeTermDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	committeeTerm, err := ctc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, committeeTerm)
}
