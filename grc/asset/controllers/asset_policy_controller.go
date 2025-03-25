package controllers

import (
	"errors"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetPolicyController struct {
	service services.AssetPolicyServiceInterface
	db      *gorm.DB
}

func NewAssetPolicyController(service services.AssetPolicyServiceInterface) *AssetPolicyController {
	return &AssetPolicyController{service: service}
}

// GetAllPolicies - Retrieves all asset policies with filtering, sorting, and pagination
func (apc *AssetPolicyController) GetAllPolicies(c *gin.Context) {
	assetRegisterIDStr := c.Param("id")
	assetRegisterID, err := strconv.ParseUint(assetRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Query parameters
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page parameter"))
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page_size parameter"))
		return
	}

	// Valid sort fields
	validSortFields := map[string]string{
		"created_at":           "asset_policies.created_at",
		"document_title_en":    "documents.title_en",
		"document_title_ar":    "documents.title_ar",
		"compliance_status_en": "compliance_statuses.title_en",
		"compliance_status_ar": "compliance_statuses.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	// Pass query params to DTO
	params := dtos.AssetPolicyListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	policies, err := apc.service.GetAll(&params, uint(assetRegisterID))
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, policies)
}

// GetPolicyDetails - Retrieves details of a specific asset policy
func (apc *AssetPolicyController) GetPolicyDetails(c *gin.Context) {
	assetRegisterIDStr := c.Param("id")
	assetRegisterID, err := strconv.ParseUint(assetRegisterIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	idStr := c.Param("policyId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	policy, err := apc.service.GetFullDetailsByID(uint(assetRegisterID), uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, policy)
}

// CreatePolicy - Creates a new asset policy
func (apc *AssetPolicyController) CreatePolicy(c *gin.Context) {
	var policyForm dtos.AssetPolicyFormDTO
	if err := c.ShouldBindJSON(&policyForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	assetIdStr := c.Param("id")
	assetId, err := strconv.ParseUint(assetIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	policyIDs, err := apc.service.Create(&policyForm, currentUserID, uint(assetId))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ids": policyIDs})
}

// DeletePolicies - Deletes multiple asset policies
func (apc *AssetPolicyController) DeletePolicies(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	failedIDs, err := apc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (apc *AssetPolicyController) GetAssetPolicyComplianceSummary(c *gin.Context) {
	assetRegisterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := apc.service.GetAssetPolicyComplianceSummary(uint(assetRegisterID))
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, summary)
}
