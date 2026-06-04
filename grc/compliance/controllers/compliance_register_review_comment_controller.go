package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceRegisterReviewCommentController struct {
	service services.ComplianceRegisterReviewCommentServiceInterface
	db      *gorm.DB
}

func NewComplianceRegisterReviewCommentController(service services.ComplianceRegisterReviewCommentServiceInterface) *ComplianceRegisterReviewCommentController {
	return &ComplianceRegisterReviewCommentController{
		service: service,
	}
}

func (crcc *ComplianceRegisterReviewCommentController) GetAllComments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "comment")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"comment_status_ids",
		"compliance_register_ids",
		"level",
		"created_by_ids",
	}

	validSortFields := map[string]string{
		"comment":                 "compliance_register_review_comments.comment",
		"level":                   "compliance_register_review_comments.level",
		"comment_status_title_en": "comment_statuses.title_ar",
		"comment_status_title_ar": "comment_statuses.title_en",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := complianceDTOs.ComplianceRegisterReviewCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	comments, err := crcc.service.GetAll(uint(complianceRegisterID), currentUserID, &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (crcc *ComplianceRegisterReviewCommentController) GetCommentDetails(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comment, err := crcc.service.GetFullDetailsByID(uint(complianceRegisterID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (crcc *ComplianceRegisterReviewCommentController) CreateComment(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO complianceDTOs.ComplianceRegisterReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	id, err := crcc.service.Create(uint(complianceRegisterID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (crcc *ComplianceRegisterReviewCommentController) UpdateComment(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO complianceDTOs.ComplianceRegisterReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	err = crcc.service.Update(uint(complianceRegisterID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (crcc *ComplianceRegisterReviewCommentController) DeleteComments(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
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
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	failedIDs, err := crcc.service.Delete(uint(complianceRegisterID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crcc.db, c, e.Error())
			return
		}
	}

	if len(failedIDs) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{"failed_ids": failedIDs})
		return
	}

	c.Status(http.StatusOK)
}

func (crcc *ComplianceRegisterReviewCommentController) ResolveComment(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	err = crcc.service.Resolve(uint(complianceRegisterID), uint(commentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (crcc *ComplianceRegisterReviewCommentController) GetDocumentPreview(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := crcc.service.GetDocumentPath(uint(complianceRegisterID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(crcc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (crcc *ComplianceRegisterReviewCommentController) DownloadDocument(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := crcc.service.GetDocumentPath(uint(complianceRegisterID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(crcc.db, c, err.Error())
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
