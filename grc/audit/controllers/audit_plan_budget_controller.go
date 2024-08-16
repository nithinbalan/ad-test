package controllers

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditPlanBudgetController struct {
	service services.AuditPlanBudgetServiceInterface
	db      *gorm.DB
}

func NewAuditPlanBudgetController(service services.AuditPlanBudgetServiceInterface) *AuditPlanBudgetController {
	return &AuditPlanBudgetController{
		service: service,
	}
}

func (apbc *AuditPlanBudgetController) GetAllAuditPlanBudgets(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))

	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"audit_plan_ids",
		"tool_type_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(apbc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(apbc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := dtos.AuditPlanBudgetListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditPlanBudgets, err := apbc.service.GetAll(uint(auditPlanID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(apbc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditPlanBudgets)
}

func (apbc *AuditPlanBudgetController) GetAuditPlanBudgetDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditPlanBudget, err := apbc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(apbc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditPlanBudget)
}

func (apbc *AuditPlanBudgetController) CreateAuditPlanBudget(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditPlanBudgetFormDTO dtos.AuditPlanBudgetFormDTO

	if err := c.ShouldBindJSON(&auditPlanBudgetFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apbc.db, c, err.Error())
		return
	}

	id, err := apbc.service.Create(uint(auditPlanID), currentUserID, &auditPlanBudgetFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apbc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (apbc *AuditPlanBudgetController) UpdateAuditPlanBudget(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditPlanBudgetFormDTO dtos.AuditPlanBudgetFormDTO

	auditPlanBudgetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditPlanBudgetFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apbc.db, c, err.Error())
		return
	}

	err = apbc.service.Update(uint(auditPlanBudgetID), uint(auditPlanID), currentUserID, &auditPlanBudgetFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apbc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (apbc *AuditPlanBudgetController) DeleteAuditPlanBudgets(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apbc.db, c, err.Error())
		return
	}

	failedIDs, err := apbc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(apbc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
