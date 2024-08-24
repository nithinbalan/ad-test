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

type AuditController struct {
	service services.AuditServiceInterface
	db      *gorm.DB
}

func NewAuditController(service services.AuditServiceInterface) *AuditController {
	return &AuditController{
		service: service,
	}
}

func (apc *AuditController) GetAllAudits(c *gin.Context) {
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

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditDTOs.AuditListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	audits, err := apc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, audits)
}

func (apc *AuditController) GetAuditDetails(c *gin.Context) {
	idStr := c.Param("auditId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	audit, err := apc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(apc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, audit)
}

func (apc *AuditController) CreateAudit(c *gin.Context) {
	var auditFormDTO auditDTOs.AuditFormDTO

	if err := c.ShouldBindJSON(&auditFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	id, err := apc.service.Create(currentUserID, &auditFormDTO)
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

func (apc *AuditController) UpdateAudit(c *gin.Context) {
	var auditFormDTO auditDTOs.AuditFormDTO

	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	err = apc.service.Update(uint(auditID), currentUserID, &auditFormDTO)
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

func (apc *AuditController) DeleteAudits(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
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

	failedIDs, err := apc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
