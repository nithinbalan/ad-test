package controllers

import (
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectStatusController struct {
	service services.ProjectStatusServiceInterface
	db      *gorm.DB
}

func NewProjectStatusController(service services.ProjectStatusServiceInterface) *ProjectStatusController {
	return &ProjectStatusController{
		service: service,
	}
}

func (csc *ProjectStatusController) GetAllProjectStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	projectStatuses, err := csc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(csc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, projectStatuses)
}

func (csc *ProjectStatusController) GetProjectStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectStatus, err := csc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, projectStatus)
}
