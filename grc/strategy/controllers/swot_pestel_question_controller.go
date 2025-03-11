package controllers

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SwotPestelQuestionController struct {
	service services.SwotPestelQuestionServiceInterface
	db      *gorm.DB
}

func NewSwotPestelQuestionController(service services.SwotPestelQuestionServiceInterface) *SwotPestelQuestionController {
	return &SwotPestelQuestionController{
		service: service,
	}
}

func (sqc *SwotPestelQuestionController) GetAllSwotPestelQuestions(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	isSwot := c.Query("is_swot")
	isPestel := c.Query("is_pestel")

	swotPestelQuestions, err := sqc.service.GetAll(sort, sortBy, keyword, status, isSwot, isPestel)

	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, swotPestelQuestions)
}

func (sqc *SwotPestelQuestionController) GetSwotPestelQuestionDetails(c *gin.Context) {
	idStr := c.Param("questionId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	swotPestelQuestion, err := sqc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, swotPestelQuestion)
}

func (sqc *SwotPestelQuestionController) CreateSwotPestelQuestion(c *gin.Context) {
	var swotPestelQuestion *models.SwotPestelQuestion

	if err := c.ShouldBindJSON(&swotPestelQuestion); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	createdSwotPestelQuestion, err := sqc.service.Create(currentUserID, swotPestelQuestion)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sqc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdSwotPestelQuestion)
}

func (sqc *SwotPestelQuestionController) UpdateSwotPestelQuestion(c *gin.Context) {
	swotPestelQuestionID, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedSwotPestelQuestion *models.SwotPestelQuestion
	if err := c.ShouldBindJSON(&updatedSwotPestelQuestion); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	updatedSwotPestelQuestion, err = sqc.service.Update(uint(swotPestelQuestionID), currentUserID, updatedSwotPestelQuestion)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sqc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedSwotPestelQuestion)
}

func (sqc *SwotPestelQuestionController) DeleteSwotPestelQuestions(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	failedIDs, err := sqc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (sqc *SwotPestelQuestionController) ActivateSwotPestelQuestions(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	err = sqc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (sqc *SwotPestelQuestionController) DeactivateSwotPestelQuestions(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	err = sqc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (sqc *SwotPestelQuestionController) UpdateAnswer(c *gin.Context) {
	var answerFormDTO strategyDTOs.AnswerFormDTO

	questionIDStr := c.Param("id")

	questionId, err := strconv.ParseUint(questionIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&answerFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}
	err = sqc.service.UpdateAnswer(uint(questionId), currentUserID, &answerFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sqc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (sqc *SwotPestelQuestionController) CreateInvitees(c *gin.Context) {
	var SwotPestelQuestionInvitee *strategyDTOs.InviteeFormDTO

	questionIDStr := c.Param("questionId")

	questionId, err := strconv.ParseUint(questionIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&SwotPestelQuestionInvitee); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	err = sqc.service.CreateInvitees(uint(questionId), currentUserID, SwotPestelQuestionInvitee)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sqc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, questionId)
}
