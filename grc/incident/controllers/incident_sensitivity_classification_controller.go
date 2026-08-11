package controllers

import (
	"jethur/grc/incident/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentSensitivityClassificationController struct {
	service services.IncidentSensitivityClassificationServiceInterface
	db      *gorm.DB
}

func NewIncidentSensitivityClassificationController(service services.IncidentSensitivityClassificationServiceInterface) *IncidentSensitivityClassificationController {
	return &IncidentSensitivityClassificationController{
		service: service,
	}
}

func (iscc *IncidentSensitivityClassificationController) GetAllIncidentSensitivityClassifications(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	incidentSensitivityClassifications, err := iscc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentSensitivityClassifications)
}

func (iscc *IncidentSensitivityClassificationController) GetIncidentSensitivityClassificationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentSensitivityClassification, err := iscc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentSensitivityClassification)
}

func (iscc *IncidentSensitivityClassificationController) CreateIncidentSensitivityClassification(c *gin.Context) {
	var classification *models.IncidentSensitivityClassification

	if err := c.ShouldBindJSON(&classification); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	createdClassification, err := iscc.service.Create(currentUserID, classification)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iscc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdClassification)
}

func (iscc *IncidentSensitivityClassificationController) UpdateIncidentSensitivityClassification(c *gin.Context) {
	var updatedClassification *models.IncidentSensitivityClassification

	classificationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedClassification); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	updatedClassification, err = iscc.service.Update(uint(classificationID), currentUserID, updatedClassification)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iscc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedClassification)
}

func (iscc *IncidentSensitivityClassificationController) DeleteIncidentSensitivityClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	failedIDs, err := iscc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (iscc *IncidentSensitivityClassificationController) ActivateIncidentSensitivityClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	err = iscc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (iscc *IncidentSensitivityClassificationController) DeactivateIncidentSensitivityClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	err = iscc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
