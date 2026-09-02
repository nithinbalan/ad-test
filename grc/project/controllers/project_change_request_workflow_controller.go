package controllers

import (
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectChangeRequestWorkflowController struct {
	service services.ProjectChangeRequestWorkflowServiceInterface
	db      *gorm.DB
}

func NewProjectChangeRequestWorkflowController(service services.ProjectChangeRequestWorkflowServiceInterface, db *gorm.DB) *ProjectChangeRequestWorkflowController {
	return &ProjectChangeRequestWorkflowController{
		service: service,
		db:      db,
	}
}

func (pcrwc *ProjectChangeRequestWorkflowController) GetWorkflowList(c *gin.Context) {
	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := pcrwc.service.GetWorkflowList(uint(changeRequestID))
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
			utility.SendInternalErrResp(pcrwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
