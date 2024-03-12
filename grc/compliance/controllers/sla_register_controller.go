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

type SLARegisterController struct {
	service services.SLARegisterServiceInterface
	db      *gorm.DB
}

func NewSLARegisterController(service services.SLARegisterServiceInterface) *SLARegisterController {
	return &SLARegisterController{
		service: service,
	}
}

func (src *SLARegisterController) GetAllSLARegisters(c *gin.Context) {

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

func (src *SLARegisterController) GetSLARegisterDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	slaRegister, err := src.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, slaRegister)
}

func (src *SLARegisterController) CreateSLARegister(c *gin.Context) {

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

	createdRegister, err := src.service.Create(currentUserID, slaRegister)
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

	c.JSON(http.StatusCreated, createdRegister)
}

func (src *SLARegisterController) UpdateSLARegister(c *gin.Context) {

	slaRegisterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedRegister *models.SLARegister
	if err := c.ShouldBindJSON(&updatedRegister); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	updatedRegister, err = src.service.Update(uint(slaRegisterID), currentUserID, updatedRegister)
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

	c.JSON(http.StatusOK, updatedRegister)
}

func (src *SLARegisterController) DeleteSLARegisters(c *gin.Context) {

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


