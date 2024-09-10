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

type AuditDocumentRequestConversationController struct {
	service services.AuditDocumentRequestConversationServiceInterface
	db      *gorm.DB
}

func NewAuditDocumentRequestConversationController(service services.AuditDocumentRequestConversationServiceInterface) *AuditDocumentRequestConversationController {
	return &AuditDocumentRequestConversationController{
		service: service,
	}
}

func (adrcc *AuditDocumentRequestConversationController) GetAllAuditDocumentRequestConversations(c *gin.Context) {
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

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"audit_document_request_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.AuditDocumentRequestConversationListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditDocumentRequestConversations, err := adrcc.service.GetAll(uint(auditID), uint(auditDocumentRequestID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(adrcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditDocumentRequestConversations)
}

func (adrcc *AuditDocumentRequestConversationController) CreateAuditDocumentRequestConversation(c *gin.Context) {
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
	var auditDocumentRequestConversationFormDTO dtos.AuditDocumentRequestConversationFormDTO

	if err := c.ShouldBindJSON(&auditDocumentRequestConversationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrcc.db, c, err.Error())
		return
	}

	id, err := adrcc.service.Create(uint(auditID), uint(auditDocumentRequestID), currentUserID, &auditDocumentRequestConversationFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}
