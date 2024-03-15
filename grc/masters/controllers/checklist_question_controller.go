package controllers

import (
	"jethur/grc/masters/models"
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChecklistQuestionController struct {
	service services.ChecklistQuestionServiceInterface
	db      *gorm.DB
}

func NewChecklistQuestionController(service services.ChecklistQuestionServiceInterface) *ChecklistQuestionController {
	return &ChecklistQuestionController{
		service: service,
	}
}

func (cqc *ChecklistQuestionController) GetAllChecklistQuestions(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	checklistQuestions, err := cqc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(cqc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, checklistQuestions)
}

func (cqc *ChecklistQuestionController) GetChecklistQuestionDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checklistQuestion, err := cqc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, checklistQuestion)
}

func (cqc *ChecklistQuestionController) CreateChecklistQuestion(c *gin.Context) {

	var checklistQuestion *models.ChecklistQuestion

	if err := c.ShouldBindJSON(&checklistQuestion); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cqc.db, c, err.Error())
		return
	}

	createdChecklistQuestion, err := cqc.service.Create(currentUserID, checklistQuestion)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cqc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdChecklistQuestion)
}

func (cqc *ChecklistQuestionController) UpdateChecklistQuestion(c *gin.Context) {

	checklistQuestionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedChecklistQuestion *models.ChecklistQuestion
	if err := c.ShouldBindJSON(&updatedChecklistQuestion); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cqc.db, c, err.Error())
		return
	}

	updatedChecklistQuestion, err = cqc.service.Update(uint(checklistQuestionID), currentUserID, updatedChecklistQuestion)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cqc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedChecklistQuestion)
}

func (cqc *ChecklistQuestionController) DeleteChecklistQuestions(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cqc.db, c, err.Error())
		return
	}

	failedIDs, err := cqc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cqc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
