package controllers

import (
	assetDTOs "jethur/grc/asset/dtos"

	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetControlAnalysisController struct {
	service services.AssetControlAnalysisServiceInterface
	db      *gorm.DB
}

func NewAssetControlAnalysisController(service services.AssetControlAnalysisServiceInterface) *AssetControlAnalysisController {
	return &AssetControlAnalysisController{
		service: service,
	}
}

func (acac *AssetControlAnalysisController) GetAllAssetControlAnalyses(c *gin.Context) {
	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := assetDTOs.AssetControlAnalysisListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	assetControlAnalyses, err := acac.service.GetAll(uint(assetID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(acac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, assetControlAnalyses)
}

func (acac *AssetControlAnalysisController) GetAssetControlAnalysisDetails(c *gin.Context) {
	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
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

	assetControlAnalysis, err := acac.service.GetFullDetailsByID(uint(assetID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(acac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, assetControlAnalysis)
}

func (acac *AssetControlAnalysisController) CreateAssetControlAnalysis(c *gin.Context) {
	var assetControlFormDTO assetDTOs.AssetControlFormDTO

	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&assetControlFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acac.db, c, err.Error())
		return
	}

	err = acac.service.Create(uint(assetID), currentUserID, &assetControlFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceConflictError:
			utility.SendConflictErrResp(c, e.Error())
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (acac *AssetControlAnalysisController) UpdateAssetControlAnalysis(c *gin.Context) {
	var assetControlUpdateFormDTO assetDTOs.AssetControlUpdateFormDTO

	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
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

	if err := c.ShouldBindJSON(&assetControlUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acac.db, c, err.Error())
		return
	}

	err = acac.service.Update(uint(assetID), uint(id), currentUserID, &assetControlUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (acac *AssetControlAnalysisController) DeleteAssetControlAnalyses(c *gin.Context) {
	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
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
		utility.SendInternalErrResp(acac.db, c, err.Error())
		return
	}

	failedIDs, err := acac.service.Delete(uint(assetID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

// func (acac *AssetControlAnalysisController) GetAssetControlSummary(c *gin.Context) {
// 	assetID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		utility.SendNotFoundErrResp(c)
// 		return
// 	}

// 	asset, err := acac.service.GetControlSummary(uint(assetID))
// 	if err != nil {
// 		utility.SendInternalErrResp(acac.db, c, err.Error())
// 		return
// 	}
// 	c.JSON(http.StatusOK, asset)
// }
