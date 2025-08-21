package controllers

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserDelegationController struct {
	service services.UserDelegationServiceInterface
	db      *gorm.DB
}

func NewUserDelegationController(service services.UserDelegationServiceInterface) *UserDelegationController {
	return &UserDelegationController{
		service: service,
	}
}

func (udc *UserDelegationController) GetUserDelegationHistory(c *gin.Context) {

	idStr := c.Param("userId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	UserDelegation, err := udc.service.GetAll(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(udc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, UserDelegation)
}

func (udc *UserDelegationController) CreateUserDelegation(c *gin.Context) {
	var UserDelegationFormDTO organizationDTOs.UserDelegationFormDTO
	if err := c.ShouldBindJSON(&UserDelegationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udc.db, c, err.Error())
		return
	}

	_, err = udc.service.Create(currentUserID, &UserDelegationFormDTO)
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
			utility.SendInternalErrResp(udc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (udc *UserDelegationController) GetLatestUserDelegation(c *gin.Context) {
	userIDStr := c.Param("userId")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	userDelegation, err := udc.service.GetLatestUserDelegation(uint(userID))
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
			utility.SendInternalErrResp(udc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, userDelegation)
}
