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

type AssetCriticalityMatrixFactorRatingController struct {
	service services.AssetCriticalityMatrixFactorRatingServiceInterface
	db      *gorm.DB
}

func NewAssetCriticalityMatrixFactorRatingController(service services.AssetCriticalityMatrixFactorRatingServiceInterface) *AssetCriticalityMatrixFactorRatingController {
	return &AssetCriticalityMatrixFactorRatingController{
		service: service,
	}
}

func (acmfrc *AssetCriticalityMatrixFactorRatingController) GetAllRatings(c *gin.Context) {
	matrixIDStr := c.Param("matrixId")
	matrixID, err := strconv.ParseUint(matrixIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	ratings, err := acmfrc.service.GetAll(uint(matrixID), sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(acmfrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, ratings)
}

func (acmfrc *AssetCriticalityMatrixFactorRatingController) GetRatingDetails(c *gin.Context) {
	matrixIDStr := c.Param("matrixId")
	matrixID, err := strconv.ParseUint(matrixIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	rating, err := acmfrc.service.GetByID(uint(matrixID), uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, rating)
}

func (acmfrc *AssetCriticalityMatrixFactorRatingController) CreateRating(c *gin.Context) {
	matrixIDStr := c.Param("matrixId")
	matrixID, err := strconv.ParseUint(matrixIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var rating *models.AssetCriticalityMatrixFactorRating

	if err := c.ShouldBindJSON(&rating); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmfrc.db, c, err.Error())
		return
	}

	createdRating, err := acmfrc.service.Create(uint(matrixID), currentUserID, rating)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acmfrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdRating)
}

func (acmfrc *AssetCriticalityMatrixFactorRatingController) UpdateRating(c *gin.Context) {
	matrixIDStr := c.Param("matrixId")
	matrixID, err := strconv.ParseUint(matrixIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedRating *models.AssetCriticalityMatrixFactorRating
	ratingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedRating); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmfrc.db, c, err.Error())
		return
	}

	updatedRating, err = acmfrc.service.Update(uint(matrixID), uint(ratingID), currentUserID, updatedRating)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acmfrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedRating)
}

func (acmfrc *AssetCriticalityMatrixFactorRatingController) DeleteRatings(c *gin.Context) {
	matrixIDStr := c.Param("matrixId")
	matrixID, err := strconv.ParseUint(matrixIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acmfrc.db, c, err.Error())
		return
	}

	failedIDs, err := acmfrc.service.Delete(uint(matrixID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acmfrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
