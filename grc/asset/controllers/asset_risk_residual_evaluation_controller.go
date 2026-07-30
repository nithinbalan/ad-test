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

type AssetRiskResidualEvaluationController struct {
	service services.AssetRiskResidualEvaluationServiceInterface
	db      *gorm.DB
}

func NewAssetRiskResidualEvaluationController(service services.AssetRiskResidualEvaluationServiceInterface) *AssetRiskResidualEvaluationController {
	return &AssetRiskResidualEvaluationController{
		service: service,
	}
}

func (arrc *AssetRiskResidualEvaluationController) GetAssetRiskResidualEvaluation(c *gin.Context) {
	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	evaluation, err := arrc.service.GetFullDetailsByID(uint(riskID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			c.JSON(http.StatusOK, nil)
			return
		default:
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, evaluation)
}

func (arrc *AssetRiskResidualEvaluationController) CreateAssetRiskResidualEvaluation(c *gin.Context) {
	riskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO assetDTOs.AssetRiskResidualEvaluationFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	id, err := arrc.service.Create(uint(riskID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (arrc *AssetRiskResidualEvaluationController) UpdateAssetRiskResidualEvaluation(c *gin.Context) {
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

	var formDTO assetDTOs.AssetRiskResidualEvaluationFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	if err := arrc.service.Update(uint(riskID), uint(id), currentUserID, &formDTO); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arrc *AssetRiskResidualEvaluationController) DeleteAssetRiskResidualEvaluation(c *gin.Context) {
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
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	failedID, err := arrc.service.Delete(uint(riskID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
