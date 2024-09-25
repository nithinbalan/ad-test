package controllers

import (
	"errors"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentController struct {
	service services.DocumentServiceInterface
	db      *gorm.DB
}

func NewDocumentController(service services.DocumentServiceInterface) *DocumentController {
	return &DocumentController{
		service: service,
	}
}

func (dc *DocumentController) GetAllDocuments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	enableAcknowledgedStatus := c.Query("enable_acknowledge_status")
	documentTypes := utility.ParseCommaSeparatedStringToArray(c.Query("document_types"))
	documentStatus := utility.ParseCommaSeparatedStringToArray(c.Query("document_status"))

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"document_category_ids",
		"document_type_ids",
		"document_access_type_ids",
		"responsible_user_ids",
		"responsible_unit_ids",
		"document_status_ids",
	}

	validSortFields := map[string]string{
		"type_ar":             "document_types.title_ar",
		"type_en":             "document_types.title_en",
		"title_ar":            "documents.title_ar",
		"title_en":            "documents.title_en",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"category_en":         "document_categories.title_en",
		"category_ar":         "document_categories.title_ar",
		"status_en":           "document_statuses.title_ar",
		"status_ar":           "document_statuses.title_en",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := documentDTOs.DocumentListQueryParamsDTO{
		Sort:                     sort,
		SortBy:                   dbField,
		Keyword:                  keyword,
		Status:                   status,
		Page:                     page,
		PageSize:                 pageSize,
		Filters:                  filters,
		QuickFilter:              quickFilter,
		EnableAcknowledgedStatus: enableAcknowledgedStatus == "true",
		DocumentTypes:            documentTypes,
		DocumentStatus:           documentStatus,
		Export:                   export,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	documents, err := dc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(dc.db, c, err.Error())
			return
		}

		excelFile, err := dc.service.ExportDocuments(documents.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(dc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Document.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, documents)
}

func (dc *DocumentController) GetDocumentDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	document, err := dc.service.GetFullDetailsByID(uint(id), &currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}
	//Add audit log for the document
	if err := dc.service.DirectDocumentAuditLog(uint(id), currentUserID, "viewed"); err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, document)
}

func (dc *DocumentController) CreateDocument(c *gin.Context) {
	var documentFormDTO documentDTOs.DocumentFormDTO
	if err := c.ShouldBindJSON(&documentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	createdDocument, err := dc.service.Create(currentUserID, &documentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdDocument)
}

func (dc *DocumentController) DeleteDocuments(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	failedIDs, err := dc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (dc *DocumentController) CheckoutDocument(c *gin.Context) {
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
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	documentPath, err := dc.service.CheckoutDocument(uint(id), currentUserID, token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
		default:
			utility.SendInternalErrResp(dc.db, c, err.Error())
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

func (dc *DocumentController) CheckinDocument(c *gin.Context) {
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
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	err = dc.service.CheckinDocument(uint(id), currentUserID, &documentCheckinFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (dc *DocumentController) RevertDocumentCheckout(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	err = dc.service.RevertDocumentCheckout(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (dc *DocumentController) GetDocumentPreview(c *gin.Context) {
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

	documentPath, err := dc.service.GetDocumentPath(uint(id), token, fileType)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(dc.db, c, err.Error())
		}
		return
	}
	// currentUserID, err := utility.GetCurrentUserID(c)
	// if err != nil {
	// 	utility.SendInternalErrResp(dc.db, c, err.Error())
	// 	return
	// }
	// // Add audit log for the document
	// if err := dc.service.DirectDocumentAuditLog(uint(id), uint(currentUserID), "previewed"); err != nil {
	// 	utility.SendInternalErrResp(dc.db, c, err.Error())
	// 	return
	// }
	c.File(documentPath)
}

func (dc *DocumentController) SubmitDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	err = dc.service.SubmitDocument(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (dc *DocumentController) PublishDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	err = dc.service.PublishDocument(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (dc *DocumentController) GetComments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := dc.service.GetComments(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (dc *DocumentController) AddComment(c *gin.Context) {
	var documentReviewCommentFormDTO documentDTOs.DocumentReviewCommentFormDTO
	if err := c.ShouldBindJSON(&documentReviewCommentFormDTO); err != nil {
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
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	err = dc.service.AddComment(uint(id), currentUserID, &documentReviewCommentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (dc *DocumentController) GetWorkflowList(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := dc.service.GetDocumentWorkflowList(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}

func (dc *DocumentController) DownloadDocument(c *gin.Context) {
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

	documentPath, err := dc.service.GetDocumentPath(uint(id), token, fileType)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(dc.db, c, err.Error())
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

func (dc *DocumentController) GetSharedUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, "Invalid page_size parameter")
		return
	}

	sharedUsers, err := dc.service.GetSharedUsers(uint(id), uint(page), uint(pageSize))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, sharedUsers)
}

func (dc *DocumentController) GetReviewUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reviewUsers, err := dc.service.GetReviewUsers(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, reviewUsers)
}

func (dc *DocumentController) GetSharedUnits(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sharedUnits, err := dc.service.GetSharedUnits(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, sharedUnits)
}

func (dc *DocumentController) GetLastEditSummary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	lastEditSummary, err := dc.service.GetDocumentLastEditSummary(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, lastEditSummary)
}

func (dc *DocumentController) GetDocumentExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	dc.GetAllDocuments(c)
}
