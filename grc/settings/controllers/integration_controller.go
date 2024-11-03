package controllers

import (
	integrationDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IntegrationController struct {
	service services.IntegrationServiceInterface
	db      *gorm.DB
}

func NewIntegrationController(service services.IntegrationServiceInterface) *IntegrationController {
	return &IntegrationController{
		service: service,
	}
}

func (icc *IntegrationController) GetAllIntegrations(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"integration_category_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := integrationDTOs.IntegrationListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	integrations, err := icc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, integrations)
}
