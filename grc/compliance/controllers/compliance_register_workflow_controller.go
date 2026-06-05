package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceRegisterWorkflowController struct {
	service services.ComplianceRegisterWorkflowServiceInterface
	db      *gorm.DB
}

func NewComplianceRegisterWorkflowController(service services.ComplianceRegisterWorkflowServiceInterface) *ComplianceRegisterWorkflowController {
	return &ComplianceRegisterWorkflowController{
		service: service,
	}
}

func (crwc *ComplianceRegisterWorkflowController) GetWorkflowList(c *gin.Context) {
	complianceRegisterIDStr := c.Param("id")
	complianceRegisterID, err := strconv.ParseUint(complianceRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := crwc.service.GetComplianceRegisterWorkflowList(uint(complianceRegisterID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
