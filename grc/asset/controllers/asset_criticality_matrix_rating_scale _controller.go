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

type AssetCriticalityMatrixRatingScaleController struct {
	service services.AssetCriticalityMatrixRatingScaleServiceInterface
	db      *gorm.DB
}

func NewAssetCriticalityMatrixRatingScaleController(service services.AssetCriticalityMatrixRatingScaleServiceInterface) *AssetCriticalityMatrixRatingScaleController {
	return &AssetCriticalityMatrixRatingScaleController{
		service: service,
	}
}

func (acmrs *AssetCriticalityMatrixRatingScaleController) GetAllScales(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	scales, err := acmrs.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(acmrs.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, scales)
}

func (acmrs *AssetCriticalityMatrixRatingScaleController) GetScaleDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	scale, err := acmrs.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, scale)
}

func (acmrs *AssetCriticalityMatrixRatingScaleController) CreateScale(c *gin.Context) {
	var scale *models.AssetCriticalityMatrixRatingScale

	if err := c.ShouldBindJSON(&scale); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmrs.db, c, err.Error())
		return
	}

	createdScale, err := acmrs.service.Create(currentUserID, scale)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acmrs.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdScale)
}

func (acmrs *AssetCriticalityMatrixRatingScaleController) UpdateScale(c *gin.Context) {
	var updatedScale *models.AssetCriticalityMatrixRatingScale
	scaleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedScale); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmrs.db, c, err.Error())
		return
	}

	updatedScale, err = acmrs.service.Update(uint(scaleID), currentUserID, updatedScale)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acmrs.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedScale)
}

func (acmrs *AssetCriticalityMatrixRatingScaleController) DeleteScales(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmrs.db, c, err.Error())
		return
	}

	failedIDs, err := acmrs.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acmrs.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
