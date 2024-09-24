package controllers

import (
	"jethur/grc/audit-insights/services"
	"jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightsDocumentRequestConversationController struct {
	service services.AuditInsightsDocumentRequestConversationServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsDocumentRequestConversationController(service services.AuditInsightsDocumentRequestConversationServiceInterface) *AuditInsightsDocumentRequestConversationController {
	return &AuditInsightsDocumentRequestConversationController{
		service: service,
	}
}

func (aidrcc *AuditInsightsDocumentRequestConversationController) GetAllAuditInsightsDocumentRequestConversations(c *gin.Context) {
	auditInsightsDocumentRequestID, err := strconv.Atoi(c.Param("documentRequestId"))
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

	auditInsightsDocumentRequestConversations, err := aidrcc.service.GetAll(uint(auditInsightsDocumentRequestID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(aidrcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditInsightsDocumentRequestConversations)
}

func (aidrcc *AuditInsightsDocumentRequestConversationController) CreateAuditInsightsDocumentRequestConversation(c *gin.Context) {
	auditInsightsDocumentRequestID, err := strconv.Atoi(c.Param("documentRequestId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditInsightsDocumentRequestConversationFormDTO dtos.AuditDocumentRequestConversationFormDTO

	if err := c.ShouldBindJSON(&auditInsightsDocumentRequestConversationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidrcc.db, c, err.Error())
		return
	}

	id, err := aidrcc.service.Create(uint(auditInsightsDocumentRequestID), currentUserID, &auditInsightsDocumentRequestConversationFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aidrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}
