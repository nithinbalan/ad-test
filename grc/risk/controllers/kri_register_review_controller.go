package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRIRegisterReviewController struct {
	service services.KRIRegisterReviewServiceInterface
	db      *gorm.DB
}

func NewKRIRegisterReviewController(service services.KRIRegisterReviewServiceInterface) *KRIRegisterReviewController {
	return &KRIRegisterReviewController{
		service: service,
	}
}

func (krrc *KRIRegisterReviewController) SubmitKRIRegister(c *gin.Context) {
	kriRegisterIDStr := c.Param("id")
	kriRegisterID, err := strconv.ParseUint(kriRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krrc.db, c, err.Error())
		return
	}

	err = krrc.service.Submit(uint(kriRegisterID), currentUserID)
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
			utility.SendInternalErrResp(krrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (krrc *KRIRegisterReviewController) ApproveKRIRegister(c *gin.Context) {
	kriRegisterIDStr := c.Param("id")
	kriRegisterID, err := strconv.ParseUint(kriRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krrc.db, c, err.Error())
		return
	}

	err = krrc.service.Approve(uint(kriRegisterID), currentUserID)
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
			utility.SendInternalErrResp(krrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (krrc *KRIRegisterReviewController) RejectKRIRegister(c *gin.Context) {
	kriRegisterIDStr := c.Param("id")
	kriRegisterID, err := strconv.ParseUint(kriRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krrc.db, c, err.Error())
		return
	}

	err = krrc.service.Reject(uint(kriRegisterID), currentUserID)
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
			utility.SendInternalErrResp(krrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
