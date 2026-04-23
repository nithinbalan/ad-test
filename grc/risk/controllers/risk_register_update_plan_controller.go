package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanController struct {
	service services.RiskRegisterUpdatePlanServiceInterface
	db      *gorm.DB
}

func NewRiskRegisterUpdatePlanController(service services.RiskRegisterUpdatePlanServiceInterface) *RiskRegisterUpdatePlanController {
	return &RiskRegisterUpdatePlanController{
		service: service,
	}
}

func (rrupc *RiskRegisterUpdatePlanController) GetAllRiskRegisterUpdatePlan(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
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
			utility.SendInternalErrResp(rrupc.db, c, "Invalid export parameter")
			return
		}
	}

	// Allowed filters
	allowedParams := []string{
		"division_ids",
		"business_unit_ids",
		"risk_register_update_plan_status_ids",
		"risk_champion_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"created_at":             "risk_register_update_plans.created_at",
		"updated_at":             "risk_register_update_plans.updated_at",
		"division_title_en":      "div.title_en",
		"division_title_ar":      "div.title_ar",
		"business_unit_title_en": "bu.title_en",
		"business_unit_title_ar": "bu.title_ar",
		"status":                 "status.title_en",
		"risk_champion_title_en": "champion.first_name_en",
		"risk_champion_title_ar": "champion.first_name_ar",
		"start_date":             "risk_register_update_plans.start_date",
		"end_date":               "risk_register_update_plans.end_date",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO{
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

	response, err := rrupc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	// If export is true, generate and return Excel file
	if export {
		// TODO: Implement export functionality
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Export not implemented yet"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (rrupc *RiskRegisterUpdatePlanController) GetRiskRegisterUpdatePlanDetailsByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	plan, err := rrupc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, plan)
}

func (rrupc *RiskRegisterUpdatePlanController) CreateRiskRegisterUpdatePlan(c *gin.Context) {
	var formDTO riskDTOs.RiskRegisterUpdatePlanFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	planID, err := rrupc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: planID})
}

func (rrupc *RiskRegisterUpdatePlanController) UpdateRiskRegisterUpdatePlan(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.RiskRegisterUpdatePlanFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrupc *RiskRegisterUpdatePlanController) DeleteRiskRegisterUpdatePlan(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	// Call service
	failedIDs, err := rrupc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"failed_ids": failedIDs,
	})
}

