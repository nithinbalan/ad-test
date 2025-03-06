package controllers

import (
	"jethur/grc/models"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SwotQuestionController struct {
	service services.SwotQuestionServiceInterface
	db      *gorm.DB
}

func NewSwotQuestionController(service services.SwotQuestionServiceInterface) *SwotQuestionController {
	return &SwotQuestionController{
		service: service,
	}
}

func (sqc *SwotQuestionController) GetAllSwotQuestions(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	swotQuestions, err := sqc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, swotQuestions)
}

func (sqc *SwotQuestionController) GetSwotQuestionDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	swotQuestion, err := sqc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, swotQuestion)
}

func (sqc *SwotQuestionController) CreateSwotQuestion(c *gin.Context) {
	var swotQuestion *models.SwotQuestion

	if err := c.ShouldBindJSON(&swotQuestion); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	createdSwotQuestion, err := sqc.service.Create(currentUserID, swotQuestion)
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

	c.JSON(http.StatusCreated, createdSwotQuestion)
}

func (sqc *SwotQuestionController) UpdateSwotQuestion(c *gin.Context) {
	swotQuestionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedSwotQuestion *models.SwotQuestion
	if err := c.ShouldBindJSON(&updatedSwotQuestion); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sqc.db, c, err.Error())
		return
	}

	updatedSwotQuestion, err = sqc.service.Update(uint(swotQuestionID), currentUserID, updatedSwotQuestion)
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

	c.JSON(http.StatusOK, updatedSwotQuestion)
}

func (sqc *SwotQuestionController) DeleteSwotQuestions(c *gin.Context) {
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

func (sqc *SwotQuestionController) ActivateSwotQuestions(c *gin.Context) {
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

func (sqc *SwotQuestionController) DeactivateSwotQuestions(c *gin.Context) {
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
