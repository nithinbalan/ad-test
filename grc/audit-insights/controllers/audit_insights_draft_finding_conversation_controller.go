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

type AuditInsightsDraftFindingConversationController struct {
	service services.AuditInsightsDraftFindingConversationServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsDraftFindingConversationController(service services.AuditInsightsDraftFindingConversationServiceInterface) *AuditInsightsDraftFindingConversationController {
	return &AuditInsightsDraftFindingConversationController{
		service: service,
	}
}

func (aidrcc *AuditInsightsDraftFindingConversationController) GetAllAuditInsightsDraftFindingConversations(c *gin.Context) {
	findingID, err := strconv.Atoi(c.Param("findingId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"audit_draft_finding_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.AuditDocumentRequestConversationListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditInsightsDraftFindingConversations, err := aidrcc.service.GetAll(uint(findingID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(aidrcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditInsightsDraftFindingConversations)
}

func (aidrcc *AuditInsightsDraftFindingConversationController) CreateAuditInsightsDraftFindingConversation(c *gin.Context) {
	findingId, err := strconv.Atoi(c.Param("findingId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditInsightsDraftFindingConversationFormDTO dtos.AuditDocumentRequestConversationFormDTO

	if err := c.ShouldBindJSON(&auditInsightsDraftFindingConversationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidrcc.db, c, err.Error())
		return
	}

	id, err := aidrcc.service.Create(uint(findingId), currentUserID, &auditInsightsDraftFindingConversationFormDTO)
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

