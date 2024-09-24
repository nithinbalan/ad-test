package controllers

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"path/filepath"
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

func (adrcc *AuditDocumentRequestConversationController) DownloadDocumentRequestConversationDocument(c *gin.Context) {
	auditIdStr := c.Param("auditId")
	auditId, err := strconv.ParseUint(auditIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	auditDocumentRequestID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	conversationID, err := strconv.Atoi(c.Param("conversationId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := adrcc.service.GetDocumentPath(uint(auditId), uint(auditDocumentRequestID), uint(conversationID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(adrcc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(documentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(documentPath)
}
