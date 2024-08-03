package controllers

import (
	"jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UnitController struct {
	service services.UnitServiceInterface
	db      *gorm.DB
}

func NewUnitController(service services.UnitServiceInterface) *UnitController {
	return &UnitController{
		service: service,
	}
}

func (uc *UnitController) GetAllUnits(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	units, err := uc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, units)
}

func (uc *UnitController) GetUnitDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	unit, err := uc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, unit)
}

func (uc *UnitController) CreateUnit(c *gin.Context) {

	var unit *models.Unit

	if err := c.ShouldBindJSON(&unit); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	createdUnit, err := uc.service.Create(currentUserID, unit)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdUnit)
}

func (uc *UnitController) UpdateUnit(c *gin.Context) {

	unitID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedUnit *models.Unit
	if err := c.ShouldBindJSON(&updatedUnit); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	updatedUnit, err = uc.service.Update(uint(unitID), currentUserID, updatedUnit)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedUnit)
}

func (uc *UnitController) DeleteUnits(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	failedIDs, err := uc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (uc *UnitController) ActivateUnits(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	err = uc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (uc *UnitController) DeactivateUnits(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	err = uc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
