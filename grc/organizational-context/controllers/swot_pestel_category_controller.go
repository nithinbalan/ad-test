package controllers

import (
	"jethur/grc/organizational-context/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SwotPestelCategoryController struct {
	service services.SwotPestelCategoryServiceInterface
	db      *gorm.DB
}

func NewSwotPestelCategoryController(service services.SwotPestelCategoryServiceInterface) *SwotPestelCategoryController {
	return &SwotPestelCategoryController{
		service: service,
	}
}

func (cesc *SwotPestelCategoryController) GetAllSwotPestelCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	swotPestelCategories, err := cesc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(cesc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, swotPestelCategories)
}

func (cesc *SwotPestelCategoryController) GetSwotPestelCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	swotPestelCategory, err := cesc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, swotPestelCategory)
}
