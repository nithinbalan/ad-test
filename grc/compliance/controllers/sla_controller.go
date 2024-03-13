package controllers

import (
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SLAController struct {
	service services.SLAServiceInterface
	db      *gorm.DB
}

func NewSLAController(service services.SLAServiceInterface) *SLAController {
	return &SLAController{
		service: service,
	}
}

func (src *SLAController) GetAllSLA(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	slaRegisters, err := src.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, slaRegisters)
}

func (src *SLAController) GetSLADetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sla, err := src.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, sla)
}

func (src *SLAController) CreateSLA(c *gin.Context) {

	var slaRegister *models.SLARegister

	if err := c.ShouldBindJSON(&slaRegister); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	createdSLA, err := src.service.Create(currentUserID, slaRegister)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(src.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdSLA)
}

func (src *SLAController) UpdateSLA(c *gin.Context) {

	slaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedSLA *models.SLARegister
	if err := c.ShouldBindJSON(&updatedSLA); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	updatedSLA, err = src.service.Update(uint(slaID), currentUserID, updatedSLA)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(src.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedSLA)
}

func (src *SLAController) DeleteSLA(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	failedIDs, err := src.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}


