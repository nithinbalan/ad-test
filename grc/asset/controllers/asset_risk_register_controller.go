package controllers

import (
	"errors"
	"net/http"
	"strconv"

	assetRiskDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskRegisterController struct {
	service services.AssetRiskRegisterServiceInterface
	db      *gorm.DB
}

func NewAssetRiskRegisterController(service services.AssetRiskRegisterServiceInterface) *AssetRiskRegisterController {
	return &AssetRiskRegisterController{
		service: service,
	}
}

func (arrc *AssetRiskRegisterController) GetAllAssetRiskRegisters(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(arrc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"division_ids",
		"section_ids",
		"risk_category_ids",
		"current_risk_rating_ids",
		"is_strategic",
		"is_operational",
		"is_internal",
		"is_external",
	}

	validSortFields := map[string]string{
		"risk_code":         "risk_register.risk_code",
		"title":             "risk_register.title",
		"registration_date": "risk_register.registration_date",
		"risk_category_en":  "risk_categories.title_en",
		"risk_category_ar":  "risk_categories.title_ar",
		"division_en":       "division.title_en",
		"division_ar":       "division.title_ar",
		"section_en":        "section.title_en",
		"section_ar":        "section.title_ar",
		"risk_rating_en":    "current_risk_rating.title_en",
		"risk_rating_ar":    "current_risk_rating.title_ar",
		"created_at":        "risk_register.created_at",
		"updated_at":        "risk_register.updated_at",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := assetRiskDTOs.AssetRiskRegisterListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
		Export:   export,
	}

	risks, err := arrc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(arrc.db, c, err.Error())
			return
		}

		excelFile, err := arrc.service.ExportRiskRegister(risks.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(arrc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Asset Risk Register.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, risks)
}

func (arrc *AssetRiskRegisterController) GetAssetRiskRegisterExport(c *gin.Context) {
	c.Set("export", true)
	arrc.GetAllAssetRiskRegisters(c)
}

func (arrc *AssetRiskRegisterController) GetAssetRiskRegisterDetails(c *gin.Context) {
	riskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	risk, err := arrc.service.GetFullDetailsByID(uint(riskID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(arrc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, risk)
}

func (arrc *AssetRiskRegisterController) CreateAssetRiskRegister(c *gin.Context) {
	var dto assetRiskDTOs.AssetRiskRegisterFormDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	id, err := arrc.service.Create(currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(arrc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{
		ID: id,
	})
}

func (arrc *AssetRiskRegisterController) UpdateAssetRiskRegister(c *gin.Context) {
	var dto assetRiskDTOs.AssetRiskRegisterFormDTO

	riskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	err = arrc.service.Update(uint(riskID), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(arrc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (arrc *AssetRiskRegisterController) DeleteAssetRiskRegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	failedIDs, err := arrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"failed_ids": failedIDs,
	})
}
