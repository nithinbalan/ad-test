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
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	ratings, err := acmfrc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(acmfrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, ratings)
}

func (acmfrc *AssetCriticalityMatrixFactorRatingController) GetRatingDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	rating, err := acmfrc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, rating)
}

func (acmfrc *AssetCriticalityMatrixFactorRatingController) CreateRating(c *gin.Context) {
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

	createdRating, err := acmfrc.service.Create(currentUserID, rating)
	if err != nil {
		utility.SendInternalErrResp(acmfrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdRating)
}

func (acmfrc *AssetCriticalityMatrixFactorRatingController) UpdateRating(c *gin.Context) {
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

	updatedRating, err = acmfrc.service.Update(uint(ratingID), currentUserID, updatedRating)
	if err != nil {
		utility.SendInternalErrResp(acmfrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, updatedRating)
}

func (acmfrc *AssetCriticalityMatrixFactorRatingController) DeleteRatings(c *gin.Context) {
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

	failedIDs, err := acmfrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acmfrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
