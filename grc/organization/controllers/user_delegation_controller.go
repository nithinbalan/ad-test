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

func (spc *UserDelegationController) GetUserDelegationHistory(c *gin.Context) {

	idStr := c.Param("userId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	UserDelegation, err := spc.service.GetAll(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, UserDelegation)
}

func (spc *UserDelegationController) CreateUserDelegation(c *gin.Context) {
	var UserDelegationFormDTO organizationDTOs.UserDelegationFormDTO
	if err := c.ShouldBindJSON(&UserDelegationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	_, err = spc.service.Create(currentUserID, &UserDelegationFormDTO)
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
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
