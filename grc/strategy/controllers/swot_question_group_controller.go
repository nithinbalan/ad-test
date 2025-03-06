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

type SwotQuestionGroupController struct {
	service services.SwotQuestionGroupServiceInterface
	db      *gorm.DB
}

func NewSwotQuestionGroupController(service services.SwotQuestionGroupServiceInterface) *SwotQuestionGroupController {
	return &SwotQuestionGroupController{
		service: service,
	}
}

func (sqgc *SwotQuestionGroupController) GetAllSwotQuestionGroups(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	swotQuestionGroups, err := sqgc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, swotQuestionGroups)
}

func (sqgc *SwotQuestionGroupController) GetSwotQuestionGroupDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	swotQuestionGroup, err := sqgc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, swotQuestionGroup)
}

func (sqgc *SwotQuestionGroupController) CreateSwotQuestionGroup(c *gin.Context) {
	var swotQuestionGroup *models.SwotQuestionGroup

	if err := c.ShouldBindJSON(&swotQuestionGroup); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}

	createdSwotQuestionGroup, err := sqgc.service.Create(currentUserID, swotQuestionGroup)
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

	c.JSON(http.StatusCreated, createdSwotQuestionGroup)
}

func (sqgc *SwotQuestionGroupController) UpdateSwotQuestionGroup(c *gin.Context) {
	swotQuestionGroupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedSwotQuestionGroup *models.SwotQuestionGroup
	if err := c.ShouldBindJSON(&updatedSwotQuestionGroup); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqgc.db, c, err.Error())
		return
	}

	updatedSwotQuestionGroup, err = sqgc.service.Update(uint(swotQuestionGroupID), currentUserID, updatedSwotQuestionGroup)
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

	c.JSON(http.StatusOK, updatedSwotQuestionGroup)
}

func (sqgc *SwotQuestionGroupController) DeleteSwotQuestionGroups(c *gin.Context) {
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

func (sqgc *SwotQuestionGroupController) ActivateSwotQuestionGroups(c *gin.Context) {
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

func (sqgc *SwotQuestionGroupController) DeactivateSwotQuestionGroups(c *gin.Context) {
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
