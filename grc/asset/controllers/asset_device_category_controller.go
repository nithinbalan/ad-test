package controllers

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetDeviceCategoryController struct {
	service services.AssetDeviceCategoryServiceInterface
	db      *gorm.DB
}

func NewAssetDeviceCategoryController(service services.AssetDeviceCategoryServiceInterface) *AssetDeviceCategoryController {
	return &AssetDeviceCategoryController{
		service: service,
	}
}

func (adcc *AssetDeviceCategoryController) GetAllAssetDeviceCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	assetDeviceCategories, err := adcc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, assetDeviceCategories)
}

func (adcc *AssetDeviceCategoryController) GetAssetDeviceCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetDeviceCategory, err := adcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, assetDeviceCategory)
}

func (adcc *AssetDeviceCategoryController) CreateAssetDeviceCategory(c *gin.Context) {
	var form dtos.AssetDeviceCategoryFormDTO

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	createdCategory, err := adcc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (adcc *AssetDeviceCategoryController) UpdateAssetDeviceCategory(c *gin.Context) {
	var form dtos.AssetDeviceCategoryFormDTO

	assetDeviceCategoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	updatedCategory, err := adcc.service.Update(uint(assetDeviceCategoryID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (adcc *AssetDeviceCategoryController) DeleteAssetDeviceCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	failedIDs, err := adcc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (adcc *AssetDeviceCategoryController) ActivateAssetDeviceCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	err = adcc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (adcc *AssetDeviceCategoryController) DeactivateAssetDeviceCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	err = adcc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}