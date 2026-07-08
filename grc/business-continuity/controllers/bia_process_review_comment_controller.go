package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BiaProcessReviewCommentController struct {
	service services.BiaProcessReviewCommentServiceInterface
	db      *gorm.DB
}

func NewBiaProcessReviewCommentController(service services.BiaProcessReviewCommentServiceInterface) *BiaProcessReviewCommentController {
	return &BiaProcessReviewCommentController{
		service: service,
	}
}

func (bprcc *BiaProcessReviewCommentController) GetAllComments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "comment")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"comment_status_ids",
		"bia_process_ids",
		"level",
		"created_by_ids",
	}

	validSortFields := map[string]string{
		"comment":                 "bia_process_review_comments.comment",
		"level":                   "bia_process_review_comments.level",
		"comment_status_title_en": "comment_statuses.title_en",
		"comment_status_title_ar": "comment_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := businessContinuityDTOs.BiaProcessReviewCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bprcc.db, c, err.Error())
		return
	}

	comments, err := bprcc.service.GetAll(uint(biaProcessID), currentUserID, &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bprcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (bprcc *BiaProcessReviewCommentController) GetCommentDetails(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
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

	comment, err := bprcc.service.GetFullDetailsByID(uint(biaProcessID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bprcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (bprcc *BiaProcessReviewCommentController) CreateComment(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO businessContinuityDTOs.BiaProcessReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bprcc.db, c, err.Error())
		return
	}

	id, err := bprcc.service.Create(uint(biaProcessID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bprcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (bprcc *BiaProcessReviewCommentController) UpdateComment(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
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

	var formDTO businessContinuityDTOs.BiaProcessReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bprcc.db, c, err.Error())
		return
	}

	err = bprcc.service.Update(uint(biaProcessID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bprcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (bprcc *BiaProcessReviewCommentController) DeleteComments(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
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
		utility.SendInternalErrResp(bprcc.db, c, err.Error())
		return
	}

	failedIDs, err := bprcc.service.Delete(uint(biaProcessID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(bprcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (bprcc *BiaProcessReviewCommentController) ResolveComment(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
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
		utility.SendInternalErrResp(bprcc.db, c, err.Error())
		return
	}

	err = bprcc.service.Resolve(uint(biaProcessID), uint(commentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(bprcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (bprcc *BiaProcessReviewCommentController) GetDocumentPreview(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
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

	filePath, err := bprcc.service.GetDocumentPath(uint(biaProcessID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(bprcc.db, c, err.Error())
		}
		return
	}

	c.File(filePath)
}

func (bprcc *BiaProcessReviewCommentController) DownloadDocument(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
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

	filePath, err := bprcc.service.GetDocumentPath(uint(biaProcessID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(bprcc.db, c, err.Error())
		}
		return
	}

	filename := filepath.Base(filePath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(filePath)
}
