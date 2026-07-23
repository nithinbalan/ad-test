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

type AssetRiskInherentEvaluationController struct {
	service services.AssetRiskInherentEvaluationServiceInterface
	db      *gorm.DB
}

func NewAssetRiskInherentEvaluationController(service services.AssetRiskInherentEvaluationServiceInterface) *AssetRiskInherentEvaluationController {
	return &AssetRiskInherentEvaluationController{
		service: service,
	}
}

func (aric *AssetRiskInherentEvaluationController) GetAssetRiskInherentEvaluation(c *gin.Context) {
	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	evaluation, err := aric.service.GetFullDetailsByID(uint(riskID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aric.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, evaluation)
}

func (aric *AssetRiskInherentEvaluationController) CreateAssetRiskInherentEvaluation(c *gin.Context) {
	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO assetDTOs.AssetRiskInherentEvaluationFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aric.db, c, err.Error())
		return
	}

	id, err := aric.service.Create(uint(riskID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aric.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (aric *AssetRiskInherentEvaluationController) UpdateAssetRiskInherentEvaluation(c *gin.Context) {
	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("evaluationId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO assetDTOs.AssetRiskInherentEvaluationFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aric.db, c, err.Error())
		return
	}

	if err := aric.service.Update(uint(riskID), uint(id), currentUserID, &formDTO); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aric.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (aric *AssetRiskInherentEvaluationController) DeleteAssetRiskInherentEvaluation(c *gin.Context) {
	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("evaluationId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aric.db, c, err.Error())
		return
	}

	failedID, err := aric.service.Delete(uint(riskID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aric.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
