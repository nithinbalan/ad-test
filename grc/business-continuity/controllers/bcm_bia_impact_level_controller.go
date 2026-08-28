package controllers

import (
	bcmDTOs "jethur/grc/business-continuity/dtos"
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

func (bbilc *BCMBIAImpactLevelController) GetBCMBIAImpactLevelByImpactAnalysisLevelID(c *gin.Context) {
	impactAnalysisLevelIDStr := c.Param("analysisLevelId")

	impactAnalysisLevelID, err := strconv.ParseUint(impactAnalysisLevelIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcmBiaImpactLevel, err := bbilc.service.GetByImpactAnalysisLevelID(uint(impactAnalysisLevelID))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, bcmBiaImpactLevel)
}

func (bbilc *BCMBIAImpactLevelController) UpdateBCMBIAImpactLevels(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bbilc.db, c, err.Error())
		return
	}

	var req bcmDTOs.BCMBIAImpactLevelBulkUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	updatedLevels, err := bbilc.service.UpdateImpactLevels(uint(currentUserID), &req)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, err)
		default:
			utility.SendInternalErrResp(bbilc.db, c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, updatedLevels)
}
