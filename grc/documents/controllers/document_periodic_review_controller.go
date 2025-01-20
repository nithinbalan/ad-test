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

type DocumentPeriodicReviewController struct {
	service services.DocumentPeriodicReviewServiceInterface
	db      *gorm.DB
}

func NewDocumentPeriodicReviewController(service services.DocumentPeriodicReviewServiceInterface) *DocumentPeriodicReviewController {
	return &DocumentPeriodicReviewController{
		service: service,
	}
}

// Get all document periodic reviews with filtering, sorting, and pagination
func (dprc *DocumentPeriodicReviewController) GetAllDocumentPeriodicReviews(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "review_date")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"review_status_ids",
	}

	validSortFields := map[string]string{
		"review_date":        "document_reviews.review_date",
		"is_change_required": "document_reviews.is_change_required",
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

	params := documentDTOs.DocumentPeriodicReviewQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	documentIdStr := c.Param("id")
	documentId, err := strconv.ParseUint(documentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentReviews, err := dprc.service.GetAll(&params, uint(documentId))
	if err != nil {
		utility.SendInternalErrResp(dprc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, documentReviews)
}

// Get details of a single document periodic review by ID
func (dprc *DocumentPeriodicReviewController) GetDocumentPeriodicReviewDetails(c *gin.Context) {
	idStr := c.Param("reviewId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := dprc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, details)
}

// Create a new document periodic review
func (dprc *DocumentPeriodicReviewController) CreateDocumentPeriodicReview(c *gin.Context) {
	var formDTO documentDTOs.DocumentPeriodicReviewFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	documentIdStr := c.Param("id")
	documentId, err := strconv.ParseUint(documentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dprc.db, c, err.Error())
		return
	}

	id, err := dprc.service.Create(currentUserID, uint(documentId), &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

// Update an existing document periodic review
func (dprc *DocumentPeriodicReviewController) UpdateDocumentPeriodicReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("reviewId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO documentDTOs.DocumentPeriodicReviewFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dprc.db, c, err.Error())
		return
	}

	err = dprc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

// Delete document periodic reviews by IDs
func (dprc *DocumentPeriodicReviewController) DeleteDocumentPeriodicReview(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dprc.db, c, err.Error())
		return
	}

	failedIDs, err := dprc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dprc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

// Download a document periodic review file
func (dprc *DocumentPeriodicReviewController) DownloadDocumentPeriodicReviewDocument(c *gin.Context) {
	idStr := c.Param("reviewId")
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := dprc.service.GetDocumentPath(uint(reviewID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(dprc.db, c, err.Error())
		}
		return
	}

	filename := filepath.Base(documentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}
