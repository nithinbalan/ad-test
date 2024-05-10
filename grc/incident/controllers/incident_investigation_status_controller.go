package controllers

import (
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentInvestigationStatusController struct {
	service services.IncidentInvestigationStatusServiceInterface
	db      *gorm.DB
}

func NewIncidentInvestigationStatusController(service services.IncidentInvestigationStatusServiceInterface) *IncidentInvestigationStatusController {
	return &IncidentInvestigationStatusController{
		service: service,
	}
}

func (iisc *IncidentInvestigationStatusController) GetAllIncidentInvestigationStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	incidentInvestigationStatuses, err := iisc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(iisc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentInvestigationStatuses)
}

func (iisc *IncidentInvestigationStatusController) GetIncidentInvestigationStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentInvestigationStatus, err := iisc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentInvestigationStatus)
}
