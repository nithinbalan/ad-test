package controllers

import (
	organizationDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"jethur/grc/asset/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskInherentControlAnalysisController struct {
	service services.AssetRiskInherentControlAnalysisServiceInterface
	db      *gorm.DB
}

func NewAssetRiskInherentControlAnalysisController(service services.AssetRiskInherentControlAnalysisServiceInterface) *AssetRiskInherentControlAnalysisController {
	return &AssetRiskInherentControlAnalysisController{
		service: service,
	}
}

func (arcac *AssetRiskInherentControlAnalysisController) GetAllAssetRiskControlAnalyses(c *gin.Context) {
	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := riskDTOs.RiskControlAnalysisListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	riskControlAnalyses, err := arcac.service.GetAll(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arcac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskControlAnalyses)
}

func (arcac *AssetRiskInherentControlAnalysisController) GetAssetRiskControlAnalysisDetails(c *gin.Context) {
	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("analysisId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskControlAnalysis, err := arcac.service.GetFullDetailsByID(uint(riskID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arcac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, riskControlAnalysis)
}

func (arcac *AssetRiskInherentControlAnalysisController) CreateAssetRiskControlAnalysis(c *gin.Context) {
	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arcac.db, c, err.Error())
		return
	}

	err = arcac.service.Create(uint(riskID), currentUserID)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arcac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arcac *AssetRiskInherentControlAnalysisController) UpdateAssetRiskControlAnalysis(c *gin.Context) {
	var controlTestFormDTO organizationDTOs.ControlTestFormDTO

	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("analysisId")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&controlTestFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arcac.db, c, err.Error())
		return
	}

	err = arcac.service.Update(uint(riskID), uint(id), currentUserID, &controlTestFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arcac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arcac *AssetRiskInherentControlAnalysisController) DeleteAssetRiskControlAnalyses(c *gin.Context) {
	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
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
		utility.SendInternalErrResp(arcac.db, c, err.Error())
		return
	}

	failedIDs, err := arcac.service.Delete(uint(riskID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(arcac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (arcac *AssetRiskInherentControlAnalysisController) GetAssetRiskControlSummary(c *gin.Context) {
	riskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	risk, err := arcac.service.GetControlSummary(uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(arcac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, risk)
}
