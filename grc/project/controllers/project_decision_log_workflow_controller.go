package controllers

import (
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectDecisionLogWorkflowController struct {
	service services.ProjectDecisionLogWorkflowServiceInterface
	db      *gorm.DB
}

func NewProjectDecisionLogWorkflowController(service services.ProjectDecisionLogWorkflowServiceInterface, db *gorm.DB) *ProjectDecisionLogWorkflowController {
	return &ProjectDecisionLogWorkflowController{
		service: service,
		db:      db,
	}
}

func (pdlwc *ProjectDecisionLogWorkflowController) GetWorkflowList(c *gin.Context) {
	decisionLogIDStr := c.Param("decisionLogId")
	decisionLogID, err := strconv.ParseUint(decisionLogIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := pdlwc.service.GetWorkflowList(uint(decisionLogID))
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
			utility.SendInternalErrResp(pdlwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
