package controllers

import (
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskTimelineController struct {
	service services.RiskTimelineServiceInterface
	db      *gorm.DB
}

func NewRiskTimelineController(service services.RiskTimelineServiceInterface) *RiskTimelineController {
	return &RiskTimelineController{
		service: service,
	}
}

func (rtc *RiskTimelineController) GetAllRiskTimelines(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := riskDTOs.RiskTimelineListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}

	riskTimelines, err := rtc.service.GetAll(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rtc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskTimelines)
}
