package controllers

import (
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskAssociatedItemController struct {
	service services.RiskAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewRiskAssociatedItemController(service services.RiskAssociatedItemServiceInterface) *RiskAssociatedItemController {
	return &RiskAssociatedItemController{
		service: service,
	}
}

func (raic *RiskAssociatedItemController) GetAllRiskControls(c *gin.Context) {

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

	// Define an array of allowed parameter names
	allowedParams := []string{
		"control_category_ids",
		"control_type_ids",
		"control_implementation_status_ids",
		"control_efficiency_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Handle multiple control_category_ids
	if categoryIDs, ok := c.Request.URL.Query()["control_category_ids"]; ok {
		for _, categoryID := range categoryIDs {
			categoryIDs := strings.Split(categoryID, ",")
			filters["controls.control_category_id"] = append(filters["controls.control_category_id"], categoryIDs...)
		}
	}

	queryParams := riskDTOs.RiskAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	risks, err := raic.service.GetAllControls(uint(riskID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(raic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, risks)
}
