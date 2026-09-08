package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectChangeRequestReviewCommentController struct {
	service services.ProjectChangeRequestReviewCommentServiceInterface
	db      *gorm.DB
}

func NewProjectChangeRequestReviewCommentController(service services.ProjectChangeRequestReviewCommentServiceInterface) *ProjectChangeRequestReviewCommentController {
	return &ProjectChangeRequestReviewCommentController{
		service: service,
	}
}

func (pcrrcc *ProjectChangeRequestReviewCommentController) GetAllComments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "comment")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"comment_status_ids",
		"project_change_request_ids",
		"level",
		"created_by_ids",
	}

	validSortFields := map[string]string{
		"comment":                 "project_change_request_review_comments.comment",
		"level":                   "project_change_request_review_comments.level",
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

	queryParams := projectDTOs.ProjectChangeRequestReviewCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcrrcc.db, c, err.Error())
		return
	}

	comments, err := pcrrcc.service.GetAll(uint(projectID), uint(changeRequestID), currentUserID, &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pcrrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (pcrrcc *ProjectChangeRequestReviewCommentController) GetCommentDetails(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
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

	comment, err := pcrrcc.service.GetFullDetailsByID(uint(projectID), uint(changeRequestID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pcrrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (pcrrcc *ProjectChangeRequestReviewCommentController) CreateComment(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO projectDTOs.ProjectChangeRequestReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcrrcc.db, c, err.Error())
		return
	}

	id, err := pcrrcc.service.Create(uint(projectID), uint(changeRequestID), currentUserID, &formDTO)
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
			utility.SendInternalErrResp(pcrrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (pcrrcc *ProjectChangeRequestReviewCommentController) UpdateComment(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
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

	var formDTO projectDTOs.ProjectChangeRequestReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcrrcc.db, c, err.Error())
		return
	}

	err = pcrrcc.service.Update(uint(projectID), uint(changeRequestID), uint(commentID), currentUserID, &formDTO)
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
			utility.SendInternalErrResp(pcrrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pcrrcc *ProjectChangeRequestReviewCommentController) DeleteComments(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
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
		utility.SendInternalErrResp(pcrrcc.db, c, err.Error())
		return
	}

	failedIDs, err := pcrrcc.service.Delete(uint(projectID), uint(changeRequestID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pcrrcc.db, c, e.Error())
			return
		}
	}

	if len(failedIDs) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{"failed_ids": failedIDs})
		return
	}

	c.Status(http.StatusOK)
}

func (pcrrcc *ProjectChangeRequestReviewCommentController) ResolveComment(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
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
		utility.SendInternalErrResp(pcrrcc.db, c, err.Error())
		return
	}

	err = pcrrcc.service.Resolve(uint(projectID), uint(changeRequestID), uint(commentID), currentUserID)
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
			utility.SendInternalErrResp(pcrrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pcrrcc *ProjectChangeRequestReviewCommentController) GetDocumentPreview(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
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

	documentPath, err := pcrrcc.service.GetDocumentPath(uint(projectID), uint(changeRequestID), uint(commentID), token)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(pcrrcc.db, c, e.Error())
		}
		return
	}
	c.File(documentPath)
}

func (pcrrcc *ProjectChangeRequestReviewCommentController) DownloadDocument(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
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

	documentPath, err := pcrrcc.service.GetDocumentPath(uint(projectID), uint(changeRequestID), uint(commentID), token)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(pcrrcc.db, c, e.Error())
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
