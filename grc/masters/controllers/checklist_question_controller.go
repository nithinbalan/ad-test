package controllers

import (
	masterDTOs "jethur/grc/masters/dtos"
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

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cqc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cqc.db, c, "Invalid page_size parameter")
		return
	}

	checklistQuestions, err := cqc.service.GetAll(sort, sortBy, keyword, page, pageSize)

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

	var checklistQuestionFormDTO *masterDTOs.ChecklistQuestionFormDTO

	if err := c.ShouldBindJSON(&checklistQuestionFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cqc.db, c, err.Error())
		return
	}

	createdChecklistQuestion, err := cqc.service.Create(currentUserID, checklistQuestionFormDTO)
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

	var checklistQuestionFormDT masterDTOs.ChecklistQuestionFormDTO
	if err := c.ShouldBindJSON(&checklistQuestionFormDT); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cqc.db, c, err.Error())
		return
	}

	err = cqc.service.Update(uint(checklistQuestionID), currentUserID, &checklistQuestionFormDT)
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

	c.Status(http.StatusOK)
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
