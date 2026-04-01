package controllers

import (
	"jethur/grc/consultation-request/dtos"
	"jethur/grc/consultation-request/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConsultationRequestController struct {
	service services.ConsultationRequestServiceInterface
	db      *gorm.DB
}

func NewConsultationRequestController(service services.ConsultationRequestServiceInterface) *ConsultationRequestController {
	return &ConsultationRequestController{
		service: service,
	}
}

// Get all consultation requests
func (crc *ConsultationRequestController) GetAllConsultationRequests(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	allowedParams := []string{
		"requesting_department_unit_ids",
		"consulted_department_unit_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Export flag
	exportFromContext, exists := c.Get("export")
	var export bool
	var err error
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(crc.db, c, "Invalid export parameter")
			return
		}
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page_size parameter")
		return
	}
	queryParams := dtos.ConsultationRequestListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
		Export:   export,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}
	requests, err := crc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	if export {
		excelFile, err := crc.service.ExportConsultationRequests(requests.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(crc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Disposition", `attachment; filename="Consultation Requests.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, requests)
}

// Get single consultation request by ID
func (crc *ConsultationRequestController) GetConsultationRequestDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	request, err := crc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, request)
}

// Create new consultation request
func (crc *ConsultationRequestController) CreateConsultationRequest(c *gin.Context) {
	var dto dtos.ConsultationRequestFormDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	createdRequest, err := crc.service.Create(currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, createdRequest)
}

// Update consultation request
func (crc *ConsultationRequestController) UpdateConsultationRequest(c *gin.Context) {
	var dto dtos.ConsultationRequestFormDTO
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.Update(uint(id), &dto, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

// Delete consultation requests
func (crc *ConsultationRequestController) DeleteConsultationRequests(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	failedIDs, err := crc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

// Activate consultation requests
func (crc *ConsultationRequestController) UpdateConsultationRequestResponsibleUser(c *gin.Context) {
	var dto struct {
		ResponsibleUserID uint `json:"responsible_user_id" binding:"required"`
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.UpdateResponsibleUser(uint(id), dto.ResponsibleUserID, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}
func (crc *ConsultationRequestController) UpdateConsultationRequestNotes(c *gin.Context) {
	var dto struct {
		Notes string `json:"notes" binding:"required"`
	}

	// Parse ID from URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Bind JSON payload
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	// Call service to update Notes
	err = crc.service.UpdateNotes(uint(id), dto.Notes, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (crc *ConsultationRequestController) CreateConsultationRequestComment(c *gin.Context) {
	var dto dtos.ConsultationRequestCommentFormDTO

	// Bind request body
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get consultation request ID from URL
	consultationRequestID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Get current user ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	// Call service
	createdComment, err := crc.service.CreateConsultationRequestComment(
		uint(consultationRequestID),
		currentUserID,
		&dto,
	)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, createdComment)
}

func (crc *ConsultationRequestController) UpdateConsultationRequestComment(c *gin.Context) {
	var dto dtos.ConsultationRequestCommentUpdateFormDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	consultationRequestID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentID, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.UpdateConsultationRequestComment(
		uint(consultationRequestID),
		uint(commentID),
		currentUserID,
		&dto,
	)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (crc *ConsultationRequestController) DeleteConsultationRequestComment(c *gin.Context) {
	consultationRequestID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	failedIDs, err := crc.service.DeleteConsultationRequestComments(uint(consultationRequestID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

// Export consultation requests
func (crc *ConsultationRequestController) GetConsultationRequestExport(c *gin.Context) {
	c.Set("export", true)

	crc.GetAllConsultationRequests(c)
}


func (crc *ConsultationRequestController) GetConsultationRequestCommentDocumentPreview(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	path, err := crc.service.GetDocumentPath(uint(requestID), uint(commentID), token)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	c.File(path)
}

func (crc *ConsultationRequestController) DownloadConsultationRequestCommentDocument(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	path, err := crc.service.GetDocumentPath(uint(requestID), uint(commentID), token)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	filename := filepath.Base(path)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(path)
}