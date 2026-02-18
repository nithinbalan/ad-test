package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlConsistencyController struct {
	service services.ControlConsistencyServiceInterface
	db      *gorm.DB
}

func NewControlConsistencyController(service services.ControlConsistencyServiceInterface) *ControlConsistencyController {
	return &ControlConsistencyController{
		service: service,
	}
}

func (ccc *ControlConsistencyController) GetAllControlConsistencies(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	controlConsistencies, err := ccc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, controlConsistencies)
}

func (ccc *ControlConsistencyController) GetControlConsistencyDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlConsistency, err := ccc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, controlConsistency)
}

func (ccc *ControlConsistencyController) ActivateControlConsistencies(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	err = ccc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ccc *ControlConsistencyController) DeactivateControlConsistencies(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	err = ccc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}


