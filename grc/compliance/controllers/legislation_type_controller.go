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

type LegislationTypeController struct {
	service services.LegislationTypeServiceInterface
	db      *gorm.DB
}

func NewLegislationTypeController(service services.LegislationTypeServiceInterface) *LegislationTypeController {
	return &LegislationTypeController{
		service: service,
	}
}

// Get all Legislation Types
func (ltc *LegislationTypeController) GetAllLegislationTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	types, err := ltc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(ltc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, types)
}

// Get Legislation Type details by ID
func (ltc *LegislationTypeController) GetLegislationTypeDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	legType, err := ltc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, legType)
}

// Create new Legislation Type
func (ltc *LegislationTypeController) CreateLegislationType(c *gin.Context) {
	var legType *models.ComplianceLegislationType
	if err := c.ShouldBindJSON(&legType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ltc.db, c, err.Error())
		return
	}

	createdType, err := ltc.service.Create(currentUserID, legType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(ltc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, createdType)
}

// Update Legislation Type
func (ltc *LegislationTypeController) UpdateLegislationType(c *gin.Context) {
	var legType *models.ComplianceLegislationType
	typeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&legType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ltc.db, c, err.Error())
		return
	}

	updatedType, err := ltc.service.Update(uint(typeID), currentUserID, legType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(ltc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, updatedType)
}

// Delete Legislation Types (bulk)
func (ltc *LegislationTypeController) DeleteLegislationTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ltc.db, c, err.Error())
		return
	}

	failedIDs, err := ltc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ltc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

// Activate Legislation Types (bulk)
func (ltc *LegislationTypeController) ActivateLegislationTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ltc.db, c, err.Error())
		return
	}

	if err := ltc.service.Activate(request.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(ltc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// Deactivate Legislation Types (bulk)
func (ltc *LegislationTypeController) DeactivateLegislationTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ltc.db, c, err.Error())
		return
	}

	if err := ltc.service.Deactivate(request.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(ltc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
