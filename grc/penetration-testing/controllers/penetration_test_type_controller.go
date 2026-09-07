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

type PenetrationTestTypeController struct {
	service services.PenetrationTestTypeServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestTypeController(service services.PenetrationTestTypeServiceInterface) *PenetrationTestTypeController {
	return &PenetrationTestTypeController{
		service: service,
	}
}

func (pttc *PenetrationTestTypeController) GetAllPenetrationTestTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	penetrationTestTypes, err := pttc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(pttc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, penetrationTestTypes)
}

func (pttc *PenetrationTestTypeController) GetPenetrationTestTypeDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	penetrationTestType, err := pttc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, penetrationTestType)
}

func (pttc *PenetrationTestTypeController) CreatePenetrationTestType(c *gin.Context) {
	var penetrationTestType *models.PenetrationTestType

	if err := c.ShouldBindJSON(&penetrationTestType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pttc.db, c, err.Error())
		return
	}

	created, err := pttc.service.Create(currentUserID, penetrationTestType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pttc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, created)
}

func (pttc *PenetrationTestTypeController) UpdatePenetrationTestType(c *gin.Context) {
	var updated *models.PenetrationTestType

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
		utility.SendInternalErrResp(pttc.db, c, err.Error())
		return
	}

	updated, err = pttc.service.Update(uint(id), currentUserID, updated)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pttc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updated)
}

func (pttc *PenetrationTestTypeController) DeletePenetrationTestTypes(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pttc.db, c, err.Error())
		return
	}

	failedIDs, err := pttc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pttc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pttc *PenetrationTestTypeController) ActivatePenetrationTestTypes(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pttc.db, c, err.Error())
		return
	}

	if err := pttc.service.Activate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(pttc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (pttc *PenetrationTestTypeController) DeactivatePenetrationTestTypes(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pttc.db, c, err.Error())
		return
	}

	if err := pttc.service.Deactivate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(pttc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
