package controllers

import (
	"jethur/grc/asset/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetCriticalityMatrixFactorController struct {
	service services.AssetCriticalityMatrixFactorServiceInterface
	db      *gorm.DB
}

func NewAssetCriticalityMatrixFactorController(service services.AssetCriticalityMatrixFactorServiceInterface) *AssetCriticalityMatrixFactorController {
	return &AssetCriticalityMatrixFactorController{
		service: service,
	}
}

func (acmfc *AssetCriticalityMatrixFactorController) GetAllFactors(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	factors, err := acmfc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(acmfc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, factors)
}

func (acmfc *AssetCriticalityMatrixFactorController) GetFactorDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	factor, err := acmfc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, factor)
}

func (acmfc *AssetCriticalityMatrixFactorController) CreateFactor(c *gin.Context) {
	var factor *models.AssetCriticalityMatrixFactor

	if err := c.ShouldBindJSON(&factor); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmfc.db, c, err.Error())
		return
	}

	createdFactor, err := acmfc.service.Create(currentUserID, factor)
	if err != nil {
		utility.SendInternalErrResp(acmfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdFactor)
}

func (acmfc *AssetCriticalityMatrixFactorController) UpdateFactor(c *gin.Context) {
	var updatedFactor *models.AssetCriticalityMatrixFactor
	factorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedFactor); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmfc.db, c, err.Error())
		return
	}

	updatedFactor, err = acmfc.service.Update(uint(factorID), currentUserID, updatedFactor)
	if err != nil {
		utility.SendInternalErrResp(acmfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, updatedFactor)
}

func (acmfc *AssetCriticalityMatrixFactorController) DeleteFactors(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmfc.db, c, err.Error())
		return
	}

	failedIDs, err := acmfc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acmfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
