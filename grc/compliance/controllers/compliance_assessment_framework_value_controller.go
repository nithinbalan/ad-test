package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentFrameworkValueController struct {
	service services.ComplianceAssessmentFrameworkValueServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentFrameworkValueController(service services.ComplianceAssessmentFrameworkValueServiceInterface) *ComplianceAssessmentFrameworkValueController {
	return &ComplianceAssessmentFrameworkValueController{
		service: service,
	}
}

func (cafc *ComplianceAssessmentFrameworkValueController) GetAllComplianceAssessmentFrameworkValues(c *gin.Context) {
	complianceAssessmentFrameworkIDStr := c.Param("complianceAssessmentFrameworkId")

	var complianceAssessmentFrameworkID uint64
	var err error

	// Check if projectID is provided in the URL
	if complianceAssessmentFrameworkIDStr != "" {
		complianceAssessmentFrameworkID, err = strconv.ParseUint(complianceAssessmentFrameworkIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		complianceAssessmentFrameworkID = 0
	}
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"compliance_assessment_framework_ids",
		"compliance_status_ids",
		"created_by_id",
		"updated_by_id",
	}

	validSortFields := map[string]string{
		"id":                   "compliance_assessment_framework_values.id",
		"title_en":             "compliance_assessment_framework_values.title_en",
		"title_ar":             "compliance_assessment_framework_values.title_ar",
		"score_min":            "compliance_assessment_framework_values.score_min",
		"score_max":            "compliance_assessment_framework_values.score_max",
		"value":                "compliance_assessment_framework_values.value",
		"compliance_status_id": "compliance_assessment_framework_values.compliance_status_id",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cafc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cafc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := complianceDTOs.ComplianceAssessmentFrameworkValueListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	complianceAssessmentFrameworkValues, err := cafc.service.GetAll(uint(complianceAssessmentFrameworkID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cafc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceAssessmentFrameworkValues)
}

func (cafc *ComplianceAssessmentFrameworkValueController) GetComplianceAssessmentFrameworkValueDetails(c *gin.Context) {
	complianceAssessmentFrameworkIDStr := c.Param("complianceAssessmentFrameworkId")

	complianceAssessmentFrameworkID, err := strconv.ParseUint(complianceAssessmentFrameworkIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("complianceAssessmentFrameworkValueId")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	complianceAssessmentFrameworkValue, err := cafc.service.GetFullDetailsByID(uint(complianceAssessmentFrameworkID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cafc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceAssessmentFrameworkValue)
}

func (cafc *ComplianceAssessmentFrameworkValueController) CreateComplianceAssessmentFrameworkValue(c *gin.Context) {
	var complianceAssessmentFrameworkValueFormDTO complianceDTOs.ComplianceAssessmentFrameworkValueFormDTO
	complianceAssessmentFrameworkID, err := strconv.ParseUint(c.Param("complianceAssessmentFrameworkId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&complianceAssessmentFrameworkValueFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cafc.db, c, err.Error())
		return
	}
	id, err := cafc.service.Create(uint(complianceAssessmentFrameworkID), currentUserID, &complianceAssessmentFrameworkValueFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cafc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (cafc *ComplianceAssessmentFrameworkValueController) UpdateComplianceAssessmentFrameworkValue(c *gin.Context) {
	var complianceAssessmentFrameworkValueFormDTO complianceDTOs.ComplianceAssessmentFrameworkValueFormDTO
	complianceAssessmentFrameworkID, err := strconv.ParseUint(c.Param("complianceAssessmentFrameworkId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	valueID, err := strconv.Atoi(c.Param("complianceAssessmentFrameworkValueId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&complianceAssessmentFrameworkValueFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cafc.db, c, err.Error())
		return
	}
	err = cafc.service.Update(uint(complianceAssessmentFrameworkID), uint(valueID), currentUserID, &complianceAssessmentFrameworkValueFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cafc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (cafc *ComplianceAssessmentFrameworkValueController) DeleteComplianceAssessmentFrameworkValues(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	complianceAssessmentFrameworkID, err := strconv.ParseUint(c.Param("complianceAssessmentFrameworkId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cafc.db, c, err.Error())
		return
	}
	failedIDs, err := cafc.service.Delete(uint(complianceAssessmentFrameworkID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cafc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
