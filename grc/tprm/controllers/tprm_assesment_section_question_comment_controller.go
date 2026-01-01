package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type TprmAssessmentSectionQuestionCommentController struct {
	service services.TprmAssessmentSectionQuestionCommentServiceInterface
	db      *gorm.DB
}

func NewTprmAssessmentSectionQuestionCommentController(
	service services.TprmAssessmentSectionQuestionCommentServiceInterface,
) *TprmAssessmentSectionQuestionCommentController {

	return &TprmAssessmentSectionQuestionCommentController{
		service: service,
	}
}

func (tasqcc *TprmAssessmentSectionQuestionCommentController) GetAllComments(c *gin.Context) {
	questionIDStr := c.Param("questionId")
	assessmentIDStr := c.Param("assessmentID")
	assessmentID, _ := strconv.ParseUint(assessmentIDStr, 10, 32)

	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Read query parameters
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	// Filters allowed
	allowedParams := []string{"created_by_ids", "level"}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	params := &dtos.TprmAssessmentSectionQuestionCommentListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Filters: filters,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	comments, err := tasqcc.service.GetAll(uint(questionID), uint(assessmentID), currentUserID, params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (tasqcc *TprmAssessmentSectionQuestionCommentController) GetAllCommentsByAssessment(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentID")
	assessmentID, _ := strconv.ParseUint(assessmentIDStr, 10, 32)

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tasqcc.db, c, err.Error())
		return
	}

	allowed := []string{"section_ids", "tprm_assessment_section_question_ids", "level", "created_by_ids"}
	filters := utility.ParseQueryFilters(allowed, c.Request.URL.Query())

	params := &dtos.TprmAssessmentSectionQuestionCommentListQueryParamsDTO{
		Keyword: c.Query("keyword"),
		Filters: filters,
	}

	comments, err := tasqcc.service.GetAllByAssessmentID(uint(assessmentID), currentUserID, params)
	if err != nil {
		utility.SendInternalErrResp(tasqcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (tasqcc *TprmAssessmentSectionQuestionCommentController) GetCommentDetails(c *gin.Context) {
	questionIDStr := c.Param("questionId")
	assessmentIDStr := c.Param("assessmentID")
	assessmentID, _ := strconv.ParseUint(assessmentIDStr, 10, 32)
	questionId, _ := strconv.ParseUint(questionIDStr, 10, 64)
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	detail, err := tasqcc.service.GetFullDetailsByID(uint(questionId), uint(assessmentID), uint(commentID))
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, detail)
}

func (tasqcc *TprmAssessmentSectionQuestionCommentController) CreateComment(c *gin.Context) {

	questionIDStr := c.Param("questionId")
	assessmentIDStr := c.Param("assessmentID")
	assessmentID, _ := strconv.ParseUint(assessmentIDStr, 10, 32)
	questionId, _ := strconv.ParseUint(questionIDStr, 10, 64)

	var dto dtos.TprmAssessmentSectionQuestionCommentFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	id, err := tasqcc.service.Create(uint(questionId), uint(assessmentID), &dto, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (tasqcc *TprmAssessmentSectionQuestionCommentController) UpdateComment(c *gin.Context) {

	assessmentIDStr := c.Param("assessmentID")
	assessmentID, _ := strconv.ParseUint(assessmentIDStr, 10, 32)
	idStr := c.Param("commentId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto dtos.TprmAssessmentSectionQuestionCommentFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	err = tasqcc.service.Update(uint(assessmentID), uint(id), currentUserID, &dto)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(nil, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment updated successfully"})
}

func (tasqcc *TprmAssessmentSectionQuestionCommentController) DeleteComments(c *gin.Context) {

	assessmentIDStr := c.Param("assessmentID")
	assessmentID, _ := strconv.ParseUint(assessmentIDStr, 10, 32)
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	failedIDs, err := tasqcc.service.Delete(uint(assessmentID), request.IDs, currentUserID)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(nil, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tasqcc *TprmAssessmentSectionQuestionCommentController) CreateAssessmentComment(c *gin.Context) {

	assessmentIDStr := c.Param("assessmentID")
	assessmentID, _ := strconv.ParseUint(assessmentIDStr, 10, 32)

	var dto dtos.TprmAssessmentSectionQuestionCommentFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	id, err := tasqcc.service.CreateAssessmentComment(uint(assessmentID), &dto, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}
