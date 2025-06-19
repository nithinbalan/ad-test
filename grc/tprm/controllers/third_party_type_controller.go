package controllers

import (
	"jethur/grc/models"
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyTypeController struct {
	service services.ThirdPartyTypeServiceInterface
	db      *gorm.DB
}

func NewThirdPartyTypeController(service services.ThirdPartyTypeServiceInterface) *ThirdPartyTypeController {
	return &ThirdPartyTypeController{
		service: service,
	}
}

func (tptc *ThirdPartyTypeController) GetAllThirdPartyTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	tprmTypes, err := tptc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tptc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, tprmTypes)
}

func (tptc *ThirdPartyTypeController) GetThirdPartyTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	tprmType, err := tptc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, tprmType)
}

func (tptc *ThirdPartyTypeController) CreateThirdPartyType(c *gin.Context) {
	var tprmType *models.ThirdPartyType

	if err := c.ShouldBindJSON(&tprmType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tptc.db, c, err.Error())
		return
	}

	createdType, err := tptc.service.Create(currentUserID, tprmType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tptc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdType)
}

func (tptc *ThirdPartyTypeController) UpdateThirdPartyType(c *gin.Context) {
	var updatedType *models.ThirdPartyType

	tprmTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tptc.db, c, err.Error())
		return
	}

	updatedType, err = tptc.service.Update(uint(tprmTypeID), currentUserID, updatedType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tptc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedType)
}

func (tptc *ThirdPartyTypeController) DeleteThirdPartyTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tptc.db, c, err.Error())
		return
	}

	failedIDs, err := tptc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tptc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tptc *ThirdPartyTypeController) ActivateThirdPartyTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tptc.db, c, err.Error())
		return
	}

	err = tptc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tptc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (tptc *ThirdPartyTypeController) DeactivateThirdPartyTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tptc.db, c, err.Error())
		return
	}

	err = tptc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tptc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
