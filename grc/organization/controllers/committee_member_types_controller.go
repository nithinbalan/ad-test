package controllers

import (
	"errors"
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
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")

	validSortFields := map[string]string{
		"id":       "committee_member_types.id",
		"title_en": "committee_member_types.title_en",
		"title_ar": "committee_member_types.title_ar",
		"status":   "committee_member_types.status",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	committeeMemberTypes, err := cmtc.service.GetAll(sort, sortField, keyword, status)
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
