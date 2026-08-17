package controllers

import (
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentClosureWorkflowController struct {
	service services.IncidentClosureWorkflowServiceInterface
	db      *gorm.DB
}

func NewIncidentClosureWorkflowController(service services.IncidentClosureWorkflowServiceInterface) *IncidentClosureWorkflowController {
	return &IncidentClosureWorkflowController{
		service: service,
	}
}

func (icwc *IncidentClosureWorkflowController) GetWorkflowList(c *gin.Context) {
	closureIDStr := c.Param("closureId")
	closureID, err := strconv.ParseUint(closureIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := icwc.service.GetIncidentClosureWorkflowList(uint(closureID))
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
