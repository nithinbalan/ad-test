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

type PlannedAdvisorySecurityClassificationController struct {
	service services.PlannedAdvisorySecurityClassificationServiceInterface
	db      *gorm.DB
}

func NewPlannedAdvisorySecurityClassificationController(service services.PlannedAdvisorySecurityClassificationServiceInterface) *PlannedAdvisorySecurityClassificationController {
	return &PlannedAdvisorySecurityClassificationController{
		service: service,
	}
}

func (pasc *PlannedAdvisorySecurityClassificationController) GetAllPlannedAdvisorySecurityClassifications(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	classifications, err := pasc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, classifications)
}

func (pasc *PlannedAdvisorySecurityClassificationController) GetPlannedAdvisorySecurityClassificationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	plannedAdvisorySecurityClassification, err := pasc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, plannedAdvisorySecurityClassification)
}

func (pasc *PlannedAdvisorySecurityClassificationController) CreatePlannedAdvisorySecurityClassification(c *gin.Context) {

	var classification *models.PlannedAdvisorySecurityClassification

	if err := c.ShouldBindJSON(&classification); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}

	createdClassification, err := pasc.service.Create(currentUserID, classification)
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

	c.JSON(http.StatusCreated, createdClassification)
}

func (pasc *PlannedAdvisorySecurityClassificationController) UpdatePlannedAdvisorySecurityClassification(c *gin.Context) {

	classificationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedClassification *models.PlannedAdvisorySecurityClassification
	if err := c.ShouldBindJSON(&updatedClassification); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}

	updatedClassification, err = pasc.service.Update(uint(classificationID), currentUserID, updatedClassification)
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

	c.JSON(http.StatusOK, updatedClassification)
}

func (pasc *PlannedAdvisorySecurityClassificationController) DeletePlannedAdvisorySecurityClassifications(c *gin.Context) {

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

func (pasc *PlannedAdvisorySecurityClassificationController) ActivatePlannedAdvisorySecurityClassifications(c *gin.Context) {
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

func (pasc *PlannedAdvisorySecurityClassificationController) DeactivatePlannedAdvisorySecurityClassifications(c *gin.Context) {
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
