package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskLevelController struct {
	service services.RiskLevelServiceInterface
	db      *gorm.DB
}

func NewRiskLevelController(service services.RiskLevelServiceInterface) *RiskLevelController {
	return &RiskLevelController{
		service: service,
	}
}

func (rlc *RiskLevelController) GetAllRiskLevels(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	riskLevels, err := rlc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rlc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskLevels)
}

func (rlc *RiskLevelController) GetRiskLevelDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskLevel, err := rlc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskLevel)
}
