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

type PlannedAdvisoryInternalStatusController struct {
	service services.PlannedAdvisoryInternalStatusServiceInterface
	db      *gorm.DB
}

func NewPlannedAdvisoryInternalStatusController(service services.PlannedAdvisoryInternalStatusServiceInterface) *PlannedAdvisoryInternalStatusController {
	return &PlannedAdvisoryInternalStatusController{
		service: service,
	}
}

func (pais *PlannedAdvisoryInternalStatusController) GetAllPlannedAdvisoryInternalStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	internalStatuses, err := pais.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(pais.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, internalStatuses)
}

func (pais *PlannedAdvisoryInternalStatusController) GetPlannedAdvisoryInternalStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	plannedAdvisoryInternalStatus, err := pais.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, plannedAdvisoryInternalStatus)
}

func (pais *PlannedAdvisoryInternalStatusController) CreatePlannedAdvisoryInternalStatus(c *gin.Context) {

	var internalStatus *models.PlannedAdvisoryInternalStatus

	if err := c.ShouldBindJSON(&internalStatus); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pais.db, c, err.Error())
		return
	}

	createdInternalStatus, err := pais.service.Create(currentUserID, internalStatus)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pais.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdInternalStatus)
}

func (pais *PlannedAdvisoryInternalStatusController) UpdatePlannedAdvisoryInternalStatus(c *gin.Context) {

	internalStatusID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedInternalStatus *models.PlannedAdvisoryInternalStatus
	if err := c.ShouldBindJSON(&updatedInternalStatus); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pais.db, c, err.Error())
		return
	}

	updatedInternalStatus, err = pais.service.Update(uint(internalStatusID), currentUserID, updatedInternalStatus)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pais.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedInternalStatus)
}

func (pais *PlannedAdvisoryInternalStatusController) DeletePlannedAdvisoryInternalStatuses(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pais.db, c, err.Error())
		return
	}

	failedIDs, err := pais.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pais.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pais *PlannedAdvisoryInternalStatusController) ActivatePlannedAdvisoryInternalStatuses(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pais.db, c, err.Error())
		return
	}

	err = pais.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pais.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (pais *PlannedAdvisoryInternalStatusController) DeactivatePlannedAdvisoryInternalStatuses(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pais.db, c, err.Error())
		return
	}

	err = pais.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pais.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
