package controllers

import (
	"jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionGroupController struct {
	service services.QuestionGroupServiceInterface
	db      *gorm.DB
}

func NewQuestionGroupController(service services.QuestionGroupServiceInterface) *QuestionGroupController {
	return &QuestionGroupController{
		service: service,
	}
}

func (qgc *QuestionGroupController) GetAllQuestionGroups(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	questionGroups, err := qgc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(qgc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, questionGroups)
}

func (qgc *QuestionGroupController) GetQuestionGroupDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionGroup, err := qgc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, questionGroup)
}

func (qgc *QuestionGroupController) CreateQuestionGroup(c *gin.Context) {

	var questionGroup *models.QuestionGroup

	if err := c.ShouldBindJSON(&questionGroup); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qgc.db, c, err.Error())
		return
	}

	createdQuestionGroup, err := qgc.service.Create(currentUserID, questionGroup)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(qgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdQuestionGroup)
}

func (qgc *QuestionGroupController) UpdateQuestionGroup(c *gin.Context) {

	questionGroupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedQuestionGroup *models.QuestionGroup
	if err := c.ShouldBindJSON(&updatedQuestionGroup); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qgc.db, c, err.Error())
		return
	}

	updatedQuestionGroup, err = qgc.service.Update(uint(questionGroupID), currentUserID, updatedQuestionGroup)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(qgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedQuestionGroup)
}

func (qgc *QuestionGroupController) DeleteQuestionGroups(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qgc.db, c, err.Error())
		return
	}

	failedIDs, err := qgc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(qgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (qgc *QuestionGroupController) ActivateQuestionGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qgc.db, c, err.Error())
		return
	}

	err = qgc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(qgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (qgc *QuestionGroupController) DeactivateQuestionGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qgc.db, c, err.Error())
		return
	}

	err = qgc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(qgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
