package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BiaProcessWorkflowController struct {
	service services.BiaProcessWorkflowServiceInterface
	db      *gorm.DB
}

func NewBiaProcessWorkflowController(service services.BiaProcessWorkflowServiceInterface) *BiaProcessWorkflowController {
	return &BiaProcessWorkflowController{
		service: service,
	}
}

func (bpwc *BiaProcessWorkflowController) GetWorkflowList(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflowList, err := bpwc.service.GetBiaProcessWorkflowList(uint(biaProcessID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bpwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflowList)
}
