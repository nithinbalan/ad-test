package controllers

import (
	"jethur/grc/models"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SwotPestelQuestionGroupController struct {
	service services.SwotPestelQuestionGroupServiceInterface
	db      *gorm.DB
}

func NewSwotPestelQuestionGroupController(service services.SwotPestelQuestionGroupServiceInterface) *SwotPestelQuestionGroupController {
	return &SwotPestelQuestionGroupController{
		service: service,
	}
}

func (sqgc *SwotPestelQuestionGroupController) GetAllSwotPestelQuestionGroups(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	isSwot := c.Query("is_swot")

	swotPestelQuestionGroups, err := sqgc.service.GetAll(sort, sortBy, keyword, status, isSwot)

	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, swotPestelQuestionGroups)
}

func (sqgc *SwotPestelQuestionGroupController) GetSwotPestelQuestionGroupDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	swotPestelQuestionGroup, err := sqgc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, swotPestelQuestionGroup)
}

func (sqgc *SwotPestelQuestionGroupController) CreateSwotPestelQuestionGroup(c *gin.Context) {
	var swotPestelQuestionGroup *models.SwotPestelQuestionGroup

	if err := c.ShouldBindJSON(&swotPestelQuestionGroup); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}

	createdSwotPestelQuestionGroup, err := sqgc.service.Create(currentUserID, swotPestelQuestionGroup)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sqgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdSwotPestelQuestionGroup)
}

func (sqgc *SwotPestelQuestionGroupController) UpdateSwotPestelQuestionGroup(c *gin.Context) {
	swotPestelQuestionGroupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedSwotPestelQuestionGroup *models.SwotPestelQuestionGroup
	if err := c.ShouldBindJSON(&updatedSwotPestelQuestionGroup); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}

	updatedSwotPestelQuestionGroup, err = sqgc.service.Update(uint(swotPestelQuestionGroupID), currentUserID, updatedSwotPestelQuestionGroup)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sqgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedSwotPestelQuestionGroup)
}

func (sqgc *SwotPestelQuestionGroupController) DeleteSwotPestelQuestionGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}

	failedIDs, err := sqgc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (sqgc *SwotPestelQuestionGroupController) ActivateSwotPestelQuestionGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}

	err = sqgc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (sqgc *SwotPestelQuestionGroupController) DeactivateSwotPestelQuestionGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}

	err = sqgc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
