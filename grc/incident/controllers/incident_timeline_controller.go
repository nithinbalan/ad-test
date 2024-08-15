package controllers

import (
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentTimelineController struct {
	service services.IncidentTimelineServiceInterface
	db      *gorm.DB
}

func NewIncidentTimelineController(service services.IncidentTimelineServiceInterface) *IncidentTimelineController {
	return &IncidentTimelineController{
		service: service,
	}
}

func (itc *IncidentTimelineController) GetAllIncidentTimelines(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := incidentDTOs.IncidentTimelineListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}

	incidentTimelines, err := itc.service.GetAll(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(itc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidentTimelines)
}
