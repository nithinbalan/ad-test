package controllers

import (
	"net/http"
	"strconv"

	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistExceptionReviewController struct {
	service services.ProjectSecurityGateChecklistExceptionReviewServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateChecklistExceptionReviewController(service services.ProjectSecurityGateChecklistExceptionReviewServiceInterface) *ProjectSecurityGateChecklistExceptionReviewController {
	return &ProjectSecurityGateChecklistExceptionReviewController{
		service: service,
	}
}

func (psgcerc *ProjectSecurityGateChecklistExceptionReviewController) SubmitProjectSecurityGateChecklistException(c *gin.Context) {
	securityGateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checklistID, err := strconv.ParseUint(c.Param("checklistId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgcerc.db, c, err.Error())
		return
	}

	err = psgcerc.service.Submit(uint(securityGateID), uint(checklistID), currentUserID)
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
			utility.SendInternalErrResp(psgcerc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgcerc *ProjectSecurityGateChecklistExceptionReviewController) ApproveProjectSecurityGateChecklistException(c *gin.Context) {
	securityGateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checklistID, err := strconv.ParseUint(c.Param("checklistId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form projectDTOs.ProjectSecurityGateChecklistExceptionReviewFormDTO
	_ = c.ShouldBindJSON(&form)

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgcerc.db, c, err.Error())
		return
	}

	err = psgcerc.service.Approve(uint(securityGateID), uint(checklistID), currentUserID, &form)
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
			utility.SendInternalErrResp(psgcerc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgcerc *ProjectSecurityGateChecklistExceptionReviewController) RevertProjectSecurityGateChecklistException(c *gin.Context) {
	securityGateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checklistID, err := strconv.ParseUint(c.Param("checklistId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form projectDTOs.ProjectSecurityGateChecklistExceptionReviewFormDTO
	_ = c.ShouldBindJSON(&form)

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgcerc.db, c, err.Error())
		return
	}

	err = psgcerc.service.Revert(uint(securityGateID), uint(checklistID), currentUserID, &form)
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
			utility.SendInternalErrResp(psgcerc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
