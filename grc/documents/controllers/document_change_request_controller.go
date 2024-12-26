package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentChangeRequestController struct {
	service services.DocumentChangeRequestServiceInterface
	db      *gorm.DB
}

func NewDocumentChangeRequestController(service services.DocumentChangeRequestServiceInterface) *DocumentChangeRequestController {
	return &DocumentChangeRequestController{
		service: service,
	}
}

// Get all document change requests with filtering, sorting, and pagination
func (dcrc *DocumentChangeRequestController) GetAllDocumentChangeRequests(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "reference_number")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"document_change_request_type_ids",
		"document_change_request_status_ids",
		"document_ids",
	}

	validSortFields := map[string]string{
		"document_change_reuqest_type_ar": "change_request_types.title_ar",
		"document_change_reuqest_type_en": "change_request_types.title_en",
		"reason_for_change":               "document_change_requests.reason_for_change",
		"reference_number":                "document_change_requests.reason_for_change",
		"status_en":                       "change_request_statuses.title_ar",
		"status_ar":                       "change_request_statuses.title_en",
		"document_ar":                     "documents.title_en",
		"document_en":                     "documents.title_en",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page parameter"))
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page_size parameter"))
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	params := documentDTOs.DocumentChangeRequestQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	documentChangeRequests, err := dcrc.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, documentChangeRequests)
}

// Get details of a single document change request by ID
func (dcrc *DocumentChangeRequestController) GetDocumentChangeRequestDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := dcrc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, details)
}

// Create a new document change request
func (dcrc *DocumentChangeRequestController) CreateDocumentChangeRequest(c *gin.Context) {
	var formDTO documentDTOs.DocumentChangeRequestFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	id, err := dcrc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

// Update an existing document change request
func (dcrc *DocumentChangeRequestController) UpdateDocumentChangeRequest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO documentDTOs.DocumentChangeRequestFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	err = dcrc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

// Delete document change requests by IDs
func (dcrc *DocumentChangeRequestController) DeleteDocumentChangeRequest(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	failedIDs, err := dcrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (dcrc *DocumentChangeRequestController) GetDocumentChangeRequestPreview(c *gin.Context) {
	fileType := "pdf"
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

	documentPath, err := dcrc.service.GetDocumentPath(uint(id), token, fileType)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(dcrc.db, c, err.Error())
		}
		return
	}

	c.File(documentPath)
}

func (dcrc *DocumentChangeRequestController) DownloadChangeRequestDocument(c *gin.Context) {
	fileType := "docx"
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

	documentPath, err := dcrc.service.GetDocumentPath(uint(id), token, fileType)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(dcrc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(documentPath)
	// currentUserID, err := utility.GetCurrentUserID(c)
	// if err != nil {
	// 	utility.SendInternalErrResp(dc.db, c, err.Error())
	// 	return
	// }
	// // Add audit log for the document
	// if err := dc.service.DirectDocumentAuditLog(uint(id), uint(currentUserID), "downloaded"); err != nil {
	// 	utility.SendInternalErrResp(dc.db, c, err.Error())
	// 	return
	// }

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

func (dcrc *DocumentChangeRequestController) CheckoutDocument(c *gin.Context) {
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

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	documentPath, err := dcrc.service.CheckoutDocument(uint(id), currentUserID, token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
		default:
			utility.SendInternalErrResp(dcrc.db, c, err.Error())
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

func (dcrc *DocumentChangeRequestController) CheckinDocument(c *gin.Context) {
	var documentCheckinFormDTO documentDTOs.DocumentCheckinFormDTO
	if err := c.ShouldBindJSON(&documentCheckinFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	err = dcrc.service.CheckinDocument(uint(id), currentUserID, &documentCheckinFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (dcrc *DocumentChangeRequestController) RevertDocumentCheckout(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	err = dcrc.service.RevertDocumentCheckout(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (dcrc *DocumentChangeRequestController) SubmitDocumentChangeRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	err = dcrc.service.SubmitDocumentChangeRequest(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (dcrc *DocumentChangeRequestController) GetLastEditSummary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	lastEditSummary, err := dcrc.service.GetDocumentCRLastEditSummary(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, lastEditSummary)
}
func (dcrc *DocumentChangeRequestController) PublishDocumentChangeRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	err = dcrc.service.PublishDocumentChangeRequest(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
