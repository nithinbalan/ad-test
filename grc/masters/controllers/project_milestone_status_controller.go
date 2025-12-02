package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectMilestoneStatusController struct {
	service services.ProjectMilestoneStatusServiceInterface
	db      *gorm.DB
}

func NewProjectMilestoneStatusController(service services.ProjectMilestoneStatusServiceInterface) *ProjectMilestoneStatusController {
	return &ProjectMilestoneStatusController{
		service: service,
	}
}

func (pms *ProjectMilestoneStatusController) GetAllProjectMilestoneStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	projectMilestoneStatuses, err := pms.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(pms.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, projectMilestoneStatuses)
}

func (pms *ProjectMilestoneStatusController) GetProjectMilestoneStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectMilestoneStatus, err := pms.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, projectMilestoneStatus)
}
