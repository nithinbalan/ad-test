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

type ComplianceRequirementController struct {
	service services.ComplianceRequirementServiceInterface
	db      *gorm.DB
}

func NewComplianceRequirementController(service services.ComplianceRequirementServiceInterface) *ComplianceRequirementController {
	return &ComplianceRequirementController{
		service: service,
	}
}

func (crc *ComplianceRequirementController) GetAllComplianceRequirements(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	var complianceID uint64
	var err error

	// Check if complianceID is provided in the URL
	if complianceIDStr != "" {
		complianceID, err = strconv.ParseUint(complianceIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		complianceID = 0
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	exportFromContext, exists := c.Get("export")
	var export bool

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(crc.db, c, "Invalid export parameter")
			return
		}
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_register_ids",
		"compliance_requirement_category_ids",
		"compliance_clause_ids",
		"severity_level_ids",
		"compliance_status_ids",
		"responsible_user_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse deadline_from, deadline_to
	deadlineFrom := c.Query("deadline_from")
	deadlineTo := c.Query("deadline_to")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := complianceDTOs.ComplianceRequirementListQueryParamsDTO{
		Sort:         sort,
		SortBy:       sortBy,
		Keyword:      keyword,
		Status:       status,
		Page:         page,
		PageSize:     pageSize,
		Filters:      filters,
		QuickFilter:  quickFilter,
		DeadlineFrom: deadlineFrom,
		DeadlineTo:   deadlineTo,
		DateFrom:     dateFrom,
		DateTo:       dateTo,
		Export:       export,
	}
	if isRegulationNotificationAssessmentStr := c.Query("is_regulation_notification_assessment"); isRegulationNotificationAssessmentStr != "" {
		val, err := strconv.ParseBool(isRegulationNotificationAssessmentStr)
		if err != nil {
			utility.SendInvalidErrResp(c, err)
			return
		}
		queryParams.IsRegulationNotificationAssessment = &val
	}
	if isRiskBasedAssessmentStr := c.Query("is_risk_based_assessment"); isRiskBasedAssessmentStr != "" {
		val, err := strconv.ParseBool(isRiskBasedAssessmentStr)
		if err != nil {
			utility.SendInvalidErrResp(c, err)
			return
		}
		queryParams.IsRiskBasedAssessment = &val
	}

	complianceRequirements, err := crc.service.GetAll(uint(complianceID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	// If export is true, generate and return the Excel file
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(crc.db, c, err.Error())
			return
		}

		excelFile, err := crc.service.ExportComplianceRequirements(complianceRequirements.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(crc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="ComplianceRequirements.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, complianceRequirements)
}

func (crc *ComplianceRequirementController) GetComplianceRequirementsByID(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "compliance_clauses.clause_no")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_register_ids",
		"compliance_requirement_category_ids",
		"compliance_clause_ids",
		"severity_level_ids",
		"compliance_status_ids",
		"responsible_user_ids",
		"responsible_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse deadline_from, deadline_to
	deadlineFrom := c.Query("deadline_from")
	deadlineTo := c.Query("deadline_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"compliance_clause_no":                     "compliance_clauses.clause_no",
		"compliance_clause_title_en":               "compliance_clauses.title_en",
		"compliance_clause_title_ar":               "compliance_clauses.title_ar",
		"title_en":                                 "compliance_requirements.title_en",
		"title_ar":                                 "compliance_requirements.title_ar",
		"compliance_requirement_category_title_en": "compliance_requirement_categories.title_en",
		"compliance_requirement_category_title_ar": "compliance_requirement_categories.title_ar",
		"severity_level_title_en":                  "severity_levels.title_en",
		"severity_level_title_ar":                  "severity_levels.title_ar",
		"deadline":                                 "compliance_requirements.deadline",
		"compliance_status_title_en":               "compliance_statuses.title_en",
		"compliance_status_title_ar":               "compliance_statuses.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := complianceDTOs.ComplianceRequirementListQueryParamsDTO{
		Sort:         sort,
		SortBy:       sortField,
		Keyword:      keyword,
		Status:       status,
		Page:         page,
		PageSize:     pageSize,
		Filters:      filters,
		QuickFilter:  quickFilter,
		DeadlineFrom: deadlineFrom,
		DeadlineTo:   deadlineTo,
	}
	if isRegulationNotificationAssessmentStr := c.Query("is_regulation_notification_assessment"); isRegulationNotificationAssessmentStr != "" {
		val, err := strconv.ParseBool(isRegulationNotificationAssessmentStr)
		if err != nil {
			utility.SendInvalidErrResp(c, err)
			return
		}
		queryParams.IsRegulationNotificationAssessment = &val
	}
	if isRiskBasedAssessmentStr := c.Query("is_risk_based_assessment"); isRiskBasedAssessmentStr != "" {
		val, err := strconv.ParseBool(isRiskBasedAssessmentStr)
		if err != nil {
			utility.SendInvalidErrResp(c, err)
			return
		}
		queryParams.IsRiskBasedAssessment = &val
	}

	complianceRequirements, err := crc.service.GetRequirementsByComplianceID(uint(complianceID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceRequirements)
}

func (crc *ComplianceRequirementController) GetComplianceRequirementDetails(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("requirementId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceRequirement, err := crc.service.GetFullDetailsByID(uint(complianceID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceRequirement)
}

func (crc *ComplianceRequirementController) CreateComplianceRequirement(c *gin.Context) {
	var ComplianceRequirementFormDTO complianceDTOs.ComplianceRequirementFormDTO

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ComplianceRequirementFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	id, err := crc.service.Create(uint(complianceID), currentUserID, &ComplianceRequirementFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (crc *ComplianceRequirementController) UpdateComplianceRequirement(c *gin.Context) {
	var ComplianceRequirementFormDTO complianceDTOs.ComplianceRequirementFormDTO

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("requirementId")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ComplianceRequirementFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.Update(uint(complianceID), uint(id), currentUserID, &ComplianceRequirementFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (crc *ComplianceRequirementController) DeleteComplianceRequirements(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	failedIDs, err := crc.service.Delete(uint(complianceID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (crc *ComplianceRequirementController) GetComplianceRequirementsExport(c *gin.Context) {
	c.Set("export", true)

	crc.GetAllComplianceRequirements(c)
}

func (crc *ComplianceRequirementController) GetRequirementsTemplate(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	file, err := crc.service.GenerateRequirementTemplate(uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Failed to generate template")
		return
	}
	defer file.Close()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", `attachment; filename="ComplianceRequirementTemplate.xlsx"`)
	c.Header("Content-Type", "application/octet-stream")
	c.File(file.Name())
}

func (crc *ComplianceRequirementController) ValidateRequirementTemplate(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Failed to read uploaded file")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Failed to open file")
		return
	}
	defer file.Close()

	validationResult, err := crc.service.ValidateRequirementTemplate(uint(complianceID), file)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, validationResult)
}

func (crc *ComplianceRequirementController) ValidateImportRequirement(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var input []map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	validationResult, err := crc.service.ValidateImportRequirement(uint(complianceID), input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, validationResult)
}

func (crc *ComplianceRequirementController) ImportRequirements(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	createdIDs, err := crc.service.ImportRequirements(uint(complianceID), input, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"created_ids": createdIDs})
}

// CreateInherentRiskAnalysis creates inherent risk analysis for a compliance requirement
func (crc *ComplianceRequirementController) CreateInherentRiskAnalysis(c *gin.Context) {
	requirementIDStr := c.Param("requirementId")
	requirementID, err := strconv.ParseUint(requirementIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, utility.NewValidationError(err, "invalid", "requirement_id"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	var riskAnalysisFormDTO []*complianceDTOs.ComplianceRequirementInherentRiskFormDTO
	if err := c.ShouldBindJSON(&riskAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = crc.service.CreateInherentRiskAnalysis(uint(requirementID), currentUserID, riskAnalysisFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

// DeleteInherentRiskAnalysis deletes all inherent risk analysis for a compliance requirement
func (crc *ComplianceRequirementController) DeleteInherentRiskAnalysis(c *gin.Context) {
	requirementIDStr := c.Param("requirementId")
	requirementID, err := strconv.ParseUint(requirementIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, utility.NewValidationError(err, "invalid", "requirement_id"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.DeleteInherentRiskAnalysis(uint(requirementID), currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
