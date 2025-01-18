package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BcpPhaseController struct {
	service services.BcpPhaseServiceInterface
	db      *gorm.DB
}

func NewBcpPhaseController(service services.BcpPhaseServiceInterface) *BcpPhaseController {
	return &BcpPhaseController{
		service: service,
	}
}

func (icc *BcpPhaseController) GetAllBcpPhases(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	bcpPhases, err := icc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, bcpPhases)
}

func (icc *BcpPhaseController) GetBcpPhaseDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcpPhase, err := icc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, bcpPhase)
}
