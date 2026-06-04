package controllers

import (
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowCustomRoleController struct {
	service services.WorkflowCustomRoleServiceInterface
	db      *gorm.DB
}

func NewWorkflowCustomRoleController(service services.WorkflowCustomRoleServiceInterface) *WorkflowCustomRoleController {
	return &WorkflowCustomRoleController{
		service: service,
	}
}

func (wcrc *WorkflowCustomRoleController) GetAllWorkflowCustomRoles(c *gin.Context) {

	workflowCustomRoles, err := wcrc.service.GetAll()

	if err != nil {
		utility.SendInternalErrResp(wcrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, workflowCustomRoles)
}
