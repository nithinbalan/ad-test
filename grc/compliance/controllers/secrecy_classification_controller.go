package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SecrecyClassificationController struct {
	service services.SecrecyClassificationServiceInterface
	db      *gorm.DB
}

func NewSecrecyClassificationController(service services.SecrecyClassificationServiceInterface) *SecrecyClassificationController {
	return &SecrecyClassificationController{
		service: service,
	}
}

func (scc *SecrecyClassificationController) GetAllSecrecyClassifications(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	classifications, err := scc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, classifications)
}

func (scc *SecrecyClassificationController) GetSecrecyClassificationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	classification, err := scc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, classification)
}

func (scc *SecrecyClassificationController) CreateSecrecyClassification(c *gin.Context) {
	var classification *models.SecrecyClassification

	if err := c.ShouldBindJSON(&classification); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}

	createdClassification, err := scc.service.Create(currentUserID, classification)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdClassification)
}

func (scc *SecrecyClassificationController) UpdateSecrecyClassification(c *gin.Context) {
	var updatedClassification *models.SecrecyClassification

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
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}

	updatedClassification, err = scc.service.Update(uint(classificationID), currentUserID, updatedClassification)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(scc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedClassification)
}

func (scc *SecrecyClassificationController) DeleteSecrecyClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}

	failedIDs, err := scc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (scc *SecrecyClassificationController) ActivateSecrecyClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}

	err = scc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (scc *SecrecyClassificationController) DeactivateSecrecyClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}

	err = scc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(scc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
