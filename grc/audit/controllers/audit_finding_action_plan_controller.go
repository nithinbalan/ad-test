package controllers

import (
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
		utility.SendInternalErrResp(afapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afapc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditFindingDTOs.AuditFindingActionPlanListQueryParamsDTO{
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
	auditFindingActionPlans, err := afapc.service.GetAll(uint(findingID), &queryParams)

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
	c.JSON(http.StatusOK, auditFindingActionPlans)
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
		utility.SendInternalErrResp(afapc.db, c, err.Error())
		return
	}

	failedIDs, err := afapc.service.Delete(uint(findingID), request.IDs, currentUserID)
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
