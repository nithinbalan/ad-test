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

type AuditDraftFindingController struct {
	service services.AuditDraftFindingServiceInterface
	db      *gorm.DB
}

func NewAuditDraftFindingController(service services.AuditDraftFindingServiceInterface) *AuditDraftFindingController {
	return &AuditDraftFindingController{
		service: service,
	}
}

func (adfc *AuditDraftFindingController) GetAllAuditDraftFindings(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_ids",
		"auditable_item_ids",
		"audit_finding_category_ids",
		"severity_level_ids",
		"audit_finding_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(adfc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(adfc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditDTOs.AuditDraftFindingListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditDraftFindings, err := adfc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(adfc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditDraftFindings)
}

func (adfc *AuditDraftFindingController) GetAuditDraftFindingDetails(c *gin.Context) {

	idStr := c.Param("findingId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditDraftFinding, err := adfc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(adfc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditDraftFinding)
}

func (adfc *AuditDraftFindingController) CreateAuditDraftFinding(c *gin.Context) {
	var auditDraftFindingFormDTO auditDTOs.AuditDraftFindingFormDTO

	if err := c.ShouldBindJSON(&auditDraftFindingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adfc.db, c, err.Error())
		return
	}

	id, err := adfc.service.Create(currentUserID, &auditDraftFindingFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adfc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (adfc *AuditDraftFindingController) UpdateAuditDraftFinding(c *gin.Context) {
	var auditDraftFindingFormDTO auditDTOs.AuditDraftFindingFormDTO

	id, err := strconv.Atoi(c.Param("findingId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditDraftFindingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adfc.db, c, err.Error())
		return
	}

	err = adfc.service.Update(uint(id), currentUserID, &auditDraftFindingFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adfc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (adfc *AuditDraftFindingController) DeleteAuditDraftFindings(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adfc.db, c, err.Error())
		return
	}

	failedIDs, err := adfc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(adfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
