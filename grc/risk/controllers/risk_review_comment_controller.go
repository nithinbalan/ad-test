package controllers

import (
	"errors"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskReviewCommentController struct {
	service services.RiskReviewCommentServiceInterface
	db      *gorm.DB
}

func NewRiskReviewCommentController(service services.RiskReviewCommentServiceInterface) *RiskReviewCommentController {
	return &RiskReviewCommentController{
		service: service,
	}
}

func (rrcc *RiskReviewCommentController) GetAllComments(c *gin.Context) {
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
		"comment":                 "risk_review_comments.comment",
		"level":                   "risk_review_comments.level",
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

	queryParams := riskDTOs.RiskReviewCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	riskIDStr := c.Param("id")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	comments, err := rrcc.service.GetAll(uint(riskID), currentUserID, &queryParams)
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

func (rrcc *RiskReviewCommentController) GetCommentDetails(c *gin.Context) {
	riskIDStr := c.Param("id")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
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

	comment, err := rrcc.service.GetFullDetailsByID(uint(riskID), uint(commentID))
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

func (rrcc *RiskReviewCommentController) CreateComment(c *gin.Context) {
	riskIDStr := c.Param("id")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.RiskReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	id, err := rrcc.service.Create(uint(riskID), currentUserID, &formDTO)
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

func (rrcc *RiskReviewCommentController) UpdateComment(c *gin.Context) {
	riskIDStr := c.Param("id")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
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

	var formDTO riskDTOs.RiskReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrcc.db, c, err.Error())
		return
	}

	err = rrcc.service.Update(uint(riskID), uint(commentID), currentUserID, &formDTO)
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

func (rrcc *RiskReviewCommentController) DeleteComment(c *gin.Context) {
	riskIDStr := c.Param("id")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
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

	failedIDs, err := rrcc.service.Delete(uint(riskID), request.IDs, currentUserID)
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

func (rrcc *RiskReviewCommentController) ResolveComment(c *gin.Context) {

	riskIDStr := c.Param("id")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
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

	err = rrcc.service.Resolve(uint(riskID), uint(commentID), currentUserID)
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
