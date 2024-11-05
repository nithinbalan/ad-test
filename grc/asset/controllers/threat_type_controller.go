package controllers

import (
	"jethur/grc/asset/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThreatTypeController struct {
	service services.ThreatTypeServiceInterface
	db      *gorm.DB
}

func NewThreatTypeController(service services.ThreatTypeServiceInterface) *ThreatTypeController {
	return &ThreatTypeController{
		service: service,
	}
}

func (ttc *ThreatTypeController) GetAllThreatTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	threatTypes, err := ttc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, threatTypes)
}

func (ttc *ThreatTypeController) GetThreatTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	threatType, err := ttc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, threatType)
}

func (ttc *ThreatTypeController) CreateThreatType(c *gin.Context) {
	var threatType *models.ThreatType

	if err := c.ShouldBindJSON(&threatType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}

	createdType, err := ttc.service.Create(currentUserID, threatType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ttc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdType)
}

func (ttc *ThreatTypeController) UpdateThreatType(c *gin.Context) {
	var updatedType *models.ThreatType

	threatTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}

	updatedType, err = ttc.service.Update(uint(threatTypeID), currentUserID, updatedType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ttc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedType)
}

func (ttc *ThreatTypeController) DeleteThreatTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}

	failedIDs, err := ttc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ttc *ThreatTypeController) ActivateThreatTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}

	err = ttc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ttc *ThreatTypeController) DeactivateThreatTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}

	err = ttc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
