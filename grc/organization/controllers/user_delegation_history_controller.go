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

type UserDelegationHistoryController struct {
	service services.UserDelegationHistoryServiceInterface
	db      *gorm.DB
}

func NewUserDelegationHistoryController(service services.UserDelegationHistoryServiceInterface) *UserDelegationHistoryController {
	return &UserDelegationHistoryController{
		service: service,
	}
}

func (spc *UserDelegationHistoryController) GetUserDelegationHistory(c *gin.Context) {

	idStr := c.Param("userId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	UserDelegationHistory, err := spc.service.GetAll(uint(id))
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
	c.JSON(http.StatusOK, UserDelegationHistory)
}

func (spc *UserDelegationHistoryController) CreateUserDelegationHistory(c *gin.Context) {
	var UserDelegationHistoryFormDTO organizationDTOs.UserDelegationHistoryFormDTO
	if err := c.ShouldBindJSON(&UserDelegationHistoryFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	_, err = spc.service.Create(currentUserID, &UserDelegationHistoryFormDTO)
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
