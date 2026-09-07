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

type ProjectDecisionLogReviewController struct {
	service services.ProjectDecisionLogReviewServiceInterface
	db      *gorm.DB
}

func NewProjectDecisionLogReviewController(service services.ProjectDecisionLogReviewServiceInterface, db *gorm.DB) *ProjectDecisionLogReviewController {
	return &ProjectDecisionLogReviewController{
		service: service,
		db:      db,
	}
}

func (pdlrc *ProjectDecisionLogReviewController) SubmitDecisionLog(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	decisionLogIDStr := c.Param("decisionLogId")
	decisionLogID, err := strconv.ParseUint(decisionLogIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pdlrc.db, c, err.Error())
		return
	}

	err = pdlrc.service.Submit(uint(projectID), uint(decisionLogID), currentUserID)
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
			utility.SendInternalErrResp(pdlrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pdlrc *ProjectDecisionLogReviewController) ResubmitDecisionLog(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	decisionLogIDStr := c.Param("decisionLogId")
	decisionLogID, err := strconv.ParseUint(decisionLogIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pdlrc.db, c, err.Error())
		return
	}

	err = pdlrc.service.Resubmit(uint(projectID), uint(decisionLogID), currentUserID)
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
			utility.SendInternalErrResp(pdlrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pdlrc *ProjectDecisionLogReviewController) ApproveDecisionLog(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	decisionLogIDStr := c.Param("decisionLogId")
	decisionLogID, err := strconv.ParseUint(decisionLogIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form projectDTOs.ProjectDecisionLogReviewFormDTO
	_ = c.ShouldBindJSON(&form)

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pdlrc.db, c, err.Error())
		return
	}

	err = pdlrc.service.Approve(uint(projectID), uint(decisionLogID), currentUserID, &form)
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
			utility.SendInternalErrResp(pdlrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pdlrc *ProjectDecisionLogReviewController) RevertDecisionLog(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	decisionLogIDStr := c.Param("decisionLogId")
	decisionLogID, err := strconv.ParseUint(decisionLogIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form projectDTOs.ProjectDecisionLogReviewFormDTO
	_ = c.ShouldBindJSON(&form)

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pdlrc.db, c, err.Error())
		return
	}

	err = pdlrc.service.Revert(uint(projectID), uint(decisionLogID), currentUserID, &form)
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
			utility.SendInternalErrResp(pdlrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pdlrc *ProjectDecisionLogReviewController) RejectDecisionLog(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	decisionLogIDStr := c.Param("decisionLogId")
	decisionLogID, err := strconv.ParseUint(decisionLogIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form projectDTOs.ProjectDecisionLogReviewFormDTO
	_ = c.ShouldBindJSON(&form)

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pdlrc.db, c, err.Error())
		return
	}

	err = pdlrc.service.Reject(uint(projectID), uint(decisionLogID), currentUserID, &form)
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
			utility.SendInternalErrResp(pdlrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
