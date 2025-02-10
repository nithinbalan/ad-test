package controllers

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementController struct {
	service services.ComplianceAssessmentRequirementServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentRequirementController(service services.ComplianceAssessmentRequirementServiceInterface) *ComplianceAssessmentRequirementController {
	return &ComplianceAssessmentRequirementController{
		service: service,
	}
}

func (carc *ComplianceAssessmentRequirementController) GetComplianceAssessmentRequirements(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "compliance_clause_no")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	validSortFields := map[string]string{
		"compliance_clause_no":                     "compliance_clauses.clause_no",
		"compliance_clause_title_en":               "compliance_clauses.title_en",
		"compliance_clause_title_ar":               "compliance_clauses.title_ar",
		"compliance_requirement_title_en":          "compliance_requirements.title_en",
		"compliance_requirement_title_ar":          "compliance_requirements.title_ar",
		"compliance_requirement_category_title_en": "compliance_requirement_categories.title_en",
		"compliance_requirement_category_title_ar": "compliance_requirement_categories.title_ar",
		"assessment_completed":                     "compliance_assessment_requirements.is_completed",
		"severity_level_en":                        "severity_levels.title_en",
		"severity_level_ar":                        "severity_levels.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_requirement_ids",
		"compliance_status_ids",
		"compliance_clause_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	requirements, err := carc.service.GetAll(uint(assessmentID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(carc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, requirements)
}

func (cac *ComplianceAssessmentRequirementController) GetComplianceAssessmentRequirementDetails(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIDStr := c.Param("requirementId")
	complianceRequirementID, err := strconv.ParseUint(requirementIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	complianceAssessment, err := cac.service.GetFullDetailsByID(uint(assessmentID), uint(complianceRequirementID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceAssessment)
}

func (cac *ComplianceAssessmentRequirementController) UpdateComplianceAssessmentRequirement(c *gin.Context) {
	var complianceAssessmentRequirementUpdateFormDTO complianceDTOs.ComplianceAssessmentRequirementUpdateFormDTO

	complianceAssessmentID, err := strconv.Atoi(c.Param("assessmentId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIDStr := c.Param("requirementId")

	complianceRequirementID, err := strconv.ParseUint(requirementIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&complianceAssessmentRequirementUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	err = cac.service.Update(uint(complianceAssessmentID), uint(complianceRequirementID), currentUserID, &complianceAssessmentRequirementUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