// ERM Department Evaluation.
func (rrupc *RiskRegisterUpdatePlanController) CreateERMDepartmentEvaluation(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	response, err := rrupc.service.CreateERMDepartmentEvaluation(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func (rrupc *RiskRegisterUpdatePlanController) UpdateERMDepartmentEvaluation(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.MapRisksToUpdatePlanDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	response, err := rrupc.service.UpdateERMDepartmentEvaluation(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, response)
}

func (rrupc *RiskRegisterUpdatePlanController) GetERMDepartmentEvaluation(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Parse query parameters
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	valueStr := c.Query("value")

	exportFromContext, exists := c.Get("export")
	var export bool
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(rrupc.db, c, "Invalid export parameter")
			return
		}
	}

	// Define allowed filter parameters
	allowedParams := []string{
		"risk_category_ids",
		"risk_source_ids",
		"responsible_unit_ids",
		"responsible_user_ids",
		"risk_status_ids",
		"review_frequency_ids",
		"inherent_risk_rating_ids",
		"residual_risk_rating_ids",
		"current_likelihood_ids",
		"current_impact_ids",
		"current_risk_rating_ids",
		"risk_response_strategy_ids",
	}

	validSortFields := map[string]string{
		"reference_number":          "risk_register.reference_number",
		"title":                     "risk_register.title",
		"description":               "risk_register.description",
		"risk_category_en":          "risk_categories.title_en",
		"risk_category_ar":          "risk_categories.title_ar",
		"organization_unit_en":      "organization_units.title_en",
		"organization_unit_ar":      "organization_units.title_ar",
		"responsible_user_en":       "users.first_name_en",
		"responsible_user_ar":       "users.first_name_ar",
		"risk_source_en":            "risk_sources.title_en",
		"risk_source_ar":            "risk_sources.title_ar",
		"risk_status_en":            "risk_statuses.title_en",
		"risk_status_ar":            "risk_statuses.title_ar",
		"inherent_likelihood_en":    "inherent_likelihood.title_en",
		"inherent_likelihood_ar":    "inherent_likelihood.title_ar",
		"inherent_impact_en":        "inherent_impact.title_en",
		"inherent_impact_ar":        "inherent_impact.title_ar",
		"inherent_risk_rating_en":   "inherent_risk_rating.title_en",
		"inherent_risk_rating_ar":   "inherent_risk_rating.title_ar",
		"risk_response_strategy_en": "risk_response_strategies.title_en",
		"risk_response_strategy_ar": "risk_response_strategies.title_ar",
		"residual_likelihood_en":    "residual_likelihood.title_en",
		"residual_likelihood_ar":    "residual_likelihood.title_ar",
		"residual_impact_en":        "residual_impact.title_en",
		"residual_impact_ar":        "residual_impact.title_ar",
		"residual_risk_rating_en":   "residual_risk_rating.title_en",
		"residual_risk_rating_ar":   "residual_risk_rating.title_ar",
		"reported_at":               "risk_register.reported_at",
		"review_frequency_en":       "review_frequencies.title_en",
		"review_frequency_ar":       "review_frequencies.title_ar",
		"next_review_date":          "risk_register.next_review_date",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	reportedAtFrom := c.Query("reported_at_from")
	reportedAtTo := c.Query("reported_at_to")
	nextReviewDateFrom := c.Query("next_review_date_from")
	nextReviewDateTo := c.Query("next_review_date_to")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	riskAppetiteToleranceLevelIDStr := c.Query("risk_appetite_tolerance_level_id")
	var riskAppetiteToleranceLevelID uint
	if riskAppetiteToleranceLevelIDStr != "" {
		parsedID, err := strconv.ParseUint(riskAppetiteToleranceLevelIDStr, 10, 32)
		if err != nil {
			utility.SendInternalErrResp(rrupc.db, c, "Invalid RiskAppetiteToleranceLevelID parameter")
			return
		}
		riskAppetiteToleranceLevelID = uint(parsedID)
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(rrupc.db, c, "Invalid value parameter")
			return
		}
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.ERMDeptEvaluationListQueryParamsDTO{
		Sort:                         sort,
		SortBy:                       dbField,
		Keyword:                      keyword,
		Page:                         page,
		PageSize:                     pageSize,
		Value:                        value,
		Filters:                      filters,
		QuickFilter:                  quickFilter,
		ReportedAtFrom:               reportedAtFrom,
		ReportedAtTo:                 reportedAtTo,
		NextReviewDateFrom:           nextReviewDateFrom,
		NextReviewDateTo:             nextReviewDateTo,
		Export:                       export,
		DateFrom:                     dateFrom,
		DateTo:                       dateTo,
		RiskAppetiteToleranceLevelID: riskAppetiteToleranceLevelID,
	}

	mappedRisks, err := rrupc.service.GetERMDepartmentEvaluation(uint(id), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, mappedRisks)
}

func (rrupc *RiskRegisterUpdatePlanController) DeleteERMDepartmentEvaluation(c *gin.Context) {
	idStr := c.Param("ermEvaluationId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.DeleteERMDepartmentEvaluation(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

// Introductory Meeting Document Operations

func (rrupc *RiskRegisterUpdatePlanController) GetIntroductoryMeetingDocumentPreview(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := rrupc.service.GetIntroductoryMeetingDocumentPath(uint(planID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
		}
		return
	}

	c.File(documentPath)
}

func (rrupc *RiskRegisterUpdatePlanController) DownloadIntroductoryMeetingDocument(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := rrupc.service.GetIntroductoryMeetingDocumentPath(uint(planID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
		}
		return
	}

	filename := filepath.Base(documentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}

func (rrupc *RiskRegisterUpdatePlanController) GetIntroductoryMeetingDocuments(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documents, err := rrupc.service.GetIntroductoryMeetingDocuments(uint(planID))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

func (rrupc *RiskRegisterUpdatePlanController) CreateIntroductoryMeetingDocuments(c *gin.Context) {
	var formDTO riskDTOs.IntroductoryMeetingDocumentFormDTO

	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.CreateIntroductoryMeetingDocuments(uint(planID), currentUserID, formDTO.Documents)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Documents updated successfully"})
}

// Business Context Operations

func (rrupc *RiskRegisterUpdatePlanController) CreateBusinessContext(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.CreateBusinessContext(uint(planID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Business contexts created successfully"})
}

func (rrupc *RiskRegisterUpdatePlanController) GetBusinessContextByPlanID(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	businessContexts, err := rrupc.service.GetBusinessContextByPlanID(uint(planID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"business_contexts": businessContexts})
}

func (rrupc *RiskRegisterUpdatePlanController) UpdateBusinessContextStatus(c *gin.Context) {
	businessContextIDStr := c.Param("businessContextId")
	businessContextID, err := strconv.ParseUint(businessContextIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.UpdateBusinessContextStatusFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.UpdateBusinessContextStatus(uint(businessContextID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Business context status updated successfully"})
}
