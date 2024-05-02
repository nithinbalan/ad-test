package controllers

import (
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentStatusController struct {
	service services.IncidentStatusServiceInterface
	db      *gorm.DB
}

func NewIncidentStatusController(service services.IncidentStatusServiceInterface) *IncidentStatusController {
	return &IncidentStatusController{
		service: service,
	}
}

func (isc *IncidentStatusController) GetAllIncidentStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	incidentStatuses, err := isc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(isc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentStatuses)
}

func (isc *IncidentStatusController) GetIncidentStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentStatus, err := isc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentStatus)
}
