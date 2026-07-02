package controllers

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InstalledApplicationController struct {
	service services.InstalledApplicationServiceInterface
	db      *gorm.DB
}

func NewInstalledApplicationController(service services.InstalledApplicationServiceInterface, db *gorm.DB) *InstalledApplicationController {
	return &InstalledApplicationController{
		service: service,
		db:      db,
	}
}

func (iac *InstalledApplicationController) GetAllInstalledApplications(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	installedApplications, err := iac.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, installedApplications)
}

func (iac *InstalledApplicationController) GetInstalledApplicationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	installedApplication, err := iac.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, installedApplication)
}

func (iac *InstalledApplicationController) CreateInstalledApplication(c *gin.Context) {
	var form dtos.InstalledApplicationFormDTO

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	createdApplication, err := iac.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdApplication)
}

func (iac *InstalledApplicationController) UpdateInstalledApplication(c *gin.Context) {
	var form dtos.InstalledApplicationFormDTO

	installedApplicationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	updatedApplication, err := iac.service.Update(uint(installedApplicationID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedApplication)
}

func (iac *InstalledApplicationController) DeleteInstalledApplications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	failedIDs, err := iac.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (iac *InstalledApplicationController) ActivateInstalledApplications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	err = iac.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (iac *InstalledApplicationController) DeactivateInstalledApplications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	err = iac.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}