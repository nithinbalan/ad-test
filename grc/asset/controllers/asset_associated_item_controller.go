package controllers

import (
	assetDTOs "jethur/grc/asset/dtos"
	assetService "jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetAssociatedItemController struct {
	service assetService.AssetAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewAssetAssociatedItemController(service assetService.AssetAssociatedItemServiceInterface) *AssetAssociatedItemController {
	return &AssetAssociatedItemController{
		service: service,
	}
}

func (aaic *AssetAssociatedItemController) GetAllAssetControls(c *gin.Context) {

	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := assetDTOs.AssetAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	controls, err := aaic.service.GetAllControls(uint(assetID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (aaic *AssetAssociatedItemController) GetAllAssetIncidents(c *gin.Context) {

	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := assetDTOs.AssetAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	controls, err := aaic.service.GetAllIncidents(uint(assetID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (aaic *AssetAssociatedItemController) GetAllAssetRisks(c *gin.Context) {

	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(aaic.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := assetDTOs.AssetAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Value:       value,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	controls, err := aaic.service.GetAllRisks(uint(assetID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}

func (aaic *AssetAssociatedItemController) GetAllAssetVulnerabilities(c *gin.Context) {

	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(aaic.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := assetDTOs.AssetAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Value:       value,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}

	// Update method to GetAllVulnerabilities
	vulnerabilities, err := aaic.service.GetAllVulnerabilities(uint(assetID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, vulnerabilities)
}

func (aaic *AssetAssociatedItemController) GetAllAssetClauses(c *gin.Context) {

	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(aaic.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := assetDTOs.AssetAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Value:       value,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}

	// Update method to GetAllVulnerabilities
	vulnerabilities, err := aaic.service.GetAllClauses(uint(assetID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, vulnerabilities)
}

func (aaic *AssetAssociatedItemController) CreateAssetVulnerabilities(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var vulnerabilities assetDTOs.AssetVulnerabilitiesFormDTO
	if err := c.ShouldBindJSON(&vulnerabilities); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	err = aaic.service.CreateAssetVulnerabilities(uint(id), vulnerabilities, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (aaic *AssetAssociatedItemController) DeleteAssetVulnerabilities(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	failedIDs, err := aaic.service.DeleteAssetVulnerabilities(uint(id), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (aaic *AssetAssociatedItemController) CreateAssetComplianceClauses(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form assetDTOs.AssetComplianceClauseFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	err = aaic.service.CreateAssetComplianceClauses(uint(id), form, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (aaic *AssetAssociatedItemController) DeleteAssetComplianceClauses(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	failedIDs, err := aaic.service.DeleteAssetComplianceClauses(uint(id), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (aaic *AssetAssociatedItemController) CreateAssetThreats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var threats assetDTOs.AssetVulnerabilityThreatsFormDTO
	if err := c.ShouldBindJSON(&threats); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	err = aaic.service.CreateAssetVulnerabilityThreats(uint(id), threats, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (aaic *AssetAssociatedItemController) DeleteAssetThreats(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	failedIDs, err := aaic.service.DeleteAssetVulnerabilityThreats(uint(id), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (aaic *AssetAssociatedItemController) CreateAssetThreatRisks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var risks assetDTOs.AssetThreatRisksFormDTO
	if err := c.ShouldBindJSON(&risks); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	err = aaic.service.CreateAssetThreatRisks(uint(id), risks, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (aaic *AssetAssociatedItemController) DeleteAssetThreatRisks(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	failedIDs, err := aaic.service.DeleteAssetThreatRisks(uint(id), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (aaic *AssetAssociatedItemController) GetRiskAnalysis(c *gin.Context) {

	assetIDStr := c.Param("id")

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	results, err := aaic.service.GetRiskAnalysis(uint(assetID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, results)
}
