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

type AssetTypeController struct {
	service services.AssetTypeServiceInterface
	db      *gorm.DB
}

func NewAssetTypeController(service services.AssetTypeServiceInterface) *AssetTypeController {
	return &AssetTypeController{
		service: service,
	}
}

func (atc *AssetTypeController) GetAllAssetTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	assetTypes, err := atc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(atc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, assetTypes)
}

func (atc *AssetTypeController) GetAssetTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetType, err := atc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, assetType)
}

func (atc *AssetTypeController) CreateAssetType(c *gin.Context) {
	var assetType *models.AssetType

	if err := c.ShouldBindJSON(&assetType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(atc.db, c, err.Error())
		return
	}

	createdType, err := atc.service.Create(currentUserID, assetType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(atc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdType)
}

func (atc *AssetTypeController) UpdateAssetType(c *gin.Context) {
	var updatedType *models.AssetType

	assetTypeID, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(atc.db, c, err.Error())
		return
	}

	updatedType, err = atc.service.Update(uint(assetTypeID), currentUserID, updatedType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(atc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedType)
}

func (atc *AssetTypeController) DeleteAssetTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(atc.db, c, err.Error())
		return
	}

	failedIDs, err := atc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(atc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (atc *AssetTypeController) ActivateAssetTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(atc.db, c, err.Error())
		return
	}

	err = atc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(atc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (atc *AssetTypeController) DeactivateAssetTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(atc.db, c, err.Error())
		return
	}

	err = atc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(atc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
