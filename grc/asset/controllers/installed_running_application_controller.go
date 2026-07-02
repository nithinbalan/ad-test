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

type InstalledRunningApplicationController struct {
	service services.InstalledRunningApplicationServiceInterface
	db      *gorm.DB
}

func NewInstalledRunningApplicationController(service services.InstalledRunningApplicationServiceInterface, db *gorm.DB) *InstalledRunningApplicationController {
	return &InstalledRunningApplicationController{
		service: service,
		db:      db,
	}
}

func (irac *InstalledRunningApplicationController) GetAllInstalledRunningApplications(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	installedRunningApplications, err := irac.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(irac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, installedRunningApplications)
}

func (irac *InstalledRunningApplicationController) GetInstalledRunningApplicationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	installedRunningApplication, err := irac.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, installedRunningApplication)
}

func (irac *InstalledRunningApplicationController) CreateInstalledRunningApplication(c *gin.Context) {
	var form dtos.InstalledRunningApplicationFormDTO

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irac.db, c, err.Error())
		return
	}

	createdApplication, err := irac.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(irac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdApplication)
}

func (irac *InstalledRunningApplicationController) UpdateInstalledRunningApplication(c *gin.Context) {
	var form dtos.InstalledRunningApplicationFormDTO

	installedRunningApplicationID, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(irac.db, c, err.Error())
		return
	}

	updatedApplication, err := irac.service.Update(uint(installedRunningApplicationID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(irac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedApplication)
}

func (irac *InstalledRunningApplicationController) DeleteInstalledRunningApplications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irac.db, c, err.Error())
		return
	}

	failedIDs, err := irac.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(irac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (irac *InstalledRunningApplicationController) ActivateInstalledRunningApplications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irac.db, c, err.Error())
		return
	}

	err = irac.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(irac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (irac *InstalledRunningApplicationController) DeactivateInstalledRunningApplications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irac.db, c, err.Error())
		return
	}

	err = irac.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(irac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}