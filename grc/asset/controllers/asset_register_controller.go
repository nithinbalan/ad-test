package controllers

import (
	"errors"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRegisterController struct {
	service services.AssetRegisterServiceInterface
	db      *gorm.DB
}

func NewAssetRegisterController(service services.AssetRegisterServiceInterface) *AssetRegisterController {
	return &AssetRegisterController{
		service: service,
	}
}

func (arc *AssetRegisterController) GetAllAssets(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	acquisitionDateFrom := c.Query("acquisition_date_from")
	acquisitionDateTo := c.Query("acquisition_date_to")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"asset_category_ids",
		"responsible_unit_ids",
		"responsible_user_ids",
		"asset_status_ids",
	}

	validSortFields := map[string]string{
		"reference_number":      "asset_register.reference_number",
		"title_en":              "asset_register.title_en",
		"title_ar":              "asset_register.title_ar",
		"description":           "asset_register.description",
		"location":              "asset_register.location",
		"asset_category_en":     "asset_categories.title_en",
		"asset_category_ar":     "asset_categories.title_ar",
		"asset_sub_category_en": "asset_sub_categories.title_en",
		"asset_sub_category_ar": "asset_sub_categories.title_ar",
		"asset_type_en":         "asset_types.title_en",
		"asset_type_ar":         "asset_types.title_ar",
		"responsible_user_en":   "responsible_user.first_name_en",
		"responsible_user_ar":   "responsible_user.first_name_ar",
		"responsible_unit_en":   "organization_units.title_en",
		"responsible_unit_ar":   "organization_units.title_ar",
		"acquisition_date":      "asset_register.acquisition_date",
		"assignment_date":       "asset_register.assignment_date",
		"warranty_expiry_date":  "asset_register.warranty_expiry_date",
		"manufacturer":          "asset_register.manufacturer",
		"model":                 "asset_register.model",
		"serial_number":         "asset_register.serial_number",
		"status_en":             "asset_statuses.title_en",
		"status_ar":             "asset_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

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

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	params := dtos.AssetRegisterListQueryParamsDTO{
		Sort:                sort,
		SortBy:              dbField,
		Keyword:             keyword,
		Status:              status,
		Page:                page,
		PageSize:            pageSize,
		Filters:             filters,
		QuickFilter:         quickFilter,
		AcquisitionDateFrom: acquisitionDateFrom,
		AcquisitionDateTo:   acquisitionDateTo,
		Export:              export,
	}

	assets, err := arc.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	// if export {
	// 	excelFile, err := arc.service.ExportAssets(assets.Items)
	// 	if err != nil {
	// 		utility.SendInternalErrResp(c, err.Error())
	// 		return
	// 	}
	// 	defer excelFile.Close()

	// 	c.Header("Content-Disposition", `attachment; filename="Assets.xlsx"`)
	// 	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	// 	c.File(excelFile.Name())
	// 	return
	// }

	c.JSON(http.StatusOK, assets)
}

func (arc *AssetRegisterController) GetAssetDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	asset, err := arc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, asset)
}

func (arc *AssetRegisterController) CreateAsset(c *gin.Context) {
	var assetForm dtos.AssetRegisterFormDTO
	if err := c.ShouldBindJSON(&assetForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}
	assetID, err := arc.service.Create(&assetForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": assetID})
}

func (arc *AssetRegisterController) UpdateAsset(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var assetForm dtos.AssetRegisterFormDTO
	if err := c.ShouldBindJSON(&assetForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	err = arc.service.Update(uint(id), &assetForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arc *AssetRegisterController) DeleteAssets(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	failedIDs, err := arc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (arc *AssetRegisterController) GetControlSummary(c *gin.Context) {
	assetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlSummary, err := arc.service.GetControlSummary(uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, controlSummary)
}

func (arc *AssetRegisterController) GetAssetsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"
	arc.GetAllAssets(c)
}

func (arc *AssetRegisterController) DownloadAssetDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetDocumentPath, err := arc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(arc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(assetDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(assetDocumentPath)
}
