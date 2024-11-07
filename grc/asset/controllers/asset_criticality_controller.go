package controllers

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetCriticalityController struct {
	service services.AssetCriticalityServiceInterface
	db      *gorm.DB
}

func NewAssetCriticalityController(service services.AssetCriticalityServiceInterface) *AssetCriticalityController {
	return &AssetCriticalityController{
		service: service,
	}
}

func (acc *AssetCriticalityController) GetAssetCriticality(c *gin.Context) {
	assetIDStr := c.Param("assetId")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetCriticality, err := acc.service.GetFullDetailsByID(uint(assetID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, assetCriticality)
}

func (acc *AssetCriticalityController) CreateAssetCriticality(c *gin.Context) {
	var assetCriticalityFormDTO assetDTOs.AssetCriticalityFormDTO

	assetIDStr := c.Param("assetId")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&assetCriticalityFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	id, err := acc.service.Create(uint(assetID), currentUserID, &assetCriticalityFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (acc *AssetCriticalityController) UpdateAssetCriticality(c *gin.Context) {
	var assetCriticalityFormDTO assetDTOs.AssetCriticalityFormDTO

	assetIDStr := c.Param("assetId")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&assetCriticalityFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	err = acc.service.Update(uint(assetID), uint(id), currentUserID, &assetCriticalityFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (acc *AssetCriticalityController) DeleteAssetCriticality(c *gin.Context) {
	assetIDStr := c.Param("assetId")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	failedID, err := acc.service.Delete(uint(assetID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
