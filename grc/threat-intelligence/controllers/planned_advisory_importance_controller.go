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

type PlannedAdvisoryImportanceController struct {
	service services.PlannedAdvisoryImportanceServiceInterface
	db      *gorm.DB
}

func NewPlannedAdvisoryImportanceController(service services.PlannedAdvisoryImportanceServiceInterface) *PlannedAdvisoryImportanceController {
	return &PlannedAdvisoryImportanceController{
		service: service,
	}
}

func (paic *PlannedAdvisoryImportanceController) GetAllPlannedAdvisoryImportances(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	importances, err := paic.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, importances)
}

func (paic *PlannedAdvisoryImportanceController) GetPlannedAdvisoryImportanceDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	plannedAdvisoryImportance, err := paic.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, plannedAdvisoryImportance)
}

func (paic *PlannedAdvisoryImportanceController) CreatePlannedAdvisoryImportance(c *gin.Context) {

	var importance *models.PlannedAdvisoryImportance

	if err := c.ShouldBindJSON(&importance); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	createdImportance, err := paic.service.Create(currentUserID, importance)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdImportance)
}

func (paic *PlannedAdvisoryImportanceController) UpdatePlannedAdvisoryImportance(c *gin.Context) {

	importanceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedImportance *models.PlannedAdvisoryImportance
	if err := c.ShouldBindJSON(&updatedImportance); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	updatedImportance, err = paic.service.Update(uint(importanceID), currentUserID, updatedImportance)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(paic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedImportance)
}

func (paic *PlannedAdvisoryImportanceController) DeletePlannedAdvisoryImportances(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	failedIDs, err := paic.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (paic *PlannedAdvisoryImportanceController) ActivatePlannedAdvisoryImportances(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	err = paic.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (paic *PlannedAdvisoryImportanceController) DeactivatePlannedAdvisoryImportances(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	err = paic.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
