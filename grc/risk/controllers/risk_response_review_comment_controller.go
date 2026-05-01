package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskResponseReviewCommentController struct {
	service services.RiskResponseReviewCommentServiceInterface
	db      *gorm.DB
}

func NewRiskResponseReviewCommentController(service services.RiskResponseReviewCommentServiceInterface) *RiskResponseReviewCommentController {
	return &RiskResponseReviewCommentController{
		service: service,
	}
}

func (rrcc *RiskResponseReviewCommentController) GetAllComments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "comment")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"comment_status_ids",
		"risk_response_ids",
		"level",
		"created_by_ids",
	}

	validSortFields := map[string]string{
		"comment":                 "risk_response_review_comments.comment",
		"level":                   "risk_response_review_comments.level",
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

	queryParams := riskDTOs.RiskResponseReviewCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	comments, err := rrcc.service.GetAll(uint(riskResponseID), currentUserID, &queryParams)
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

func (rrcc *RiskResponseReviewCommentController) GetCommentDetails(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
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

	comment, err := rrcc.service.GetFullDetailsByID(uint(riskResponseID), uint(commentID))
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

func (rrcc *RiskResponseReviewCommentController) CreateComment(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.RiskResponseReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	id, err := rrcc.service.Create(uint(riskResponseID), currentUserID, &formDTO)
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

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (rrcc *RiskResponseReviewCommentController) UpdateComment(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
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

	var formDTO riskDTOs.RiskResponseReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	err = rrcc.service.Update(uint(riskResponseID), uint(commentID), currentUserID, &formDTO)
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

	c.Status(http.StatusOK)
}

func (rrcc *RiskResponseReviewCommentController) DeleteComments(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
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

	failedIDs, err := rrcc.service.Delete(uint(riskResponseID), request.IDs, currentUserID)
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

	if len(failedIDs) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{"failed_ids": failedIDs})
		return
	}

	c.Status(http.StatusOK)
}

func (rrcc *RiskResponseReviewCommentController) ResolveComment(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
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

	err = rrcc.service.Resolve(uint(riskResponseID), uint(commentID), currentUserID)
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

	c.Status(http.StatusOK)
}

func (rrcc *RiskResponseReviewCommentController) GetDocumentPreview(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
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

	documentPath, err := rrcc.service.GetDocumentPath(uint(riskResponseID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrcc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (rrcc *RiskResponseReviewCommentController) DownloadDocument(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
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

	documentPath, err := rrcc.service.GetDocumentPath(uint(riskResponseID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrcc.db, c, err.Error())
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
