package controllers

import (
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskResponseWorkflowController struct {
	service services.AssetRiskResponseWorkflowServiceInterface
	db      *gorm.DB
}

func NewAssetRiskResponseWorkflowController(service services.AssetRiskResponseWorkflowServiceInterface) *AssetRiskResponseWorkflowController {
	return &AssetRiskResponseWorkflowController{
		service: service,
	}
}

func (arwc *AssetRiskResponseWorkflowController) GetWorkflowList(c *gin.Context) {
	riskIDStr := c.Param("id")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := arwc.service.GetRiskResponseWorkflowList(uint(riskID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
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
