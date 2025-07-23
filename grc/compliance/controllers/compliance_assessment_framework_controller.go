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

type ComplianceAssessmentFrameworkController struct {
	service services.ComplianceAssessmentFrameworkServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentFrameworkController(service services.ComplianceAssessmentFrameworkServiceInterface) *ComplianceAssessmentFrameworkController {
	return &ComplianceAssessmentFrameworkController{
		service: service,
	}
}

func (cafc *ComplianceAssessmentFrameworkController) GetAllComplianceAssessmentFrameworks(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_assessment_framework_ids",
		"compliance_assessment_framework_type_ids",
	}
	validSortFields := map[string]string{
		"title_en":    "compliance_assessment_frameworks.title_en",
		"title_ar":    "compliance_assessment_frameworks.title_ar",
		"description": "compliance_assessment_frameworks.description",
		"created_at":  "compliance_assessment_frameworks.created_at",
		"updated_at":  "compliance_assessment_frameworks.updated_at",
		"compliance_assessment_framework_type_title_en": "compliance_assessment_framework_types.title_en",
		"compliance_assessment_framework_type_title_ar": "compliance_assessment_framework_types.title_ar",
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
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := complianceDTOs.ComplianceAssessmentFrameworkListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	complianceAssessmentFramework, err := cafc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(cafc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, complianceAssessmentFramework)
}

func (cafc *ComplianceAssessmentFrameworkController) GetComplianceAssessmentFrameworkDetails(c *gin.Context) {
	idStr := c.Param("complianceAssessmentFrameworkId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceAssessmentFramework, err := cafc.service.GetFullDetailsByID(uint(id))
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
	c.JSON(http.StatusOK, complianceAssessmentFramework)
}

func (cafc *ComplianceAssessmentFrameworkController) CreateComplianceAssessmentFramework(c *gin.Context) {

	var complianceAssessmentFrameworkForm complianceDTOs.ComplianceAssessmentFrameworkFormDTO

	if err := c.ShouldBindJSON(&complianceAssessmentFrameworkForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cafc.db, c, err.Error())
		return
	}

	id, err := cafc.service.Create(currentUserID, &complianceAssessmentFrameworkForm)

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

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (cafc *ComplianceAssessmentFrameworkController) UpdateComplianceAssessmentFramework(c *gin.Context) {
	var complianceAssessmentFrameworkForm complianceDTOs.ComplianceAssessmentFrameworkFormDTO

	complianceAssessmentFrameworkID, err := strconv.Atoi(c.Param("complianceAssessmentFrameworkId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&complianceAssessmentFrameworkForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cafc.db, c, err.Error())
		return
	}

	err = cafc.service.Update(uint(complianceAssessmentFrameworkID), currentUserID, &complianceAssessmentFrameworkForm)
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

func (cafc *ComplianceAssessmentFrameworkController) DeleteComplianceAssessmentFramework(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
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

	failedIDs, err := cafc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cafc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

// func (cafc *ComplianceAssessmentFrameworkController) GetComplianceAssessmentFrameworkByType(c *gin.Context) {
// 	status := c.Query("status")

// 	complianceAssessmentFramework, err := cafc.service.GetByType(status)

// 	if err != nil {
// 		utility.SendInternalErrResp(cafc.db, c, err.Error())
// 		return
// 	}
// 	c.JSON(http.StatusOK, complianceAssessmentFramework)
// }
