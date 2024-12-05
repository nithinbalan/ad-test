package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IntegrationCategoryController struct {
	service services.IntegrationCategoryServiceInterface
	db      *gorm.DB
}

func NewIntegrationCategoryController(service services.IntegrationCategoryServiceInterface) *IntegrationCategoryController {
	return &IntegrationCategoryController{
		service: service,
	}
}

func (icc *IntegrationCategoryController) GetAllIntegrationCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	integrationCategories, err := icc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, integrationCategories)
}
