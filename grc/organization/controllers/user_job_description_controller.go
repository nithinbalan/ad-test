package controllers

import (
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserJobDescriptionController struct {
	service services.UserJobDescriptionServiceInterface
	db      *gorm.DB
}

func NewUserJobDescriptionController(service services.UserJobDescriptionServiceInterface) *UserJobDescriptionController {
	return &UserJobDescriptionController{
		service: service,
	}
}

func (ujdc *UserJobDescriptionController) AssignJobDescriptionsToUser(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utility.SendInvalidErrResp(c, utility.NewValidationError(err, "user_id", "Invalid user ID"))
		return
	}

	var input dtos.AssignJobDescriptionsDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, utility.NewValidationError(err, "job_description_ids", "Invalid input"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ujdc.db, c, err.Error())
		return
	}

	if err := ujdc.service.Create(uint(userID), &input, currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ujdc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (ctrl *UserJobDescriptionController) DeleteJobDescriptionUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	jdID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ctrl.db, c, err.Error())
		return
	}

	if err := ctrl.service.Delete(uint(userID), uint(jdID), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ctrl.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}
