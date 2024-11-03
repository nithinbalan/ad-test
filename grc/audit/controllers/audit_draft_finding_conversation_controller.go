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

type AuditDraftFindingConversationController struct {
	service services.AuditDraftFindingConversationServiceInterface
	db      *gorm.DB
}

func NewAuditDraftFindingConversationController(service services.AuditDraftFindingConversationServiceInterface) *AuditDraftFindingConversationController {
	return &AuditDraftFindingConversationController{
		service: service,
	}
}

func (adfcc *AuditDraftFindingConversationController) GetAllAuditDraftFindingConversations(c *gin.Context) {
	auditDraftFindingID, err := strconv.Atoi(c.Param("findingId"))
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

	queryParams := dtos.AuditDraftFindingConversationListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditDraftFindingConversations, err := adfcc.service.GetAll(uint(auditDraftFindingID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(adfcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditDraftFindingConversations)
}
func (adfcc *AuditDraftFindingConversationController) GetAuditDraftFindingConversationDetails(c *gin.Context) {
	auditDraftFindingID, err := strconv.Atoi(c.Param("findingId"))
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

	auditDraftFinding, err := adfcc.service.GetFullDetailsByID(uint(auditDraftFindingID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(adfcc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditDraftFinding)
}

func (adfcc *AuditDraftFindingConversationController) CreateAuditDraftFindingConversation(c *gin.Context) {
	auditDraftFindingID, err := strconv.Atoi(c.Param("findingId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditDraftFindingConversationFormDTO dtos.AuditDraftFindingConversationFormDTO

	if err := c.ShouldBindJSON(&auditDraftFindingConversationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adfcc.db, c, err.Error())
		return
	}

	id, err := adfcc.service.Create(uint(auditDraftFindingID), currentUserID, &auditDraftFindingConversationFormDTO)
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
			utility.SendInternalErrResp(adfcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (adfcc *AuditDraftFindingConversationController) UpdateAuditDraftFindingConversation(c *gin.Context) {
	var AuditDraftFindingConversationFormDTO dtos.AuditDraftFindingConversationFormDTO

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
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

	if err := c.ShouldBindJSON(&AuditDraftFindingConversationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adfcc.db, c, err.Error())
		return
	}

	err = adfcc.service.Update(uint(findingID), uint(id), currentUserID, &AuditDraftFindingConversationFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adfcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (adfcc *AuditDraftFindingConversationController) DeleteAuditDraftFindingConversations(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
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
		utility.SendInternalErrResp(adfcc.db, c, err.Error())
		return
	}

	failedIDs, err := adfcc.service.Delete(uint(findingID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adfcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (adfcc *AuditDraftFindingConversationController) GetAuditDraftFindingConversationDocumentPreview(c *gin.Context) {

	findingIdStr := c.Param("findingId")
	findingId, err := strconv.ParseUint(findingIdStr, 10, 32)
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

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := adfcc.service.GetDocumentPath(uint(findingId), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(adfcc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (adfcc *AuditDraftFindingConversationController) DownloadAuditDraftFindingConversationDocument(c *gin.Context) {
	findingIdStr := c.Param("findingId")
	findingId, err := strconv.ParseUint(findingIdStr, 10, 32)
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
	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := adfcc.service.GetDocumentPath(uint(findingId), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(adfcc.db, c, err.Error())
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
