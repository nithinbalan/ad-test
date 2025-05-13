package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionSetController struct {
	service services.QuestionSetServiceInterface
	db      *gorm.DB
}

func NewQuestionSetController(service services.QuestionSetServiceInterface) *QuestionSetController {
	return &QuestionSetController{service: service}
}

// GET /question-sets
func (qsc *QuestionSetController) GetAllQuestionSets(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")

	allowedParams := []string{
		"third_party_question_category_ids",
		"third_party_question_library_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("Invalid page parameter"))
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("Invalid page_size parameter"))
		return
	}

	validSortFields := map[string]string{
		"title_en": "question_sets.title_en",
		"title_ar": "question_sets.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("Invalid sort_by parameter"))
		return
	}

	query := dtos.QuestionSetListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	result, err := qsc.service.GetAll(&query)
	if err != nil {
		utility.SendInternalErrResp(qsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// GET /question-sets/:id
func (qsc *QuestionSetController) GetQuestionSetDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := qsc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(qsc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

// POST /question-sets
func (qsc *QuestionSetController) CreateQuestionSet(c *gin.Context) {
	var form dtos.QuestionSetFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qsc.db, c, err.Error())
		return
	}

	id, err := qsc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(qsc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// PUT /question-sets/:id
func (qsc *QuestionSetController) UpdateQuestionSet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.QuestionSetFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qsc.db, c, err.Error())
		return
	}

	err = qsc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(qsc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

// DELETE /question-sets
func (qsc *QuestionSetController) DeleteQuestionSets(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qsc.db, c, err.Error())
		return
	}

	failedIDs, err := qsc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(qsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (qsc *QuestionSetController) AddQuestionsToSet(c *gin.Context) {
	questionSetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("Invalid question_set_id"))
		return
	}

	var form dtos.QuestionSetAddQuestionsDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(qsc.db, c, err.Error())
		return
	}

	err = qsc.service.AddQuestionsToSet(uint(questionSetID), currentUserID, &form)
	if err != nil {
		utility.SendInternalErrResp(qsc.db, c, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (qsc *QuestionSetController) DeleteQuestionFromSet(c *gin.Context) {
	questionSetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("Invalid question_set_id"))
		return
	}
	questionID, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("Invalid question id"))
		return
	}

	err = qsc.service.DeleteQuestionFromSet(uint(questionSetID), uint(questionID))
	if err != nil {
		utility.SendInternalErrResp(qsc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
