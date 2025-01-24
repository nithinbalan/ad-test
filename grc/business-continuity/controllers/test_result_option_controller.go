package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestResultOptionController struct {
	service services.TestResultOptionServiceInterface
	db      *gorm.DB
}

func NewTestResultOptionController(service services.TestResultOptionServiceInterface) *TestResultOptionController {
	return &TestResultOptionController{
		service: service,
	}
}

func (tesc *TestResultOptionController) GetAllTestResultOptions(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	testResultOptions, err := tesc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tesc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, testResultOptions)
}

func (tesc *TestResultOptionController) GetTestResultOptionDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	testResultOption, err := tesc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, testResultOption)
}
