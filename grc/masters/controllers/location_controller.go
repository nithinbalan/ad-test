package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LocationController struct {
	service services.LocationServiceInterface
	db      *gorm.DB
}

func NewLocationController(service services.LocationServiceInterface) *LocationController {
	return &LocationController{
		service: service,
	}
}

func (lc *LocationController) GetAllLocations(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	locations, err := lc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(lc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, locations)
}

func (lc *LocationController) GetLocationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	location, err := lc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, location)
}

func (lc *LocationController) CreateLocation(c *gin.Context) {

	var location *models.Location

	if err := c.ShouldBindJSON(&location); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lc.db, c, err.Error())
		return
	}

	createdLocation, err := lc.service.Create(currentUserID, location)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(lc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdLocation.ID})
}

func (lc *LocationController) UpdateLocation(c *gin.Context) {

	locationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedLocation *models.Location
	if err := c.ShouldBindJSON(&updatedLocation); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lc.db, c, err.Error())
		return
	}

	updatedLocation, err = lc.service.Update(uint(locationID), currentUserID, updatedLocation)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(lc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedLocation)
}

func (lc *LocationController) DeleteLocations(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lc.db, c, err.Error())
		return
	}

	failedIDs, err := lc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(lc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (lc *LocationController) ActivateLocations(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lc.db, c, err.Error())
		return
	}

	err = lc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(lc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (lc *LocationController) DeactivateLocations(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lc.db, c, err.Error())
		return
	}

	err = lc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(lc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
