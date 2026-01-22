package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LegislativeBodyController struct {
	service services.LegislativeBodyServiceInterface
	db      *gorm.DB
}

func NewLegislativeBodyController(service services.LegislativeBodyServiceInterface) *LegislativeBodyController {
	return &LegislativeBodyController{
		service: service,
	}
}

// Get all Legislative Bodies
func (lbc *LegislativeBodyController) GetAllLegislativeBodies(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	bodies, err := lbc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(lbc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, bodies)
}

// Get Legislative Body details by ID
func (lbc *LegislativeBodyController) GetLegislativeBodyDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	body, err := lbc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, body)
}

// Create new Legislative Body
func (lbc *LegislativeBodyController) CreateLegislativeBody(c *gin.Context) {
	var body *models.LegislativeBody
	if err := c.ShouldBindJSON(&body); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lbc.db, c, err.Error())
		return
	}

	createdBody, err := lbc.service.Create(currentUserID, body)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(lbc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, createdBody)
}

// Update Legislative Body
func (lbc *LegislativeBodyController) UpdateLegislativeBody(c *gin.Context) {
	var body *models.LegislativeBody
	bodyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lbc.db, c, err.Error())
		return
	}

	updatedBody, err := lbc.service.Update(uint(bodyID), currentUserID, body)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(lbc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, updatedBody)
}

// Delete Legislative Bodies (bulk)
func (lbc *LegislativeBodyController) DeleteLegislativeBodies(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lbc.db, c, err.Error())
		return
	}

	failedIDs, err := lbc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(lbc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

// Activate Legislative Bodies (bulk)
func (lbc *LegislativeBodyController) ActivateLegislativeBodies(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lbc.db, c, err.Error())
		return
	}

	if err := lbc.service.Activate(request.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(lbc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// Deactivate Legislative Bodies (bulk)
func (lbc *LegislativeBodyController) DeactivateLegislativeBodies(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lbc.db, c, err.Error())
		return
	}

	if err := lbc.service.Deactivate(request.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(lbc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
