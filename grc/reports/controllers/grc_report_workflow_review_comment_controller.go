package controllers

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GRCReportWorkflowReviewCommentController struct {
	service services.GRCReportWorkflowReviewCommentServiceInterface
	db      *gorm.DB
}

func NewGRCReportWorkflowReviewCommentController(service services.GRCReportWorkflowReviewCommentServiceInterface) *GRCReportWorkflowReviewCommentController {
	return &GRCReportWorkflowReviewCommentController{service: service}
}

func (grwrcc *GRCReportWorkflowReviewCommentController) GetAllComments(c *gin.Context) {
	grcReportIDStr := c.Param("grcReportId")
	grcReportID, err := strconv.ParseUint(grcReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "comment")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	section := c.Query("type")

	allowedParams := []string{
		"grc_report_workflow_ids",
	}
	validSortFields := map[string]string{
		"comment": "grc_report_workflow_review_comments.comment",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := dtos.GRCReportWorkflowReviewCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Filters:     filters,
		QuickFilter: quickFilter,
		Section:     section,
	}

	comments, err := grwrcc.service.GetAll(uint(grcReportID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(grwrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (grwrcc *GRCReportWorkflowReviewCommentController) GetCommentDetails(c *gin.Context) {
	grcReportIDStr := c.Param("grcReportId")
	grcReportID, err := strconv.ParseUint(grcReportIDStr, 10, 32)
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

	comment, err := grwrcc.service.GetFullDetailsByID(uint(grcReportID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(grwrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (grwrcc *GRCReportWorkflowReviewCommentController) CreateComment(c *gin.Context) {
	grcReportIDStr := c.Param("grcReportId")
	grcReportID, err := strconv.ParseUint(grcReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.GRCReportWorkflowReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(grwrcc.db, c, err.Error())
		return
	}

	id, err := grwrcc.service.Create(uint(grcReportID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(grwrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (grwrcc *GRCReportWorkflowReviewCommentController) UpdateComment(c *gin.Context) {
	grcReportIDStr := c.Param("grcReportId")
	grcReportID, err := strconv.ParseUint(grcReportIDStr, 10, 32)
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

	var formDTO dtos.GRCReportWorkflowReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(grwrcc.db, c, err.Error())
		return
	}

	err = grwrcc.service.Update(uint(grcReportID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(grwrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (grwrcc *GRCReportWorkflowReviewCommentController) DeleteComment(c *gin.Context) {
	grcReportIDStr := c.Param("grcReportId")
	grcReportID, err := strconv.ParseUint(grcReportIDStr, 10, 32)
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
		utility.SendInternalErrResp(grwrcc.db, c, err.Error())
		return
	}

	failedIDs, err := grwrcc.service.Delete(uint(grcReportID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(grwrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (grwrcc *GRCReportWorkflowReviewCommentController) GRCReportWorkflowReviewCommentResolve(c *gin.Context) {

	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
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
		utility.SendInternalErrResp(grwrcc.db, c, err.Error())
		return
	}

	err = grwrcc.service.GRCReportSectionCommentResolve(uint(grcReportID), uint(commentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(grwrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (grwrcc *GRCReportWorkflowReviewCommentController) GRCReportWorkflowReviewCommentClose(c *gin.Context) {

	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
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
		utility.SendInternalErrResp(grwrcc.db, c, err.Error())
		return
	}

	err = grwrcc.service.GRCReportSectionCommentClose(uint(grcReportID), uint(commentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(grwrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (grwrcc *GRCReportWorkflowReviewCommentController) GRCReportSectionReviewed(c *gin.Context) {

	fmt.Println("123")
	var commentForm dtos.GRCReportWorkflowReviewCommentFormDTO

	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&commentForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(grwrcc.db, c, err.Error())
		return
	}

	err = grwrcc.service.GRCReportSectionReviewed(uint(grcReportID), currentUserID, &commentForm)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(grwrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (grwrcc *GRCReportWorkflowReviewCommentController) GRCReportSectionApproved(c *gin.Context) {
	var commentForm dtos.GRCReportWorkflowReviewCommentFormDTO

	grcReportID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&commentForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(grwrcc.db, c, err.Error())
		return
	}

	err = grwrcc.service.GRCReportSectionApproved(uint(grcReportID), currentUserID, &commentForm)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(grwrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
