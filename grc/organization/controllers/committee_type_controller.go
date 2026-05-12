package controllers

import (
	"jethur/grc/organization/dtos"
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

	queryParams := dtos.CommitteeTypeListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Type:    typeFilter,
	}

	committeeTypes, err := ctc.service.GetAll(&queryParams)

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

func (ctc *CommitteeTypeController) CreateCommitteeType(c *gin.Context) {
	var ctDTO dtos.CommitteeTypeFormDTO

	if err := c.ShouldBindJSON(&ctDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}

	id, err := ctc.service.Create(currentUserID, &ctDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(ctc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (ctc *CommitteeTypeController) UpdateCommitteeType(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var ctDTO dtos.CommitteeTypeFormDTO
	if err := c.ShouldBindJSON(&ctDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}

	err = ctc.service.Update(uint(id), currentUserID, &ctDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ctc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (ctc *CommitteeTypeController) DeleteCommitteeTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}

	failedIds, err := ctc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}

	if len(failedIds) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Some committee types could not be deleted",
			"failed_ids": failedIds,
		})
		return
	}

	c.Status(http.StatusOK)
}
