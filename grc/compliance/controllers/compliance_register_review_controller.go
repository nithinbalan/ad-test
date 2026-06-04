package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceRegisterReviewController struct {
	service services.ComplianceRegisterReviewServiceInterface
	db      *gorm.DB
}

func NewComplianceRegisterReviewController(service services.ComplianceRegisterReviewServiceInterface) *ComplianceRegisterReviewController {
	return &ComplianceRegisterReviewController{
		service: service,
	}
}

func (crrc *ComplianceRegisterReviewController) SubmitComplianceRegister(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crrc.db, c, err.Error())
		return
	}

	err = crrc.service.Submit(uint(complianceRegisterID), currentUserID)
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
			utility.SendInternalErrResp(crrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (crrc *ComplianceRegisterReviewController) ApproveComplianceRegister(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crrc.db, c, err.Error())
		return
	}

	err = crrc.service.Approve(uint(complianceRegisterID), currentUserID)
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
			utility.SendInternalErrResp(crrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (crrc *ComplianceRegisterReviewController) RejectComplianceRegister(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crrc.db, c, err.Error())
		return
	}

	err = crrc.service.Reject(uint(complianceRegisterID), currentUserID)
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
			utility.SendInternalErrResp(crrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (crrc *ComplianceRegisterReviewController) ReopenComplianceRegister(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crrc.db, c, err.Error())
		return
	}

	err = crrc.service.Reopen(uint(complianceRegisterID), currentUserID)
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
			utility.SendInternalErrResp(crrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
