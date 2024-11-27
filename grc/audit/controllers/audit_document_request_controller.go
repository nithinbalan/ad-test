package controllers

import (
	"errors"
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditDocumentRequestController struct {
	service services.AuditDocumentRequestServiceInterface
	db      *gorm.DB
}

func NewAuditDocumentRequestController(service services.AuditDocumentRequestServiceInterface) *AuditDocumentRequestController {
	return &AuditDocumentRequestController{
		service: service,
	}
}

func (adrc *AuditDocumentRequestController) GetAllAuditDocumentRequests(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	requestedAt := c.Query("requested_at")

	allowedParams := []string{
		"audit_ids",
		"responsible_user_ids",
		"responsible_unit_ids",
		"requested_by_ids",
		"audit_document_request_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"title":                                  "audit_document_requests.title",
		"requested_at":                           "audit_document_requests.requested_at",
		"responsible_unit_en":                    "organization_units.title_en",
		"responsible_unit_ar":                    "organization_units.title_ar",
		"responsible_user_first_name_en":         "responsible_user.first_name_en",
		"responsible_user_first_name_ar":         "responsible_user.first_name_ar",
		"responsible_user_last_name_en":          "responsible_user.last_name_en",
		"responsible_user_last_name_ar":          "responsible_user.last_name_ar",
		"requested_user_first_name_en":           "requested_user.first_name_en",
		"requested_user_first_name_ar":           "requested_user.first_name_ar",
		"requested_user_last_name_en":            "requested_user.last_name_en",
		"requested_user_last_name_ar":            "requested_user.last_name_ar",
		"audit_document_request_status_title_en": "audit_document_request_statuses.title_en",
		"audit_document_request_status_title_ar": "audit_document_request_statuses.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := dtos.AuditDocumentRequestListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		RequestedAt: requestedAt,
	}

	auditDocumentRequests, err := adrc.service.GetAll(uint(auditID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditDocumentRequests)
}

func (adrc *AuditDocumentRequestController) GetAuditDocumentRequestDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditDocumentRequest, err := adrc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(adrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditDocumentRequest)
}

func (adrc *AuditDocumentRequestController) CreateAuditDocumentRequest(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var auditDocumentRequestFormDTO dtos.AuditDocumentRequestFormDTO

	if err := c.ShouldBindJSON(&auditDocumentRequestFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}

	ids, err := adrc.service.Create(uint(auditID), currentUserID, &auditDocumentRequestFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"ids": ids})
}

func (adrc *AuditDocumentRequestController) UpdateAuditDocumentRequest(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditDocumentRequestID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var auditDocumentRequestFormDTO dtos.AuditDocumentRequestFormDTO

	if err := c.ShouldBindJSON(&auditDocumentRequestFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}

	err = adrc.service.Update(uint(auditID), uint(auditDocumentRequestID), currentUserID, &auditDocumentRequestFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (adrc *AuditDocumentRequestController) DeleteAuditDocumentRequests(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}

	failedIDs, err := adrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (adrc *AuditDocumentRequestController) SubmitAuditDocumentRequest(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditDocumentRequestID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrc.db, c, err.Error())
		return
	}

	err = adrc.service.Submit(uint(auditID), uint(auditDocumentRequestID), currentUserID)
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
			utility.SendInternalErrResp(adrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
