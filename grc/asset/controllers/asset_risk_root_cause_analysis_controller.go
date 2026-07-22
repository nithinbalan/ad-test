package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	riskServices "jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskRootCauseAnalysisController struct {
	service riskServices.RiskRootCauseAnalysisServiceInterface
	db      *gorm.DB
}

func NewAssetRiskRootCauseAnalysisController(service riskServices.RiskRootCauseAnalysisServiceInterface) *AssetRiskRootCauseAnalysisController {
	return &AssetRiskRootCauseAnalysisController{
		service: service,
	}
}

func (arrcac *AssetRiskRootCauseAnalysisController) GetAllAssetRiskRootCauseAnalyses(c *gin.Context) {
	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := riskDTOs.RiskRootCauseAnalysisListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	riskRootCauseAnalyses, err := arrcac.service.GetAll(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arrcac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskRootCauseAnalyses)
}

func (arrcac *AssetRiskRootCauseAnalysisController) CreateAssetRiskRootCauseAnalysis(c *gin.Context) {
	var riskRootCauseAnalysisFormDTO riskDTOs.RiskRootCauseAnalysisFormDTO

	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&riskRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrcac.db, c, err.Error())
		return
	}

	id, err := arrcac.service.Create(uint(riskID), currentUserID, &riskRootCauseAnalysisFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arrcac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (arrcac *AssetRiskRootCauseAnalysisController) CreateAssetRiskRootCauseAnalyses(c *gin.Context) {
	var rcaFormDTO []riskDTOs.RiskRootCauseAnalysisFormDTO

	riskIDStr := c.Param("id")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&rcaFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrcac.db, c, err.Error())
		return
	}

	err = arrcac.service.CreateMultiple(uint(riskID), currentUserID, rcaFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(arrcac.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusCreated)
}

func (arrcac *AssetRiskRootCauseAnalysisController) UpdateAssetRiskRootCauseAnalysis(c *gin.Context) {
	var riskRootCauseAnalysisFormDTO riskDTOs.RiskRootCauseAnalysisFormDTO

	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("rcaId")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&riskRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrcac.db, c, err.Error())
		return
	}

	err = arrcac.service.Update(uint(riskID), uint(id), currentUserID, &riskRootCauseAnalysisFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arrcac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arrcac *AssetRiskRootCauseAnalysisController) DeleteAssetRiskRootCauseAnalysis(c *gin.Context) {
	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("rcaId")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrcac.db, c, err.Error())
		return
	}

	failedID, err := arrcac.service.Delete(uint(riskID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arrcac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
