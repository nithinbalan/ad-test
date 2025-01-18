package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BcpStatusController struct {
	service services.BcpStatusServiceInterface
	db      *gorm.DB
}

func NewBcpStatusController(service services.BcpStatusServiceInterface) *BcpStatusController {
	return &BcpStatusController{
		service: service,
	}
}

func (bsc *BcpStatusController) GetAllBcpStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	bcpStatuses, err := bsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(bsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, bcpStatuses)
}

func (bsc *BcpStatusController) GetBcpStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcpStatus, err := bsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, bcpStatus)
}
