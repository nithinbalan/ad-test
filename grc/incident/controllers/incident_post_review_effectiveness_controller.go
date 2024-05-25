package controllers

import (
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentPostReviewEffectivenessController struct {
	service services.IncidentPostReviewEffectivenessServiceInterface
	db      *gorm.DB
}

func NewIncidentPostReviewEffectivenessController(service services.IncidentPostReviewEffectivenessServiceInterface) *IncidentPostReviewEffectivenessController {
	return &IncidentPostReviewEffectivenessController{
		service: service,
	}
}

func (iprec *IncidentPostReviewEffectivenessController) GetAllIncidentPostReviewEffectivenesses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	incidentPostReviewEffectivenesses, err := iprec.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(iprec.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentPostReviewEffectivenesses)
}

func (iprec *IncidentPostReviewEffectivenessController) GetIncidentPostReviewEffectivenessDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentPostReviewEffectiveness, err := iprec.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentPostReviewEffectiveness)
}
