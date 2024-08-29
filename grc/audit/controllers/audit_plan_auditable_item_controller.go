package controllers

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditPlanAuditableItemController struct {
	service services.AuditPlanAuditableItemServiceInterface
	db      *gorm.DB
}

func NewAuditPlanAuditableItemController(service services.AuditPlanAuditableItemServiceInterface) *AuditPlanAuditableItemController {
	return &AuditPlanAuditableItemController{
		service: service,
	}
}

func (apaic *AuditPlanAuditableItemController) GetAllAuditPlanAuditableItems(c *gin.Context) {
	auditPlanIDStr := c.Param("planId")

	auditPlanID, err := strconv.ParseUint(auditPlanIDStr, 10, 64)
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
		"audit_category_ids",
		"audit_team_ids",
		"team_lead_ids",
		"audit_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(apaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(apaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditDTOs.AuditPlanAuditableItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditPlans, err := apaic.service.GetAll(uint(auditPlanID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(apaic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditPlans)
}

func (apaic *AuditPlanAuditableItemController) GetAllAuditPlanAuditableItemDetails(c *gin.Context) {
	auditPlanIDStr := c.Param("planId")

	_, err := strconv.ParseUint(auditPlanIDStr, 10, 64)
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

	auditableItem, err := apaic.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(apaic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, auditableItem)
}

func (apaic *AuditPlanAuditableItemController) CreateAuditPlanAuditableItem(c *gin.Context) {
	var auditPlanFormDTO auditDTOs.AuditPlanAuditableItemFormDTO
	auditPlanIDStr := c.Param("planId")

	auditPlanID, err := strconv.ParseUint(auditPlanIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apaic.db, c, err.Error())
		return
	}

	id, err := apaic.service.Create(uint(auditPlanID), currentUserID, &auditPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(apaic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (apaic *AuditPlanAuditableItemController) DeleteAuditPlanAuditableItems(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	auditPlanIDStr := c.Param("planId")

	auditPlanID, err := strconv.ParseUint(auditPlanIDStr, 10, 64)
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
		utility.SendInternalErrResp(apaic.db, c, err.Error())
		return
	}

	failedIDs, err := apaic.service.Delete(uint(auditPlanID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(apaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
