package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskHeatmapController struct {
	service services.RiskHeatmapServiceInterface
	db      *gorm.DB
}

func NewRiskHeatmapController(service services.RiskHeatmapServiceInterface) *RiskHeatmapController {
	return &RiskHeatmapController{
		service: service,
	}
}

func (rhc *RiskHeatmapController) GetProcessRiskHeatmap(c *gin.Context) {
	idStr := c.Param("processId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskHeatmap, err := rhc.service.GetRiskHeatmapByProcessID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rhc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, riskHeatmap)
}
