package controllers

import (
	"errors"
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskActionPlanController struct {
	service services.AssetRiskActionPlanServiceInterface
	db      *gorm.DB
}

func NewAssetRiskActionPlanController(service services.AssetRiskActionPlanServiceInterface) *AssetRiskActionPlanController {
	return &AssetRiskActionPlanController{
		service: service,
	}
}

func (arapc *AssetRiskActionPlanController) GetAllAssetRiskActionPlans(c *gin.Context) {
	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}

	isPassiveUser := utility.IsPassiveAccess(c)

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Convert the export parameter to a boolean
	exportFromContext, exists := c.Get("export")
	var export bool

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(arapc.db, c, "Invalid export parameter")
			return
		}
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
		"responsible_user_ids",
		"action_plan_status_ids",
		"risk_register_ids",
	}

	validSortFields := map[string]string{
		"title":               "risk_action_plans.title",
		"risk_title":          "risk_register.title",
		"severity_level_en":   "severity_levels.title_en",
		"severity_level_ar":   "severity_levels.title_ar",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"target_date":         "risk_action_plans.target_date",
		"percentage":          "risk_action_plans.percentage",
		"status_en":           "action_plan_statuses.title_en",
		"status_ar":           "action_plan_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(arapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(arapc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arapc.db, c, err.Error())
		return
	}

	queryParams := assetDTOs.AssetRiskActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         dbField,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
		Export:         export,
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		IsPassiveUser:  isPassiveUser,
	}
	riskActionPlans, err := arapc.service.GetAll(uint(riskID), currentUserID, &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arapc.db, c, e.Error())
			return
		}
	}

	if export {
		excelFile, err := arapc.service.ExportAssetRiskActionPlan(riskActionPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(arapc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Risk action plan.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, riskActionPlans)
}

func (arapc *AssetRiskActionPlanController) GetAssetRiskActionPlanDetails(c *gin.Context) {
	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arapc.db, c, err.Error())
		return
	}

	isPassiveUser := utility.IsPassiveAccess(c)

	riskActionPlan, err := arapc.service.GetFullDetailsByID(uint(riskID), uint(id), currentUserID, isPassiveUser)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(arapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, riskActionPlan)
}

func (arapc *AssetRiskActionPlanController) CreateAssetRiskActionPlan(c *gin.Context) {
	var formDTO assetDTOs.AssetRiskActionPlanFormDTO

	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arapc.db, c, err.Error())
		return
	}

	id, err := arapc.service.Create(uint(riskID), currentUserID, &formDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (arapc *AssetRiskActionPlanController) UpdateAssetRiskActionPlan(c *gin.Context) {
	var formDTO assetDTOs.AssetRiskActionPlanFormDTO

	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arapc.db, c, err.Error())
		return
	}

	err = arapc.service.Update(uint(riskID), uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arapc *AssetRiskActionPlanController) DeleteAssetRiskActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	riskIDStr := c.Param("id")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arapc.db, c, err.Error())
		return
	}

	failedIDs, err := arapc.service.Delete(uint(riskID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (arapc *AssetRiskActionPlanController) GetAssetRiskActionPlanExport(c *gin.Context) {
	c.Set("export", true)

	arapc.GetAllAssetRiskActionPlans(c)
}
