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

type ProcessGroupController struct {
	service services.ProcessGroupServiceInterface
	db      *gorm.DB
}

func NewProcessGroupController(service services.ProcessGroupServiceInterface) *ProcessGroupController {
	return &ProcessGroupController{
		service: service,
	}
}

func (pgc *ProcessGroupController) GetAllProcessGroups(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	processGroups, err := pgc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(pgc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, processGroups)
}

func (pgc *ProcessGroupController) GetProcessGroupDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	processGroup, err := pgc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, processGroup)
}

func (pgc *ProcessGroupController) CreateProcessGroup(c *gin.Context) {
	var processGroup *models.ProcessGroup

	if err := c.ShouldBindJSON(&processGroup); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pgc.db, c, err.Error())
		return
	}

	createdProcessGroup, err := pgc.service.Create(currentUserID, processGroup)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdProcessGroup)
}

func (pgc *ProcessGroupController) UpdateProcessGroup(c *gin.Context) {
	processGroupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedProcessGroup *models.ProcessGroup
	if err := c.ShouldBindJSON(&updatedProcessGroup); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pgc.db, c, err.Error())
		return
	}

	updatedProcessGroup, err = pgc.service.Update(uint(processGroupID), currentUserID, updatedProcessGroup)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedProcessGroup)
}

func (pgc *ProcessGroupController) DeleteProcessGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pgc.db, c, err.Error())
		return
	}

	failedIDs, err := pgc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pgc *ProcessGroupController) ActivateProcessGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pgc.db, c, err.Error())
		return
	}

	err = pgc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (pgc *ProcessGroupController) DeactivateProcessGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pgc.db, c, err.Error())
		return
	}

	err = pgc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
