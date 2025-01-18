package controllers

import (
	"errors"
	auditFindingDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFindingActionPlanController struct {
	service services.AuditFindingActionPlanServiceInterface
	db      *gorm.DB
}

func NewAuditFindingActionPlanController(service services.AuditFindingActionPlanServiceInterface) *AuditFindingActionPlanController {
	return &AuditFindingActionPlanController{
		service: service,
	}
}

func (afapc *AuditFindingActionPlanController) GetAllAuditFindingActionPlans(c *gin.Context) {
	var findingID *uint64
	var err error

	// Retrieve the finding ID from the request, if provided
	findingIDStr := c.Param("findingId")
	if findingIDStr != "" {
		parsedID, err := strconv.ParseUint(findingIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
		findingID = &parsedID
	}

	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(afapc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
		"responsible_user_ids",
		"action_plan_status_ids",
		"audit_finding_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afapc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"id":                             "audit_finding_action_plans.id",
		"audit_finding_id":               "audit_finding_action_plans.audit_finding_id",
		"audit_finding_reference_number": "audit_findings.reference_number",
		"audit_finding_title":            "audit_findings.title",
		"audit_finding_description":      "audit_findings.description",
		"title":                          "audit_finding_action_plans.title",
		"description":                    "audit_finding_action_plans.description",
		"severity_level_title_en":        "severity_levels.title_en",
		"severity_level_title_ar":        "severity_levels.title_ar",
		"responsible_user_id":            "responsible_users.id",
		"responsible_user_first_name_en": "responsible_users.first_name_en",
		"responsible_user_first_name_ar": "responsible_users.first_name_ar",
		"responsible_user_last_name_en":  "responsible_users.last_name_en",
		"responsible_user_last_name_ar":  "responsible_users.last_name_ar",
		"responsible_user_email":         "responsible_users.email",
		"review_user_id":                 "review_users.id",
		"review_user_first_name_en":      "review_users.first_name_en",
		"review_user_first_name_ar":      "review_users.first_name_ar",
		"review_user_last_name_en":       "review_users.last_name_en",
		"review_user_last_name_ar":       "review_users.last_name_ar",
		"review_user_email":              "review_users.email",
		"target_date":                    "audit_finding_action_plans.target_date",
		"percentage":                     "audit_finding_action_plans.percentage",
		"comments":                       "audit_finding_action_plans.comments",
		"action_plan_status_title_en":    "action_plan_statuses.title_en",
		"action_plan_status_title_ar":    "action_plan_statuses.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := auditFindingDTOs.AuditFindingActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortField,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
		Export:         export,
	}

	// Handle the optional findingID
	var findingIDValue uint
	if findingID != nil {
		findingIDValue = uint(*findingID)
	}

	auditFindingActionPlans, err := afapc.service.GetAll(findingIDValue, &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}

	// If export is true, generate and return the Excel file
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(afapc.db, c, err.Error())
			return
		}

		excelFile, err := afapc.service.ExportAuditFindingActionPlans(auditFindingActionPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(afapc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Audit Finding Action Plans.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, auditFindingActionPlans)
}

func (afapc *AuditFindingActionPlanController) GetAuditFindingActionPlanExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	afapc.GetAllAuditFindingActionPlans(c)
}

func (afapc *AuditFindingActionPlanController) GetAuditFindingActionPlanDetails(c *gin.Context) {

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditFindingActionPlan, err := afapc.service.GetFullDetailsByID(uint(findingID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, auditFindingActionPlan)
}

func (afapc *AuditFindingActionPlanController) CreateAuditFindingActionPlan(c *gin.Context) {
	var AuditFindingActionPlanFormDTO auditFindingDTOs.AuditFindingActionPlanFormDTO

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&AuditFindingActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afapc.db, c, err.Error())
		return
	}

	id, err := afapc.service.Create(uint(findingID), currentUserID, &AuditFindingActionPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (afapc *AuditFindingActionPlanController) UpdateAuditFindingActionPlan(c *gin.Context) {
	var AuditFindingActionPlanFormDTO auditFindingDTOs.AuditFindingActionPlanFormDTO

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&AuditFindingActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afapc.db, c, err.Error())
		return
	}

	err = afapc.service.Update(uint(findingID), uint(id), currentUserID, &AuditFindingActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (afapc *AuditFindingActionPlanController) DeleteAuditFindingActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	findingIDStr := c.Param("findingId")

	_, err := strconv.ParseUint(findingIDStr, 10, 64)
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
		utility.SendInternalErrResp(afapc.db, c, err.Error())
		return
	}

	failedIDs, err := afapc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (afapc *AuditFindingActionPlanController) GetAuditFindingActionPlanCountByStatus(c *gin.Context) {
	findingIDStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanCounts, err := afapc.service.GetAuditFindingActionPlanCountByStatus(uint(findingId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, actionPlanCounts)
}
