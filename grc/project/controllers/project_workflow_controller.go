package controllers

import (
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectWorkflowController struct {
	service services.ProjectWorkflowServiceInterface
	db      *gorm.DB
}

func NewProjectWorkflowController(service services.ProjectWorkflowServiceInterface) *ProjectWorkflowController {
	return &ProjectWorkflowController{
		service: service,
	}
}

func (pwc *ProjectWorkflowController) GetWorkflowList(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := pwc.service.GetWorkflowList(uint(projectID))
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
			utility.SendInternalErrResp(pwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
