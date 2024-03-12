package controllers

import (
	"jethur/grc/libraries/models"
	"jethur/grc/libraries/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionLibraryController struct {
	service services.QuestionLibraryServiceInterface
	db      *gorm.DB
}

func NewQuestionLibraryController(service services.QuestionLibraryServiceInterface) *QuestionLibraryController {
	return &QuestionLibraryController{
		service: service,
	}
}

func (qlc *QuestionLibraryController) GetAllQuestionLibrary(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	questionLibrary, err := qlc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(qlc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, questionLibrary)
}

func (qlc *QuestionLibraryController) GetQuestionLibraryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionLibrary, err := qlc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, questionLibrary)
}

func (qlc *QuestionLibraryController) CreateQuestionLibrary(c *gin.Context) {

	var questionLibrary *models.QuestionLibrary

	if err := c.ShouldBindJSON(&questionLibrary); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qlc.db, c, err.Error())
		return
	}

	createdQuestionLibrary, err := qlc.service.Create(currentUserID, questionLibrary)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(qlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdQuestionLibrary)
}

func (qlc *QuestionLibraryController) UpdateQuestionLibrary(c *gin.Context) {

	questionLibraryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedQuestionLibrary *models.QuestionLibrary
	if err := c.ShouldBindJSON(&updatedQuestionLibrary); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qlc.db, c, err.Error())
		return
	}

	updatedQuestionLibrary, err = qlc.service.Update(uint(questionLibraryID), currentUserID, updatedQuestionLibrary)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(qlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedQuestionLibrary)
}

func (qlc *QuestionLibraryController) DeleteQuestionLibraries(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qlc.db, c, err.Error())
		return
	}

	failedIDs, err := qlc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(qlc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
