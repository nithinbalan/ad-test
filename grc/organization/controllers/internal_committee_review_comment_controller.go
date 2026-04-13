package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InternalCommitteeReviewCommentController struct {
	service services.InternalCommitteeReviewCommentServiceInterface
	db      *gorm.DB
}

func NewInternalCommitteeReviewCommentController(service services.InternalCommitteeReviewCommentServiceInterface) *InternalCommitteeReviewCommentController {
	return &InternalCommitteeReviewCommentController{
		service: service,
	}
}

func (rrcc *InternalCommitteeReviewCommentController) GetAllComments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "comment")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"comment_status_ids",
		"risk_register_ids",
		"level",
		"created_by_ids",
	}

	validSortFields := map[string]string{
		"comment":                 "internal_committee_review_comments.comment",
		"level":                   "internal_committee_review_comments.level",
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

	queryParams := organizationDTOs.InternalCommitteeReviewCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	comments, err := rrcc.service.GetAll(uint(committeeID), currentUserID, &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (rrcc *InternalCommitteeReviewCommentController) GetCommentDetails(c *gin.Context) {
	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
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

	comment, err := rrcc.service.GetFullDetailsByID(uint(committeeID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (rrcc *InternalCommitteeReviewCommentController) CreateComment(c *gin.Context) {
	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO organizationDTOs.InternalCommitteeReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	id, err := rrcc.service.Create(uint(committeeID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (rrcc *InternalCommitteeReviewCommentController) UpdateComment(c *gin.Context) {
	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
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

	var formDTO organizationDTOs.InternalCommitteeReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	err = rrcc.service.Update(uint(committeeID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrcc *InternalCommitteeReviewCommentController) DeleteComment(c *gin.Context) {
	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
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
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	failedIDs, err := rrcc.service.Delete(uint(committeeID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rrcc *InternalCommitteeReviewCommentController) ResolveComment(c *gin.Context) {

	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
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
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	err = rrcc.service.Resolve(uint(committeeID), uint(commentID), currentUserID)
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
			utility.SendInternalErrResp(rrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (iarrcc *InternalCommitteeReviewCommentController) GetDocumentPreview(c *gin.Context) {
	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentIDstr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDstr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := iarrcc.service.GetDocumentPath(uint(committeeID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(iarrcc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (iarrcc *InternalCommitteeReviewCommentController) DownloadDocument(c *gin.Context) {
	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	commentIDstr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDstr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := iarrcc.service.GetDocumentPath(uint(committeeID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(iarrcc.db, c, err.Error())
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
