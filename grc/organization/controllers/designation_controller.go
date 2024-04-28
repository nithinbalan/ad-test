package controllers

import (
	"jethur/grc/models"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DesignationController struct {
	service services.DesignationServiceInterface
	db      *gorm.DB
}

func NewDesignationController(service services.DesignationServiceInterface) *DesignationController {
	return &DesignationController{
		service: service,
	}
}

func (dc *DesignationController) GetAllDesignations(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	designations, err := dc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, designations)
}

func (dc *DesignationController) GetDesignationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	designation, err := dc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, designation)
}

func (dc *DesignationController) CreateDesignation(c *gin.Context) {
	var designation *models.Designation

	if err := c.ShouldBindJSON(&designation); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	createdDesignation, err := dc.service.Create(currentUserID, designation)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdDesignation)
}

func (dc *DesignationController) UpdateDesignation(c *gin.Context) {
	designationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedDesignation *models.Designation
	if err := c.ShouldBindJSON(&updatedDesignation); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	updatedDesignation, err = dc.service.Update(uint(designationID), currentUserID, updatedDesignation)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedDesignation)
}

func (dc *DesignationController) DeleteDesignations(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	failedIDs, err := dc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (dc *DesignationController) ActivateDesignations(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	err = dc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (dc *DesignationController) DeactivateDesignations(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	err = dc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
