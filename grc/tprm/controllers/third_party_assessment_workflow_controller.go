package controllers

import (
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyAssessmentWorkflowController struct {
	service services.ThirdPartyAssessmentWorkflowServiceInterface
	db      *gorm.DB
}

func NewThirdPartyAssessmentWorkflowController(service services.ThirdPartyAssessmentWorkflowServiceInterface) *ThirdPartyAssessmentWorkflowController {
	return &ThirdPartyAssessmentWorkflowController{
		service: service,
	}
}

func (arwc *ThirdPartyAssessmentWorkflowController) GetWorkflowList(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentID")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := arwc.service.GetTprmWorkflowList(uint(assessmentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
