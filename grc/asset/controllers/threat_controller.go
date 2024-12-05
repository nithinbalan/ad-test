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

type ThreatController struct {
	service services.ThreatServiceInterface
	db      *gorm.DB
}

func NewThreatController(service services.ThreatServiceInterface) *ThreatController {
	return &ThreatController{
		service: service,
	}
}

func (tc *ThreatController) GetAllThreats(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	detectionDateFrom := c.Query("detection_date_from")
	detectionDateTo := c.Query("detection_date_to")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, "Invalid export parameter")
		return
	}

	allowedParams := []string{
		"threat_type_ids",
		"severity_level_ids",
		"threat_status_ids",
		"responsible_unit_ids",
		"responsible_user_ids",
	}

	validSortFields := map[string]string{
		"id":                             "threats.id",
		"title_en":                       "threats.title_en",
		"title_ar":                       "threats.title_ar",
		"description":                    "threats.description",
		"threat_type":                    "threat_types.title_en",
		"severity_level":                 "severity_levels.title_en",
		"threat_status":                  "threat_statuses.title_en",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_unit":               "organization_units.title_en",
		"detection_date":                 "threats.detection_date",
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

	params := dtos.ThreatListQueryParamsDTO{
		Sort:              sort,
		SortBy:            dbField,
		Keyword:           keyword,
		Status:            status,
		Page:              page,
		PageSize:          pageSize,
		Filters:           filters,
		QuickFilter:       quickFilter,
		DetectionDateFrom: detectionDateFrom,
		DetectionDateTo:   detectionDateTo,
		Export:            export,
	}

	threats, err := tc.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}

	// Optional export functionality
	// if export {
	// 	excelFile, err := tc.service.ExportThreats(threats.Items)
	// 	if err != nil {
	// 		utility.SendInternalErrResp(c, err.Error())
	// 		return
	// 	}
	// 	defer excelFile.Close()

	// 	c.Header("Content-Disposition", `attachment; filename="Threats.xlsx"`)
	// 	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	// 	c.File(excelFile.Name())
	// 	return
	// }

	c.JSON(http.StatusOK, threats)
}

func (tc *ThreatController) GetThreatDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	threat, err := tc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, threat)
}

func (tc *ThreatController) CreateThreat(c *gin.Context) {
	var threatForm dtos.ThreatFormDTO
	if err := c.ShouldBindJSON(&threatForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}
	threatID, err := tc.service.Create(&threatForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": threatID})
}

func (tc *ThreatController) UpdateThreat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var threatForm dtos.ThreatFormDTO
	if err := c.ShouldBindJSON(&threatForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}

	err = tc.service.Update(uint(id), &threatForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tc *ThreatController) DeleteThreats(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}

	failedIDs, err := tc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tc *ThreatController) GetThreatsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"
	tc.GetAllThreats(c)
}

func (tc *ThreatController) GetThreatSummary(c *gin.Context) {
	idStr := c.Param("id")
	assetId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	summary, err := tc.service.GetThreatSummary(uint(assetId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}
