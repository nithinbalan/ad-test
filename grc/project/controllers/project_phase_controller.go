package controllers

import (
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectPhaseController struct {
	service services.ProjectPhaseServiceInterface
	db      *gorm.DB
}

func NewProjectPhaseController(service services.ProjectPhaseServiceInterface) *ProjectPhaseController {
	return &ProjectPhaseController{
		service: service,
	}
}

func (ppc *ProjectPhaseController) GetAllProjectPhases(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	projectPhases, err := ppc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ppc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, projectPhases)
}

func (ppc *ProjectPhaseController) GetProjectPhaseDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectPhase, err := ppc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, projectPhase)
}