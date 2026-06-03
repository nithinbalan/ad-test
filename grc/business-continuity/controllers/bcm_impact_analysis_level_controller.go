package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BCMImpactAnalysisLevelController struct {
	service services.BCMImpactAnalysisLevelServiceInterface
	db      *gorm.DB
}

func NewBCMImpactAnalysisLevelController(service services.BCMImpactAnalysisLevelServiceInterface) *BCMImpactAnalysisLevelController {
	return &BCMImpactAnalysisLevelController{
		service: service,
	}
}

func (bialc *BCMImpactAnalysisLevelController) GetAllBCMImpactAnalysisLevels(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	bcmImpactAnalysisLevels, err := bialc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(bialc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, bcmImpactAnalysisLevels)
}

func (bialc *BCMImpactAnalysisLevelController) GetBCMImpactAnalysisLevelDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcmImpactAnalysisLevel, err := bialc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, bcmImpactAnalysisLevel)
}
