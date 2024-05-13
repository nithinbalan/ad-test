package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskAppetiteToleranceLevelController struct {
	service services.RiskAppetiteToleranceLevelServiceInterface
	db      *gorm.DB
}

func NewRiskAppetiteToleranceLevelController(service services.RiskAppetiteToleranceLevelServiceInterface) *RiskAppetiteToleranceLevelController {
	return &RiskAppetiteToleranceLevelController{
		service: service,
	}
}

func (ratlc *RiskAppetiteToleranceLevelController) GetAllRiskAppetiteToleranceLevels(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	riskAppetiteToleranceLevels, err := ratlc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ratlc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskAppetiteToleranceLevels)
}

func (ratlc *RiskAppetiteToleranceLevelController) GetRiskAppetiteToleranceLevelDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskAppetiteToleranceLevel, err := ratlc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskAppetiteToleranceLevel)
}
