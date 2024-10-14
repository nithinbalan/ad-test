package controllers

import (
	"jethur/grc/audit-insights/services"
	auditFindingDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightsDraftFindingActionPlanController struct {
	service services.AuditInsightsDraftFindingActionPlanServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsDraftFindingActionPlanController(service services.AuditInsightsDraftFindingActionPlanServiceInterface) *AuditInsightsDraftFindingActionPlanController {
	return &AuditInsightsDraftFindingActionPlanController{
		service: service,
	}
}

func (aidfapc *AuditInsightsDraftFindingActionPlanController) GetAllAuditInsightsDraftFindingActionPlans(c *gin.Context) {

	var findingID uint64
	var err error

	findingIDStr := c.Param("findingId")

	findingID, err = strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

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
		utility.SendInternalErrResp(aidfapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aidfapc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditFindingDTOs.AuditDraftFindingActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortBy,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
	}
	
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidfapc.db, c, err.Error())
		return
	}
	auditInsightsDraftFindingActionPlans, err := aidfapc.service.GetAll(uint(findingID), currentUserID, &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aidfapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditInsightsDraftFindingActionPlans)
}

func (aidfapc *AuditInsightsDraftFindingActionPlanController) GetAuditInsightsDraftFindingActionPlanDetails(c *gin.Context) {

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditInsightsDraftFindingActionPlan, err := aidfapc.service.GetFullDetailsByID(uint(findingID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aidfapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, auditInsightsDraftFindingActionPlan)
}

func (aidfapc *AuditInsightsDraftFindingActionPlanController) CreateAuditInsightsDraftFindingActionPlan(c *gin.Context) {
	var AuditInsightsDraftFindingActionPlanFormDTO auditFindingDTOs.AuditDraftFindingActionPlanFormDTO

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&AuditInsightsDraftFindingActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidfapc.db, c, err.Error())
		return
	}

	id, err := aidfapc.service.Create(uint(findingID), currentUserID, &AuditInsightsDraftFindingActionPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aidfapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (aidfapc *AuditInsightsDraftFindingActionPlanController) UpdateAuditInsightsDraftFindingActionPlan(c *gin.Context) {
	var AuditInsightsDraftFindingActionPlanFormDTO auditFindingDTOs.AuditDraftFindingActionPlanFormDTO

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&AuditInsightsDraftFindingActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidfapc.db, c, err.Error())
		return
	}

	err = aidfapc.service.Update(uint(findingID), uint(id), currentUserID, &AuditInsightsDraftFindingActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aidfapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (aidfapc *AuditInsightsDraftFindingActionPlanController) DeleteAuditInsightsDraftFindingActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
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
		utility.SendInternalErrResp(aidfapc.db, c, err.Error())
		return
	}

	failedIDs, err := aidfapc.service.Delete(uint(findingID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aidfapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
