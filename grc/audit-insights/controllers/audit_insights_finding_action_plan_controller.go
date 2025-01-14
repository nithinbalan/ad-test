package controllers

import (
	"errors"
	"jethur/grc/audit-insights/services"
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightsFindingActionPlanController struct {
	service services.AuditInsightsFindingActionPlanServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsFindingActionPlanController(service services.AuditInsightsFindingActionPlanServiceInterface) *AuditInsightsFindingActionPlanController {
	return &AuditInsightsFindingActionPlanController{
		service: service,
	}
}

func (aifapc *AuditInsightsFindingActionPlanController) GetAllAuditInsightsFindingActionPlans(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	// Check if `findingId` exists in the route
	findingIDStr := c.Param("findingId")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_finding_ids",
		"audit_ids",
		"auditable_item_ids",
		"audit_finding_category_ids",
		"severity_level_ids",
		"audit_finding_status_ids",
	}

	// If `findingId` exists in the route, include "audit_finding_ids" in allowedParams
	if findingIDStr != "" {
		allowedParams = append(allowedParams, "audit_finding_ids")
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Directly add the `findingId` filter as a string (no need to parse it to uint64)
	if findingIDStr != "" {
		filters["audit_finding_ids"] = []string{findingIDStr}
	}
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"reference_number": "audit_findings.reference_number",
		"title":            "audit_findings.title",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := auditDTOs.AuditFindingActionPlanListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, err.Error())
		return
	}
	audits, err := aifapc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, audits)
}

func (aifapc *AuditInsightsFindingActionPlanController) GetAuditInsightsFindingActionPlanDetails(c *gin.Context) {
	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	audit, err := aifapc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aifapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, audit)
}

func (aifapc *AuditInsightsFindingActionPlanController) CreateAuditInsightsFindingActionPlan(c *gin.Context) {
	var auditInsightsFindingActionPlanFormDTO auditDTOs.AuditFindingActionPlanFormDTO
	if err := c.ShouldBindJSON(&auditInsightsFindingActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get findingId from route if available
	var findingId uint64
	if findingIDStr := c.Param("findingId"); findingIDStr != "" {
		var err error
		findingId, err = strconv.ParseUint(findingIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		// Otherwise, use the value from the request body
		findingId = uint64(*auditInsightsFindingActionPlanFormDTO.FindingID)
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, err.Error())
		return
	}

	id, err := aifapc.service.Create(uint(findingId), currentUserID, &auditInsightsFindingActionPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aifapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (aifapc *AuditInsightsFindingActionPlanController) UpdateAuditInsightsFindingActionPlan(c *gin.Context) {
	var auditInsightsFindingActionPlanFormDTO auditDTOs.AuditFindingActionPlanFormDTO

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditInsightsFindingActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	// Get findingId from route if available
	var findingId uint64
	if findingIDStr := c.Param("findingId"); findingIDStr != "" {
		var err error
		findingId, err = strconv.ParseUint(findingIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		// Otherwise, use the value from the request body
		findingId = uint64(*auditInsightsFindingActionPlanFormDTO.FindingID)
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, err.Error())
		return
	}

	err = aifapc.service.Update(uint(findingId), uint(id), currentUserID, &auditInsightsFindingActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aifapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (aifapc *AuditInsightsFindingActionPlanController) DeleteAuditInsightsFindingActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, err.Error())
		return
	}

	failedIDs, err := aifapc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aifapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (aifapc *AuditInsightsFindingActionPlanController) GetAuditInsightsFindingActionPlanCountByStatus(c *gin.Context) {
	findingIDStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanCounts, err := aifapc.service.GetAuditInsightsFindingActionPlanCountByStatus(uint(findingId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aifapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, actionPlanCounts)
}
