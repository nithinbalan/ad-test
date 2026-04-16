package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlEffectivenessRatingController struct {
	service services.ControlEffectivenessRatingServiceInterface
	db      *gorm.DB
}

func NewControlEffectivenessRatingController(service services.ControlEffectivenessRatingServiceInterface) *ControlEffectivenessRatingController {
	return &ControlEffectivenessRatingController{
		service: service,
	}
}

func (cerc *ControlEffectivenessRatingController) GetAllControlEffectivenessRatings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	typeFilter := c.Query("type")

	controlEffectivenessRatings, err := cerc.service.GetAll(sort, sortBy, keyword, status, typeFilter)

	if err != nil {
		utility.SendInternalErrResp(cerc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, controlEffectivenessRatings)
}

func (cerc *ControlEffectivenessRatingController) GetControlEffectivenessRatingDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlEffectivenessRating, err := cerc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, controlEffectivenessRating)
}
