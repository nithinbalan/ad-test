package controllers

import (
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskMatrixController struct {
	service services.AssetRiskMatrixServiceInterface
	db      *gorm.DB
}

func NewAssetRiskMatrixController(service services.AssetRiskMatrixServiceInterface) *AssetRiskMatrixController {
	return &AssetRiskMatrixController{
		service: service,
	}
}

func (armc *AssetRiskMatrixController) GetAllAssetRiskMatrixImpacts(c *gin.Context) {
	impacts, err := armc.service.GetAllImpacts()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(armc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, impacts)
}

func (armc *AssetRiskMatrixController) GetAllAssetRiskMatrixLikelihoods(c *gin.Context) {
	likelihoods, err := armc.service.GetAllLikelihoods()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(armc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, likelihoods)
}

func (armc *AssetRiskMatrixController) GetAllAssetRiskMatrixRiskRatings(c *gin.Context) {
	riskRatings, err := armc.service.GetAllRiskRatings()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(armc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskRatings)
}

func (armc *AssetRiskMatrixController) GetAllAssetRiskMatrixRiskRatingValues(c *gin.Context) {
	riskRatingValues, err := armc.service.GetAllRiskRatingValues()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(armc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskRatingValues)
}

func (armc *AssetRiskMatrixController) GetAssetRiskMatrixConfiguration(c *gin.Context) {
	riskMatrix, err := armc.service.GetRiskMatrix()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(armc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskMatrix)
}
