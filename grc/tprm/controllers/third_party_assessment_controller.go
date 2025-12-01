package controllers

import (
	"errors"
	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyAssessmentController struct {
	service services.ThirdPartyAssessmentServiceInterface
	db      *gorm.DB
}

func NewThirdPartyAssessmentController(service services.ThirdPartyAssessmentServiceInterface) *ThirdPartyAssessmentController {
	return &ThirdPartyAssessmentController{
		service: service,
	}
}

func (tprc *ThirdPartyAssessmentController) GetAllThirdPartyAssessment(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	exportFromContext, exists := c.Get("export")
	var export bool
	var err error

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(tprc.db, c, "Invalid export parameter")
			return
		}
	}
	allowedParams := []string{
		"third_party_register_ids", "question_set_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page parameter"))
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page_size parameter"))
		return
	}

	validSortFields := map[string]string{
		"third_party_register_title": "third_party_register.title",
		"title":                      "third_party_assessments.title",
		"reference_number":           "third_party_assessments.reference_number",
		"due_date":                   "third_party_assessments.due_date",
		"question_set_title_en":      "question_sets.title_en",
		"question_set_title_ar":      "question_sets.title_ar",
		"created_user_first_name_en": "created_user.first_name_en",
		"created_user_last_name_en":  "created_user.last_name_en",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid sort_by parameter"))
		return
	}

	query := dtos.ThirdPartyAssessmentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}

	thirdParties, err := tprc.service.GetAll(&query)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, thirdParties)
}

func (tprc *ThirdPartyAssessmentController) GetThirdPartyAssessmentExport(c *gin.Context) {
	c.Set("export", true)

	tprc.GetAllThirdPartyAssessment(c)
}

func (tprc *ThirdPartyAssessmentController) GetThirdPartyAssessmentDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	result, err := tprc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (tprc *ThirdPartyAssessmentController) CreateThirdPartyAssessment(c *gin.Context) {
	var form dtos.ThirdPartyAssessmentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	id, err := tprc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (tprc *ThirdPartyAssessmentController) UpdateThirdPartyAssessment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ThirdPartyAssessmentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	err = tprc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tprc *ThirdPartyAssessmentController) DeleteThirdPartyAssessment(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	failedIDs, err := tprc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tprc *ThirdPartyAssessmentController) ShareLink(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	err = tprc.service.ShareLink(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
