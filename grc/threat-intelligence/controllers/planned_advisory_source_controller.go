package controllers

import (
	"jethur/grc/models"
	"jethur/grc/threat-intelligence/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlannedAdvisorySourceController struct {
	service services.PlannedAdvisorySourceServiceInterface
	db      *gorm.DB
}

func NewPlannedAdvisorySourceController(service services.PlannedAdvisorySourceServiceInterface) *PlannedAdvisorySourceController {
	return &PlannedAdvisorySourceController{
		service: service,
	}
}

func (pasc *PlannedAdvisorySourceController) GetAllPlannedAdvisorySources(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	sources, err := pasc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, sources)
}

func (pasc *PlannedAdvisorySourceController) GetPlannedAdvisorySourceDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	plannedAdvisorySource, err := pasc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, plannedAdvisorySource)
}

func (pasc *PlannedAdvisorySourceController) CreatePlannedAdvisorySource(c *gin.Context) {

	var source *models.PlannedAdvisorySource

	if err := c.ShouldBindJSON(&source); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}

	createdSource, err := pasc.service.Create(currentUserID, source)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pasc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdSource)
}

func (pasc *PlannedAdvisorySourceController) UpdatePlannedAdvisorySource(c *gin.Context) {

	sourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedSource *models.PlannedAdvisorySource
	if err := c.ShouldBindJSON(&updatedSource); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}

	updatedSource, err = pasc.service.Update(uint(sourceID), currentUserID, updatedSource)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pasc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedSource)
}

func (pasc *PlannedAdvisorySourceController) DeletePlannedAdvisorySources(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}

	failedIDs, err := pasc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pasc *PlannedAdvisorySourceController) ActivatePlannedAdvisorySources(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}

	err = pasc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (pasc *PlannedAdvisorySourceController) DeactivatePlannedAdvisorySources(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}

	err = pasc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
