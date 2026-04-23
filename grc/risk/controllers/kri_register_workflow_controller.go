package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRIRegisterWorkflowController struct {
	service services.KRIRegisterWorkflowServiceInterface
	db      *gorm.DB
}

func NewKRIRegisterWorkflowController(service services.KRIRegisterWorkflowServiceInterface) *KRIRegisterWorkflowController {
	return &KRIRegisterWorkflowController{
		service: service,
	}
}

func (krwc *KRIRegisterWorkflowController) GetWorkflowList(c *gin.Context) {
	kriRegisterIDStr := c.Param("id")
	kriRegisterID, err := strconv.ParseUint(kriRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := krwc.service.GetKRIRegisterWorkflowList(uint(kriRegisterID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(krwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
