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

type PlannedAdvisoryClassificationController struct {
	service services.PlannedAdvisoryClassificationServiceInterface
	db      *gorm.DB
}

func NewPlannedAdvisoryClassificationController(service services.PlannedAdvisoryClassificationServiceInterface) *PlannedAdvisoryClassificationController {
	return &PlannedAdvisoryClassificationController{
		service: service,
	}
}

func (pacc *PlannedAdvisoryClassificationController) GetAllPlannedAdvisoryClassifications(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	classifications, err := pacc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(pacc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, classifications)
}

func (pacc *PlannedAdvisoryClassificationController) GetPlannedAdvisoryClassificationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	plannedAdvisoryClassification, err := pacc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, plannedAdvisoryClassification)
}

func (pacc *PlannedAdvisoryClassificationController) CreatePlannedAdvisoryClassification(c *gin.Context) {

	var classification *models.PlannedAdvisoryClassification

	if err := c.ShouldBindJSON(&classification); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pacc.db, c, err.Error())
		return
	}

	createdClassification, err := pacc.service.Create(currentUserID, classification)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pacc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdClassification)
}

func (pacc *PlannedAdvisoryClassificationController) UpdatePlannedAdvisoryClassification(c *gin.Context) {

	classificationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedClassification *models.PlannedAdvisoryClassification
	if err := c.ShouldBindJSON(&updatedClassification); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pacc.db, c, err.Error())
		return
	}

	updatedClassification, err = pacc.service.Update(uint(classificationID), currentUserID, updatedClassification)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pacc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedClassification)
}

func (pacc *PlannedAdvisoryClassificationController) DeletePlannedAdvisoryClassifications(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pacc.db, c, err.Error())
		return
	}

	failedIDs, err := pacc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pacc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pacc *PlannedAdvisoryClassificationController) ActivatePlannedAdvisoryClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pacc.db, c, err.Error())
		return
	}

	err = pacc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pacc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (pacc *PlannedAdvisoryClassificationController) DeactivatePlannedAdvisoryClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pacc.db, c, err.Error())
		return
	}

	err = pacc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pacc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
