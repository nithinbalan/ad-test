package controllers

import (
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentClosureReviewController struct {
	service services.IncidentClosureReviewServiceInterface
	db      *gorm.DB
}

func NewIncidentClosureReviewController(service services.IncidentClosureReviewServiceInterface) *IncidentClosureReviewController {
	return &IncidentClosureReviewController{
		service: service,
	}
}

func (icrc *IncidentClosureReviewController) SubmitIncidentClosure(c *gin.Context) {
	closureIDStr := c.Param("closureId")
	closureID, err := strconv.ParseUint(closureIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icrc.db, c, err.Error())
		return
	}

	err = icrc.service.Submit(uint(closureID), currentUserID)
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
			utility.SendInternalErrResp(icrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (icrc *IncidentClosureReviewController) ApproveIncidentClosure(c *gin.Context) {
	closureIDStr := c.Param("closureId")
	closureID, err := strconv.ParseUint(closureIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icrc.db, c, err.Error())
		return
	}

	err = icrc.service.Approve(uint(closureID), currentUserID)
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
			utility.SendInternalErrResp(icrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (icrc *IncidentClosureReviewController) RevertIncidentClosure(c *gin.Context) {
	closureIDStr := c.Param("closureId")
	closureID, err := strconv.ParseUint(closureIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icrc.db, c, err.Error())
		return
	}

	err = icrc.service.Revert(uint(closureID), currentUserID)
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
			utility.SendInternalErrResp(icrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
