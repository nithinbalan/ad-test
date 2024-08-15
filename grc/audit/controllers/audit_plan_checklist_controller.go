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

type AuditPlanChecklistController struct {
	service services.AuditPlanChecklistServiceInterface
	db      *gorm.DB
}

func NewAuditPlanChecklistController(service services.AuditPlanChecklistServiceInterface) *AuditPlanChecklistController {
	return &AuditPlanChecklistController{
		service: service,
	}
}

func (apc *AuditPlanChecklistController) GetAllAuditPlanChecklists(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
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
		"auditable_item_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := auditDTOs.AuditPlanChecklistListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditPlanChecklists, err := apc.service.GetAll(uint(auditPlanID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditPlanChecklists)
}

func (apc *AuditPlanChecklistController) CreateAuditPlanChecklists(c *gin.Context) {
	var auditPlanChecklistFormDTO auditDTOs.AuditPlanChecklistFormDTO

	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditPlanChecklistFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	id, err := apc.service.Create(uint(auditPlanID), currentUserID, &auditPlanChecklistFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (apc *AuditPlanChecklistController) UpdateAuditPlanChecklist(c *gin.Context) {
	var auditPlanChecklistFormDTO auditDTOs.AuditPlanChecklistFormDTO

	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditPlanChecklistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditPlanChecklistFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	err = apc.service.Update(uint(auditPlanID), uint(auditPlanChecklistID), currentUserID, &auditPlanChecklistFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (apc *AuditPlanChecklistController) DeleteAuditPlanChecklists(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	auditPlanID, err := strconv.Atoi(c.Param("planId"))
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
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	failedIDs, err := apc.service.Delete(uint(auditPlanID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
