package controllers

import (
	"jethur/grc/models"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestEnvironmentReadinessController struct {
	service services.PenetrationTestEnvironmentReadinessServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestEnvironmentReadinessController(service services.PenetrationTestEnvironmentReadinessServiceInterface) *PenetrationTestEnvironmentReadinessController {
	return &PenetrationTestEnvironmentReadinessController{
		service: service,
	}
}

func (pterc *PenetrationTestEnvironmentReadinessController) GetAllPenetrationTestEnvironmentReadinesses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	penetrationTestEnvironmentReadinesses, err := pterc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(pterc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, penetrationTestEnvironmentReadinesses)
}

func (pterc *PenetrationTestEnvironmentReadinessController) GetPenetrationTestEnvironmentReadinessDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	penetrationTestEnvironmentReadiness, err := pterc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, penetrationTestEnvironmentReadiness)
}

func (pterc *PenetrationTestEnvironmentReadinessController) CreatePenetrationTestEnvironmentReadiness(c *gin.Context) {
	var penetrationTestEnvironmentReadiness *models.PenetrationTestEnvironmentReadiness

	if err := c.ShouldBindJSON(&penetrationTestEnvironmentReadiness); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pterc.db, c, err.Error())
		return
	}

	created, err := pterc.service.Create(currentUserID, penetrationTestEnvironmentReadiness)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pterc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, created)
}

func (pterc *PenetrationTestEnvironmentReadinessController) UpdatePenetrationTestEnvironmentReadiness(c *gin.Context) {
	var updated *models.PenetrationTestEnvironmentReadiness

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updated); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pterc.db, c, err.Error())
		return
	}

	updated, err = pterc.service.Update(uint(id), currentUserID, updated)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pterc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updated)
}

func (pterc *PenetrationTestEnvironmentReadinessController) DeletePenetrationTestEnvironmentReadinesses(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pterc.db, c, err.Error())
		return
	}

	failedIDs, err := pterc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pterc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pterc *PenetrationTestEnvironmentReadinessController) ActivatePenetrationTestEnvironmentReadinesses(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pterc.db, c, err.Error())
		return
	}

	if err := pterc.service.Activate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(pterc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (pterc *PenetrationTestEnvironmentReadinessController) DeactivatePenetrationTestEnvironmentReadinesses(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pterc.db, c, err.Error())
		return
	}

	if err := pterc.service.Deactivate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(pterc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
