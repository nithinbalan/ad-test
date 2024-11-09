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

type AssetMaintenanceController struct {
	service services.AssetMaintenanceServiceInterface
	db      *gorm.DB
}

func NewAssetMaintenanceController(service services.AssetMaintenanceServiceInterface) *AssetMaintenanceController {
	return &AssetMaintenanceController{
		service: service,
	}
}

func (amc *AssetMaintenanceController) GetAllMaintenances(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	performedDateFrom := c.Query("performed_date_from")
	performedDateTo := c.Query("performed_date_to")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"asset_register_ids",
		"performed_by_ids",
	}

	validSortFields := map[string]string{
		"title":             "asset_maintenance.title",
		"type":              "asset_maintenance.type",
		"performed_date":    "asset_maintenance.performed_date",
		"performed_by_name": "performed_by.first_name_en",
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

	params := dtos.AssetMaintenanceListQueryParamsDTO{
		Sort:              sort,
		SortBy:            dbField,
		Keyword:           keyword,
		Page:              page,
		PageSize:          pageSize,
		Filters:           filters,
		PerformedDateFrom: performedDateFrom,
		PerformedDateTo:   performedDateTo,
		Export:            export,
	}

	assetIdStr := c.Param("id")
	assetId, err := strconv.ParseUint(assetIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	maintenances, err := amc.service.GetAll(&params, uint(assetId))
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, maintenances)
}

func (amc *AssetMaintenanceController) GetMaintenanceDetails(c *gin.Context) {
	idStr := c.Param("maintenanceId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	maintenance, err := amc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, maintenance)
}

func (amc *AssetMaintenanceController) CreateMaintenance(c *gin.Context) {
	var maintenanceForm dtos.AssetMaintenanceFormDTO
	if err := c.ShouldBindJSON(&maintenanceForm); err != nil {
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
		utility.SendInternalErrResp(amc.db, c, err.Error())
		return
	}

	maintenanceID, err := amc.service.Create(&maintenanceForm, currentUserID, uint(assetId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(amc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": maintenanceID})
}

func (amc *AssetMaintenanceController) UpdateMaintenance(c *gin.Context) {
	idStr := c.Param("maintenanceId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var maintenanceForm dtos.AssetMaintenanceFormDTO
	if err := c.ShouldBindJSON(&maintenanceForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, err.Error())
		return
	}

	err = amc.service.Update(uint(id), &maintenanceForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(amc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (amc *AssetMaintenanceController) DeleteMaintenances(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, err.Error())
		return
	}

	failedIDs, err := amc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (amc *AssetMaintenanceController) GetMaintenancesExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"
	amc.GetAllMaintenances(c)
}
