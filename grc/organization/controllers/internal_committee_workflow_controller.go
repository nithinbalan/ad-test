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

type InternalCommitteeWorkflowController struct {
	service services.InternalCommitteeWorkflowServiceInterface
	db      *gorm.DB
}

func NewInternalCommitteeWorkflowController(service services.InternalCommitteeWorkflowServiceInterface) *InternalCommitteeWorkflowController {
	return &InternalCommitteeWorkflowController{
		service: service,
	}
}

func (icwc *InternalCommitteeWorkflowController) GetWorkflowList(c *gin.Context) {
	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := icwc.service.GetInternalCommitteeWorkflowList(uint(committeeID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}

func (icwc *InternalCommitteeWorkflowController) GetPreviousVersionWorkflowHistory(c *gin.Context) {
	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	history, err := icwc.service.GetPreviousVersionWorkflowHistory(uint(committeeID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(icwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, history)
}

func (icwc *InternalCommitteeWorkflowController) UpdateNextLevelParticipant(c *gin.Context) {
	committeeIDStr := c.Param("committeeId")
	committeeID, err := strconv.ParseUint(committeeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto dtos.UpdateNextLevelParticipantDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icwc.db, c, err.Error())
		return
	}

	if err := icwc.service.UpdateNextLevelParticipant(uint(committeeID), currentUserID, &dto); err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Next level participant updated successfully"})
}
