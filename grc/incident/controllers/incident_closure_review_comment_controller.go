package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentClosureReviewCommentController struct {
	service services.IncidentClosureReviewCommentServiceInterface
	db      *gorm.DB
}

func NewIncidentClosureReviewCommentController(service services.IncidentClosureReviewCommentServiceInterface) *IncidentClosureReviewCommentController {
	return &IncidentClosureReviewCommentController{
		service: service,
	}
}

func (icrcc *IncidentClosureReviewCommentController) GetAllComments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "comment")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"comment_status_ids",
		"incident_closure_ids",
		"level",
		"created_by_ids",
	}

	validSortFields := map[string]string{
		"comment":                 "incident_closure_review_comments.comment",
		"level":                   "incident_closure_review_comments.level",
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

	queryParams := incidentDTOs.IncidentClosureReviewCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	incidentClosureIDStr := c.Param("closureId")
	incidentClosureID, err := strconv.ParseUint(incidentClosureIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icrcc.db, c, err.Error())
		return
	}

	comments, err := icrcc.service.GetAll(uint(incidentClosureID), currentUserID, &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(icrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (icrcc *IncidentClosureReviewCommentController) GetCommentDetails(c *gin.Context) {
	incidentClosureIDStr := c.Param("closureId")
	incidentClosureID, err := strconv.ParseUint(incidentClosureIDStr, 10, 32)
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

	comment, err := icrcc.service.GetFullDetailsByID(uint(incidentClosureID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(icrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (icrcc *IncidentClosureReviewCommentController) CreateComment(c *gin.Context) {
	incidentClosureIDStr := c.Param("closureId")
	incidentClosureID, err := strconv.ParseUint(incidentClosureIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO incidentDTOs.IncidentClosureReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icrcc.db, c, err.Error())
		return
	}

	id, err := icrcc.service.Create(uint(incidentClosureID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (icrcc *IncidentClosureReviewCommentController) UpdateComment(c *gin.Context) {
	incidentClosureIDStr := c.Param("closureId")
	incidentClosureID, err := strconv.ParseUint(incidentClosureIDStr, 10, 32)
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

	var formDTO incidentDTOs.IncidentClosureReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icrcc.db, c, err.Error())
		return
	}

	err = icrcc.service.Update(uint(incidentClosureID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (icrcc *IncidentClosureReviewCommentController) DeleteComments(c *gin.Context) {
	incidentClosureIDStr := c.Param("closureId")
	incidentClosureID, err := strconv.ParseUint(incidentClosureIDStr, 10, 32)
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
		utility.SendInternalErrResp(icrcc.db, c, err.Error())
		return
	}

	failedIDs, err := icrcc.service.Delete(uint(incidentClosureID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(icrcc.db, c, e.Error())
			return
		}
	}

	if len(failedIDs) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{"failed_ids": failedIDs})
		return
	}

	c.Status(http.StatusOK)
}

func (icrcc *IncidentClosureReviewCommentController) ResolveComment(c *gin.Context) {
	incidentClosureIDStr := c.Param("closureId")
	incidentClosureID, err := strconv.ParseUint(incidentClosureIDStr, 10, 32)
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
		utility.SendInternalErrResp(icrcc.db, c, err.Error())
		return
	}

	err = icrcc.service.Resolve(uint(incidentClosureID), uint(commentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (icrcc *IncidentClosureReviewCommentController) GetDocumentPreview(c *gin.Context) {
	incidentClosureIDStr := c.Param("closureId")
	incidentClosureID, err := strconv.ParseUint(incidentClosureIDStr, 10, 32)
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

	documentPath, err := icrcc.service.GetDocumentPath(uint(incidentClosureID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(icrcc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (icrcc *IncidentClosureReviewCommentController) DownloadDocument(c *gin.Context) {
	incidentClosureIDStr := c.Param("closureId")
	incidentClosureID, err := strconv.ParseUint(incidentClosureIDStr, 10, 32)
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

	documentPath, err := icrcc.service.GetDocumentPath(uint(incidentClosureID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(icrcc.db, c, err.Error())
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
