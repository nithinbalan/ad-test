package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SwotPestelQuestionFactorController struct {
	service services.SwotPestelQuestionFactorServiceInterface
	db      *gorm.DB
}

func NewSwotPestelQuestionFactorController(service services.SwotPestelQuestionFactorServiceInterface) *SwotPestelQuestionFactorController {
	return &SwotPestelQuestionFactorController{
		service: service,
	}
}

func (brpc *SwotPestelQuestionFactorController) GetAllSwotPestelQuestionFactors(c *gin.Context) {

	idStr := c.Param("questionId")

	questionId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "description")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
		"bcp_phase_ids",
		"responsible_user_ids",
	}
	validSortFields := map[string]string{
		"description":                   "swot_pestel_question_factors.description",
		"swot_pestel_question_title_en": "swot_pestel_questions.title_en",
		"swot_pestel_question_title_ar": "swot_pestel_questions.title_ar",
		"swot_pestel_category_title_en": "swot_pestel_categories.title_en",
		"swot_pestel_category_title_ar": "swot_pestel_categories.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := strategyDTOs.SwotPestelQuestionFactorListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
	swotPestelQuestionFactors, err := brpc.service.GetAll(uint(questionId), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, swotPestelQuestionFactors)
}

func (brpc *SwotPestelQuestionFactorController) GetSwotPestelQuestionFactorDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	swotPestelQuestionFactor, err := brpc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(brpc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, swotPestelQuestionFactor)
}

func (brpc *SwotPestelQuestionFactorController) CreateSwotPestelQuestionFactor(c *gin.Context) {

	idStr := c.Param("questionId")

	questionId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var swotPestelQuestionFactorFormDTO strategyDTOs.SwotPestelQuestionFactorFormDTO

	if err := c.ShouldBindJSON(&swotPestelQuestionFactorFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, err.Error())
		return
	}

	id, err := brpc.service.Create(uint(questionId), currentUserID, &swotPestelQuestionFactorFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(brpc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (brpc *SwotPestelQuestionFactorController) UpdateSwotPestelQuestionFactor(c *gin.Context) {
	var swotPestelQuestionFactorFormDTO strategyDTOs.SwotPestelQuestionFactorFormDTO

	swotPestelQuestionFactorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&swotPestelQuestionFactorFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, err.Error())
		return
	}

	err = brpc.service.Update(uint(swotPestelQuestionFactorID), currentUserID, &swotPestelQuestionFactorFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(brpc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (brpc *SwotPestelQuestionFactorController) DeleteSwotPestelQuestionFactors(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, err.Error())
		return
	}

	failedIDs, err := brpc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
