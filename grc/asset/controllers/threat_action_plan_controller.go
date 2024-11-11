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

type ThreatActionPlanController struct {
	service services.ThreatActionPlanServiceInterface
	db      *gorm.DB
}

func NewThreatActionPlanController(service services.ThreatActionPlanServiceInterface) *ThreatActionPlanController {
	return &ThreatActionPlanController{
		service: service,
	}
}

func (tapc *ThreatActionPlanController) GetAllThreatActionPlans(c *gin.Context) {
	threatIDStr := c.Param("id")

	var threatID uint64
	var err error

	if threatIDStr != "" {
		threatID, err = strconv.ParseUint(threatIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		threatID = 0
	}
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(tapc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
		"responsible_user_ids",
		"action_plan_status_ids",
		"threat_register_ids",
	}

	validSortFields := map[string]string{
		"title":               "threat_action_plans.title",
		"threat_title_en":     "threat_register.title_en",
		"threat_title_ar":     "threat_register.title_ar",
		"severity_level_en":   "severity_levels.title_en",
		"severity_level_ar":   "severity_levels.title_ar",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"target_date":         "threat_action_plans.target_date",
		"percentage":          "threat_action_plans.percentage",
		"status_en":           "action_plan_statuses.title_en",
		"status_ar":           "action_plan_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tapc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := assetDTOs.ThreatActionPlanListQueryParamsDTO{
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
	}
	actionPlans, err := tapc.service.GetAll(uint(threatID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tapc.db, c, e.Error())
			return
		}
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(tapc.db, c, err.Error())
			return
		}

		excelFile, err := tapc.service.ExportThreatActionPlan(actionPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(tapc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Threat action plan.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, actionPlans)
}

func (tapc *ThreatActionPlanController) GetThreatActionPlanDetails(c *gin.Context) {
	threatIDStr := c.Param("id")

	threatID, err := strconv.ParseUint(threatIDStr, 10, 64)
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

	actionPlan, err := tapc.service.GetFullDetailsByID(uint(threatID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, actionPlan)
}

func (tapc *ThreatActionPlanController) CreateThreatActionPlan(c *gin.Context) {
	var formDTO assetDTOs.ThreatActionPlanFormDTO

	threatIDStr := c.Param("id")

	threatID, err := strconv.ParseUint(threatIDStr, 10, 64)
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
		utility.SendInternalErrResp(tapc.db, c, err.Error())
		return
	}

	id, err := tapc.service.Create(uint(threatID), currentUserID, &formDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (tapc *ThreatActionPlanController) UpdateThreatActionPlan(c *gin.Context) {
	var formDTO assetDTOs.ThreatActionPlanFormDTO

	threatIDStr := c.Param("id")

	threatID, err := strconv.ParseUint(threatIDStr, 10, 64)
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
		utility.SendInternalErrResp(tapc.db, c, err.Error())
		return
	}

	err = tapc.service.Update(uint(threatID), uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tapc *ThreatActionPlanController) DeleteThreatActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	threatIDStr := c.Param("id")

	threatID, err := strconv.ParseUint(threatIDStr, 10, 64)
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
		utility.SendInternalErrResp(tapc.db, c, err.Error())
		return
	}

	failedIDs, err := tapc.service.Delete(uint(threatID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tapc *ThreatActionPlanController) GetThreatActionPlanSummary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := tapc.service.GetThreatActionPlanSummary(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}

func (tapc *ThreatActionPlanController) GetThreatActionPlanExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	tapc.GetAllThreatActionPlans(c)
}
