package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProcessDependencyTypeController struct {
	service services.ProcessDependencyTypeServiceInterface
	db      *gorm.DB
}

func NewProcessDependencyTypeController(service services.ProcessDependencyTypeServiceInterface) *ProcessDependencyTypeController {
	return &ProcessDependencyTypeController{
		service: service,
	}
}

func (pdtc *ProcessDependencyTypeController) GetAllProcessDependencyTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	processDependencyTypes, err := pdtc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(pdtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, processDependencyTypes)
}

func (pdtc *ProcessDependencyTypeController) GetProcessDependencyTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	processDependencyType, err := pdtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, processDependencyType)
}
