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

type RiskUpdatePlanReviewCommentController struct {
	service services.RiskUpdatePlanReviewCommentServiceInterface
	db      *gorm.DB
}

func NewRiskUpdatePlanReviewCommentController(service services.RiskUpdatePlanReviewCommentServiceInterface) *RiskUpdatePlanReviewCommentController {
	return &RiskUpdatePlanReviewCommentController{
		service: service,
	}
}

func (rupcc *RiskUpdatePlanReviewCommentController) GetAllComments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "comment")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"comment_status_ids",
		"risk_register_update_plan_ids",
		"level",
		"created_by_ids",
	}

	validSortFields := map[string]string{
		"comment":                 "risk_update_plan_review_comments.comment",
		"level":                   "risk_update_plan_review_comments.level",
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

	queryParams := riskDTOs.RiskUpdatePlanReviewCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rupcc.db, c, err.Error())
		return
	}

	comments, err := rupcc.service.GetAll(uint(riskRegisterUpdatePlanID), currentUserID, &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rupcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (rupcc *RiskUpdatePlanReviewCommentController) GetCommentDetails(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
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

	comment, err := rupcc.service.GetFullDetailsByID(uint(riskRegisterUpdatePlanID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rupcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (rupcc *RiskUpdatePlanReviewCommentController) CreateComment(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.RiskUpdatePlanReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rupcc.db, c, err.Error())
		return
	}

	id, err := rupcc.service.Create(uint(riskRegisterUpdatePlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rupcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (rupcc *RiskUpdatePlanReviewCommentController) UpdateComment(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
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

	var formDTO riskDTOs.RiskUpdatePlanReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rupcc.db, c, err.Error())
		return
	}

	err = rupcc.service.Update(uint(riskRegisterUpdatePlanID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rupcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rupcc *RiskUpdatePlanReviewCommentController) DeleteComments(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
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
		utility.SendInternalErrResp(rupcc.db, c, err.Error())
		return
	}

	failedIDs, err := rupcc.service.Delete(uint(riskRegisterUpdatePlanID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rupcc.db, c, e.Error())
			return
		}
	}

	if len(failedIDs) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{"failed_ids": failedIDs})
		return
	}

	c.Status(http.StatusOK)
}

func (rupcc *RiskUpdatePlanReviewCommentController) ResolveComment(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
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
		utility.SendInternalErrResp(rupcc.db, c, err.Error())
		return
	}

	err = rupcc.service.Resolve(uint(riskRegisterUpdatePlanID), uint(commentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rupcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rupcc *RiskUpdatePlanReviewCommentController) GetDocumentPreview(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
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

	documentPath, err := rupcc.service.GetDocumentPath(uint(riskRegisterUpdatePlanID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rupcc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (rupcc *RiskUpdatePlanReviewCommentController) DownloadDocument(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
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

	documentPath, err := rupcc.service.GetDocumentPath(uint(riskRegisterUpdatePlanID), uint(commentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rupcc.db, c, err.Error())
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
