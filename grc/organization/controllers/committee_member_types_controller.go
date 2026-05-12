package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommitteeMemberTypeController struct {
	service services.CommitteeMemberTypeServiceInterface
	db      *gorm.DB
}

func NewCommitteeMemberTypeController(service services.CommitteeMemberTypeServiceInterface) *CommitteeMemberTypeController {
	return &CommitteeMemberTypeController{
		service: service,
	}
}

func (cmtc *CommitteeMemberTypeController) GetAllCommitteeMemberTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	committeeMemberTypes, err := cmtc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(cmtc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, committeeMemberTypes)
}

func (cmtc *CommitteeMemberTypeController) GetCommitteeMemberTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	committeeMemberType, err := cmtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, committeeMemberType)
}
