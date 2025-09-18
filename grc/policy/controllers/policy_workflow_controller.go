package controllers

import (
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyWorkflowController struct {
	service services.DocumentWorkflowServiceInterface
	db      *gorm.DB
}

func NewPolicyWorkflowController(service services.DocumentWorkflowServiceInterface) *PolicyWorkflowController {
	return &PolicyWorkflowController{
		service: service,
	}
}

func (dwc *PolicyWorkflowController) GetWorkflowList(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := dwc.service.GetDocumentWorkflowList(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
