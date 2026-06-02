package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BCMBIAImpactLevelController struct {
	service services.BCMBIAImpactLevelServiceInterface
	db      *gorm.DB
}

func NewBCMBIAImpactLevelController(service services.BCMBIAImpactLevelServiceInterface) *BCMBIAImpactLevelController {
	return &BCMBIAImpactLevelController{
		service: service,
	}
}

func (bbilc *BCMBIAImpactLevelController) GetAllBCMBIAImpactLevels(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	bcmBiaImpactLevels, err := bbilc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(bbilc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, bcmBiaImpactLevels)
}

func (bbilc *BCMBIAImpactLevelController) GetBCMBIAImpactLevelFullDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcmBiaImpactLevel, err := bbilc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, bcmBiaImpactLevel)
}
