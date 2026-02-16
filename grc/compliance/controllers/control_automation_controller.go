package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlAutomationController struct {
	service services.ControlAutomationServiceInterface
	db      *gorm.DB
}

func NewControlAutomationController(service services.ControlAutomationServiceInterface) *ControlAutomationController {
	return &ControlAutomationController{
		service: service,
	}
}

func (cac *ControlAutomationController) GetAllControlAutomations(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	controlAutomations, err := cac.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, controlAutomations)
}

func (cac *ControlAutomationController) GetControlAutomationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlAutomation, err := cac.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, controlAutomation)
}

func (cac *ControlAutomationController) ActivateControlAutomations(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	err = cac.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (cac *ControlAutomationController) DeactivateControlAutomations(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	err = cac.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}