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

	queryParams := riskDTOs.RiskAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	risks, err := raic.service.GetAllControls(uint(riskID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(raic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, risks)
}

func (raic *RiskAssociatedItemController) GetAllRiskComplianceRequirements(c *gin.Context) {

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

	queryParams := riskDTOs.RiskAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	risks, err := raic.service.GetAllComplianceRequirements(uint(riskID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(raic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, risks)
}
