package controllers

import (
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectChangeRequestReviewController struct {
	service services.ProjectChangeRequestReviewServiceInterface
	db      *gorm.DB
}

func NewProjectChangeRequestReviewController(service services.ProjectChangeRequestReviewServiceInterface, db *gorm.DB) *ProjectChangeRequestReviewController {
	return &ProjectChangeRequestReviewController{
		service: service,
		db:      db,
	}
}

func (pcrrc *ProjectChangeRequestReviewController) SubmitChangeRequest(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcrrc.db, c, err.Error())
		return
	}

	err = pcrrc.service.Submit(uint(projectID), uint(changeRequestID), currentUserID)
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
			utility.SendInternalErrResp(pcrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pcrrc *ProjectChangeRequestReviewController) ResubmitChangeRequest(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcrrc.db, c, err.Error())
		return
	}

	err = pcrrc.service.Resubmit(uint(projectID), uint(changeRequestID), currentUserID)
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
			utility.SendInternalErrResp(pcrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pcrrc *ProjectChangeRequestReviewController) ApproveChangeRequest(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form projectDTOs.ProjectChangeRequestReviewFormDTO
	_ = c.ShouldBindJSON(&form)

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcrrc.db, c, err.Error())
		return
	}

	err = pcrrc.service.Approve(uint(projectID), uint(changeRequestID), currentUserID, &form)
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
			utility.SendInternalErrResp(pcrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pcrrc *ProjectChangeRequestReviewController) RevertChangeRequest(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form projectDTOs.ProjectChangeRequestReviewFormDTO
	_ = c.ShouldBindJSON(&form)

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcrrc.db, c, err.Error())
		return
	}

	err = pcrrc.service.Revert(uint(projectID), uint(changeRequestID), currentUserID, &form)
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
			utility.SendInternalErrResp(pcrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pcrrc *ProjectChangeRequestReviewController) RejectChangeRequest(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form projectDTOs.ProjectChangeRequestReviewFormDTO
	_ = c.ShouldBindJSON(&form)

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcrrc.db, c, err.Error())
		return
	}

	err = pcrrc.service.Reject(uint(projectID), uint(changeRequestID), currentUserID, &form)
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
			utility.SendInternalErrResp(pcrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
