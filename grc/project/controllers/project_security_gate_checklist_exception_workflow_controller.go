package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistExceptionWorkflowController struct {
	service services.ProjectSecurityGateChecklistExceptionWorkflowServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateChecklistExceptionWorkflowController(service services.ProjectSecurityGateChecklistExceptionWorkflowServiceInterface) *ProjectSecurityGateChecklistExceptionWorkflowController {
	return &ProjectSecurityGateChecklistExceptionWorkflowController{
		service: service,
	}
}

func (psgcewc *ProjectSecurityGateChecklistExceptionWorkflowController) GetProjectSecurityGateChecklistExceptionWorkflowList(c *gin.Context) {
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

	workflows, err := psgcewc.service.GetWorkflowList(uint(securityGateID), uint(checklistID))
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
			utility.SendInternalErrResp(psgcewc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
