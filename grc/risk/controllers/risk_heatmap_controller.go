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

type RiskHeatmapController struct {
	service services.RiskHeatmapServiceInterface
	db      *gorm.DB
}

func NewRiskHeatmapController(service services.RiskHeatmapServiceInterface) *RiskHeatmapController {
	return &RiskHeatmapController{
		service: service,
	}
}

func (rhc *RiskHeatmapController) GetProcessRiskHeatmap(c *gin.Context) {
	idStr := c.Param("processId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskHeatmap, err := rhc.service.GetRiskHeatmapByProcessID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rhc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, riskHeatmap)
}

func (rhc *RiskHeatmapController) GetControlRiskHeatmap(c *gin.Context) {
	idStr := c.Param("controlId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskHeatmap, err := rhc.service.GetRiskHeatmapByControlID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rhc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, riskHeatmap)
}

func (rhc *RiskHeatmapController) GetComplianceRiskHeatmap(c *gin.Context) {
	idStr := c.Param("complianceId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskHeatmap, err := rhc.service.GetRiskHeatmapByComplianceID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rhc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, riskHeatmap)
}

func (rhc *RiskHeatmapController) GetRiskHeatmapByResponsibleUnit(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_category_ids",
		"risk_source_ids",
		"owner_unit_ids",
		"owner_ids",
		"risk_status_ids",
		"review_frequency_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	reportedAtFrom := c.Query("reported_at")
	reportedAtTo := c.Query("reported_at")
	nextReviewDateFrom := c.Query("next_review_date_from")
	nextReviewDateTo := c.Query("next_review_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rhc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rhc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.RiskRegisterListQueryParamsDTO{
		Sort:               sort,
		SortBy:             sortBy,
		Keyword:            keyword,
		Status:             status,
		Page:               page,
		PageSize:           pageSize,
		Filters:            filters,
		QuickFilter:        quickFilter,
		ReportedAtFrom:     reportedAtFrom,
		ReportedAtTo:       reportedAtTo,
		NextReviewDateFrom: nextReviewDateFrom,
		NextReviewDateTo:   nextReviewDateTo,
	}

	riskHeatmap, err := rhc.service.GetRiskHeatmapByResponsibleUnit(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rhc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, riskHeatmap)
}
