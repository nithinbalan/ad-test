package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
)

type ComplianceEvidenceSubmissionController struct {
	service services.ComplianceEvidenceSubmissionServiceInterface
	db      *gorm.DB
}

func NewComplianceEvidenceSubmissionController(service services.ComplianceEvidenceSubmissionServiceInterface) *ComplianceEvidenceSubmissionController {
	return &ComplianceEvidenceSubmissionController{
		service: service,
	}
}

func (cpac *ComplianceEvidenceSubmissionController) GetAllAssessments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	isPassiveUser := utility.IsPassiveAccess(c)

	var export bool
	var err error
	if val, exists := c.Get("export"); exists {
		export = val.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(cpac.db, c, "Invalid export parameter")
			return
		}
	}

	validSortFields := map[string]string{
		"id":                             "id",
		"compliance_clause_title_en":     "compliance_clause_title_en",
		"compliance_clause_title_ar":     "compliance_clause_title_ar",
		"compliance_register_title_en":   "compliance_register_title_en",
		"compliance_register_title_ar":   "compliance_register_title_ar",
		"responsible_unit_title_en":      "responsible_unit_title_en",
		"responsible_unit_title_ar":      "responsible_unit_title_ar",
		"responsible_user_first_name_en": "responsible_user_first_name_en",
		"responsible_user_first_name_ar": "responsible_user_first_name_ar",
		"created_at":                     "created_at",
		"is_completed":                   "is_completed",
		"is_acknowledged":                "is_acknowledged",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	allowedParams := []string{
		"compliance_assessment_requirement_status_ids",
		"responsible_user_ids",
		"responsible_unit_ids",
		"compliance_requirement_ids",
		"compliance_assessment_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cpac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cpac.db, c, "Invalid page_size parameter")
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cpac.db, c, err.Error())
		return
	}

	queryParams := complianceDTOs.ComplianceEvidenceSubmissionListQueryParamsDTO{
		Sort:          sort,
		SortBy:        sortField,
		Keyword:       keyword,
		Filters:       filters,
		QuickFilter:   quickFilter,
		Page:          page,
		PageSize:      pageSize,
		Export:        export,
		IsPassiveUser: isPassiveUser,
	}

	result, err := cpac.service.GetAllAssessments(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cpac.db, c, err.Error())
		return
	}

	if export {
		excelFile, err := cpac.service.Export(result.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(cpac.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Department Responses.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (cpac *ComplianceEvidenceSubmissionController) GetAllAssessmentsExport(c *gin.Context) {
	c.Set("export", true)
	cpac.GetAllAssessments(c)
}

func (cpac *ComplianceEvidenceSubmissionController) GetRequirementAssessmentFullDetails(c *gin.Context) {
	responsibleUnitIDStr := c.Param("id")
	responsibleUnitID, err := strconv.ParseUint(responsibleUnitIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	responsibleUnit, err := cpac.service.GetRequirementAssessmentFullDetailsByID(uint(responsibleUnitID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			_ = e
			utility.SendNotFoundErrResp(c)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(cpac.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, responsibleUnit)
}
