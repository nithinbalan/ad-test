package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/risk/services"
	"jethur/utility"

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

func (rhc *RiskHeatmapController) GetControlRiskHeatmap(c *gin.Context) {
	idStr := c.Param("controlId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskHeatmap, err := rhc.service.GetRiskHeatmapByControlID(uint(id))
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

func (rhc *RiskHeatmapController) GetComplianceRiskHeatmap(c *gin.Context) {
	idStr := c.Param("complianceId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskHeatmap, err := rhc.service.GetRiskHeatmapByComplianceID(uint(id))
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

func (rhc *RiskHeatmapController) GetSLARiskHeatmap(c *gin.Context) {
	idStr := c.Param("slaId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskHeatmap, err := rhc.service.GetRiskHeatmapBySLAID(uint(id))
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

func (rhc *RiskHeatmapController) GetRiskHeatmapByResponsibleUnit(c *gin.Context) {
	riskHeatmap, err := rhc.service.GetRiskHeatmapByResponsibleUnit()
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

func (rhc *RiskHeatmapController) GetRiskHeatmapByRiskCategory(c *gin.Context) {
	riskHeatmap, err := rhc.service.GetRiskHeatmapByRiskCategory()
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

func (rhc *RiskHeatmapController) GetRiskHeatmapByRiskSource(c *gin.Context) {
	riskHeatmap, err := rhc.service.GetRiskHeatmapByRiskSource()
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

func (rhc *RiskHeatmapController) GetRiskHeatmapByRiskMatrix(c *gin.Context) {
	riskHeatmap, err := rhc.service.GetRiskHeatmapByRiskMatrix()
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

func (rhc *RiskHeatmapController) GetRiskHeatmapByResponseStrategy(c *gin.Context) {
	riskHeatmap, err := rhc.service.GetRiskHeatmapByResponseStrategy()
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

func (rhc *RiskHeatmapController) GetRiskHeatmapByResponseStatus(c *gin.Context) {
	riskHeatmap, err := rhc.service.GetRiskHeatmapByResponseStatus()
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

func (rhc *RiskHeatmapController) GetRiskHeatmapByResponseFrequency(c *gin.Context) {
	riskHeatmap, err := rhc.service.GetRiskHeatmapByResponseFrequency()
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
