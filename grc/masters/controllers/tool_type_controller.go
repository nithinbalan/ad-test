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

type ToolTypeController struct {
	service services.ToolTypeServiceInterface
	db      *gorm.DB
}

func NewToolTypeController(service services.ToolTypeServiceInterface) *ToolTypeController {
	return &ToolTypeController{
		service: service,
	}
}

func (ttc *ToolTypeController) GetAllToolTypes(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	toolTypes, err := ttc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, toolTypes)
}

func (ttc *ToolTypeController) GetToolTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	toolType, err := ttc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, toolType)
}

func (ttc *ToolTypeController) CreateToolType(c *gin.Context) {

	var toolType *models.ToolType

	if err := c.ShouldBindJSON(&toolType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}

	createdToolType, err := ttc.service.Create(currentUserID, toolType)
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

	c.JSON(http.StatusCreated, createdToolType)
}

func (ttc *ToolTypeController) UpdateToolType(c *gin.Context) {

	toolTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedToolType *models.ToolType
	if err := c.ShouldBindJSON(&updatedToolType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ttc.db, c, err.Error())
		return
	}

	updatedToolType, err = ttc.service.Update(uint(toolTypeID), currentUserID, updatedToolType)
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

	c.JSON(http.StatusOK, updatedToolType)
}

func (ttc *ToolTypeController) DeleteToolTypes(c *gin.Context) {

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

func (ttc *ToolTypeController) ActivateToolTypes(c *gin.Context) {
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

func (ttc *ToolTypeController) DeactivateToolTypes(c *gin.Context) {
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
