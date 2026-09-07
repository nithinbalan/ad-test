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

type PlannedAdvisoryTypeController struct {
	service services.PlannedAdvisoryTypeServiceInterface
	db      *gorm.DB
}

func NewPlannedAdvisoryTypeController(service services.PlannedAdvisoryTypeServiceInterface) *PlannedAdvisoryTypeController {
	return &PlannedAdvisoryTypeController{
		service: service,
	}
}

func (patc *PlannedAdvisoryTypeController) GetAllPlannedAdvisoryTypes(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	advisoryTypes, err := patc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(patc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, advisoryTypes)
}

func (patc *PlannedAdvisoryTypeController) GetPlannedAdvisoryTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	plannedAdvisoryType, err := patc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, plannedAdvisoryType)
}

func (patc *PlannedAdvisoryTypeController) CreatePlannedAdvisoryType(c *gin.Context) {

	var advisoryType *models.PlannedAdvisoryType

	if err := c.ShouldBindJSON(&advisoryType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(patc.db, c, err.Error())
		return
	}

	createdAdvisoryType, err := patc.service.Create(currentUserID, advisoryType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(patc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdAdvisoryType)
}

func (patc *PlannedAdvisoryTypeController) UpdatePlannedAdvisoryType(c *gin.Context) {

	advisoryTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedAdvisoryType *models.PlannedAdvisoryType
	if err := c.ShouldBindJSON(&updatedAdvisoryType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(patc.db, c, err.Error())
		return
	}

	updatedAdvisoryType, err = patc.service.Update(uint(advisoryTypeID), currentUserID, updatedAdvisoryType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(patc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedAdvisoryType)
}

func (patc *PlannedAdvisoryTypeController) DeletePlannedAdvisoryTypes(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(patc.db, c, err.Error())
		return
	}

	failedIDs, err := patc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(patc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (patc *PlannedAdvisoryTypeController) ActivatePlannedAdvisoryTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(patc.db, c, err.Error())
		return
	}

	err = patc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(patc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (patc *PlannedAdvisoryTypeController) DeactivatePlannedAdvisoryTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(patc.db, c, err.Error())
		return
	}

	err = patc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(patc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
